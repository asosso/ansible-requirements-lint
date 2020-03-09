package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/atosatto/ansible-requirements-lint/requirements"
)

const (
	// DefaultAnsibleGalaxyURL is the URL of the upstream Ansible Galaxy
	// server managed by Red Hat.
	DefaultAnsibleGalaxyURL = "https://galaxy.ansible.com"
)

// AnsibleGalaxy fetches Ansible Roles information
// for the Ansible Galaxy APIs.
type AnsibleGalaxy struct {
	baseURL string
}

// NewAnsibleGalaxy creates a new AnsibleGalaxy provider.
// If baseURL is a nil string, DefaulAnsibleGalaxyURL
// will be used as baseURL for all the requests to the
// AnsibleGalaxy APIs.
func NewAnsibleGalaxy(baseURL string) AnsibleGalaxy {
	g := AnsibleGalaxy{}
	if baseURL == "" {
		g.baseURL = DefaultAnsibleGalaxyURL
	} else {
		g.baseURL = baseURL
	}
	return g
}

// VersionsForRole returns the list of versions available on AnsibleGalaxy for the Role r.
func (g AnsibleGalaxy) VersionsForRole(ctx context.Context, r requirements.Role) ([]string, error) {
	client := &http.Client{Timeout: time.Second * 10}

	// Ansible Galaxy URL
	baseURL, err := url.Parse(g.baseURL + "/api/v1/search/roles/")
	if err != nil {
		return nil, err
	}

	// keywords to be used for the search on Ansible Galaxy
	var keywords string
	if len(r.Source) != 0 {
		keywords = r.Source
	} else {
		keywords = r.Name
	}

	// namespace to be used to filter the Ansible Galaxy results
	var namespace string
	split := strings.Split(keywords, ".")
	if len(split) > 0 {
		namespace = split[0]
	}

	// set the Ansible Galaxy search parameters
	params := url.Values{}
	params.Add("keywords", keywords)
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "ansible-requirements-lint")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected Ansible Galaxy response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type galaxyResult struct {
		SummaryFields struct {
			Versions []struct {
				Name string `json:"name"`
			} `json:"versions"`
			Namespace struct {
				Name string `json:"name"`
			} `json:"namespace"`
		} `json:"summary_fields"`
	}

	var results struct {
		Count   int            `json:"count"`
		Results []galaxyResult `json:"results"`
	}
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	// search for the best match in the results list
	var matching *galaxyResult
	if len(results.Results) == 1 {
		matching = &results.Results[0]
	} else {
		for _, r := range results.Results {
			if r.SummaryFields.Namespace.Name == namespace {
				matching = &r
				break
			}
		}
	}
	if matching == nil {
		return nil, fmt.Errorf("unable to find role in Ansible Galaxy: %s", keywords)
	}

	// get the latest version of the role
	versions := make([]string, len(matching.SummaryFields.Versions))
	for i, v := range matching.SummaryFields.Versions {
		versions[i] = v.Name
	}
	return versions, nil
}

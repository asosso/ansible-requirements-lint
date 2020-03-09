# ansible-requirements-lint - keep you Ansible dependencies up to date

`ansible-requirements-lint` is a simple command-line tool to check if your Ansible dependencies are up to date.

[![GoDoc](https://godoc.org/github.com/atosatto/ansible-requirements-lint?status.svg)](https://godoc.org/github.com/atosatto/ansible-requirements-lint)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/atosatto/ansible-requirements-lint)](https://goreportcard.com/report/github.com/atosatto/ansible-requirements-lint)
![GitHub All Releases](https://img.shields.io/github/downloads/atosatto/ansible-requirements-lint/total)

## Installation

You can download the latest release of `ansible-requirements-lint` from
https://github.com/atosatto/ansible-requirements-lint/releases/latest.

## Usage

Given the following `requirements.yml` file in your current working directory

```bash
$ cat requirements.yml

# Prometheus
- name: atosatto.prometheus
  version: v1.0.0

# Alertmanager
- name: atosatto.alertmanager
  version: v1.0.0

# Grafana
- name: atosatto.grafana
  version: v1.0.0
```

`ansible-requirements-lint` can be used to detect updates to the list of requirements with

```bash
$ bin/ansible-requirements-lint requirements.yml
WARN: atosatto.prometheus: role not at the latest version, upgrade to v1.1.0.
WARN: atosatto.grafana: role not at the latest version, upgrade to v1.1.0.
```

# Terraform ServiceNow Provider

A Terraform provider for managing configuration in [ServiceNow](https://www.servicenow.com/) instances via the Table API. Covers **63 resources** and **49 data sources** spanning users and groups, ACLs and roles, UI scripts and macros, business rules, scripted REST APIs, service catalog items, scheduled jobs, REST/JDBC/HTTP connections, system properties, transform maps, widgets, and more.

[![Release](https://img.shields.io/github/v/release/tylerhatton/terraform-provider-servicenow?sort=semver)](https://github.com/tylerhatton/terraform-provider-servicenow/releases/latest)
[![Registry](https://img.shields.io/badge/registry-tylerhatton%2Fservicenow-623CE4?logo=terraform)](https://registry.terraform.io/providers/tylerhatton/servicenow/latest)
[![Build](https://github.com/tylerhatton/terraform-provider-servicenow/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/tylerhatton/terraform-provider-servicenow/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tylerhatton/terraform-provider-servicenow)](https://goreportcard.com/report/github.com/tylerhatton/terraform-provider-servicenow)
[![License](https://img.shields.io/github/license/tylerhatton/terraform-provider-servicenow.svg)](LICENSE)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) `>= 1.0`
- A ServiceNow instance and an account with permission to read/write the tables you intend to manage
- For local development only: [Go](https://golang.org/doc/install) `>= 1.25`

## Usage

```hcl
terraform {
  required_providers {
    servicenow = {
      source  = "tylerhatton/servicenow"
      version = "~> 0.10"
    }
  }
}

provider "servicenow" {
  instance_url = "https://dev00000.service-now.com/"
  username     = var.servicenow_username
  password     = var.servicenow_password
}

resource "servicenow_user" "jane" {
  user_name  = "jane.doe"
  first_name = "Jane"
  last_name  = "Doe"
  email      = "jane.doe@example.com"
  active     = true
}
```

Full documentation for every resource and data source is published on the [Terraform Registry](https://registry.terraform.io/providers/tylerhatton/servicenow/latest/docs).

## Authentication

The provider uses HTTP Basic authentication against the ServiceNow Table API. Credentials can be supplied in the provider block or — recommended — via environment variables:

| Variable | Maps to |
|---|---|
| `SERVICENOW_INSTANCE_URL` | `instance_url` |
| `SERVICENOW_USERNAME` | `username` |
| `SERVICENOW_PASSWORD` | `password` |

```sh
export SERVICENOW_INSTANCE_URL="https://dev00000.service-now.com/"
export SERVICENOW_USERNAME="admin"
export SERVICENOW_PASSWORD="..."
terraform plan
```

The user must have `soap` and `web_service_admin` roles (or equivalent table-level ACLs) to manage records through the API.

## Importing existing records

Every managed resource supports `terraform import` using the record's `sys_id`:

```sh
terraform import servicenow_user.jane 5137153cc611227c000bbd1bd8cd2007
```

## Building from source

```sh
git clone https://github.com/tylerhatton/terraform-provider-servicenow.git
cd terraform-provider-servicenow
make install
```

`make install` builds the provider and places the binary under `~/.terraform.d/plugins/tyler.sh/tylerhatton/servicenow/<version>/<os_arch>/` for use with a local override. To use the local build, point `required_providers.servicenow.source` at `tyler.sh/tylerhatton/servicenow`.

## Development

```sh
make build      # build the provider
make test       # unit tests with mocked client
make testacc    # acceptance tests against a live instance (requires TF_ACC=1 and the env vars above)
```

Acceptance tests hit a real ServiceNow instance and create/destroy records — use a developer instance, not production.

## Contributing

Issues and pull requests are welcome. For non-trivial changes, please open an issue first to discuss the approach.

## License

[MIT](LICENSE) © Tyler Hatton

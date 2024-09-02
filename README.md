# Qwilt Terraform Provider

> ⚠️ **Disclaimer**: the project is still in the 0.x.x version, which means it’s still in the experimental phase.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the Provider

## Developing the Provider

Configure `~/.terraformrc` with the name of your provider and go
installation. For example:

```
provider_installation {

  dev_overrides {
      "qwilt.com/qwiltinc/qwilt" = "/path/to/binary/of/qwilt/provider"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
Set the local provider name in main.go:
```
providerName = "qwilt.com/qwiltinc/qwilt"
```

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. Don'r run it on production environment.

```shell
make testacc
```

# Terraform Provider for Twitter

Website: [registry.terraform.io/providers/sebastianmarines/twitter/latest](https://registry.terraform.io/providers/sebastianmarines/twitter/latest)

Documentation: [registry.terraform.io/providers/sebastianmarines/twitter/latest/docs](https://registry.terraform.io/providers/sebastianmarines/twitter/latest/docs)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.17

## Using the provider

### Installation

First install the provider and run `terraform init` to download it

```hcl
terraform {
  required_providers {
    twitter = {
      source  = "hashicorp.com/local/twitter"
      version = "~> 0.1.0"
    }
  }
}
```

### Configuration

To configure the provider set the following variables in the provider configuration:

```hcl
provider "twitter" {
  api_key             = "YOUR_API_KEY"
  api_secret          = "YOUR_API_SECRET"
  access_token        = "YOUR_ACCESS_TOKEN"
  access_token_secret = "YOUR_ACCESS_TOKEN_SECRET"
}
```

Or set the following environment variables:

- TWITTER_API_KEY
- TWITTER_API_SECRET_KEY
- TWITTER_ACCESS_TOKEN
- TWITTER_ACCESS_TOKEN_SECRET

> In order to get the required keys go to https://developer.twitter.com/ and apply for a developer account

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make install`. This will build the provider and put the provider binary in your current directory.

To install the provider for local development run `make install`, this will build the provider and copy the binary to the `~/.terraform.d/plugins/` directory.

To use the locally installed provider specify the required provider like this:
```hcl
terraform {
  required_providers {
    twitter = {
      source  = "hashicorp.com/local/twitter"
      version = "~> 0.1.0"
    }
  }
}
```

To generate or update documentation, run `go generate`.


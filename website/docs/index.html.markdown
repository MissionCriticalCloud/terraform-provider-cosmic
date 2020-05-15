---
layout: "cosmic"
page_title: "Provider: Cosmic"
sidebar_current: "docs-cosmic-index"
description: |-
  The Cosmic provider is used to interact with the many resources supported by Cosmic. The provider needs to be configured with a URL pointing to a running Cosmic API and the proper credentials before it can be used.
---

# Cosmic Provider

The Cosmic provider is used to interact with the many resources
supported by Cosmic. The provider needs to be configured with a
URL pointing to a running Cosmic API and the proper credentials
before it can be used.

In order to provide the required configuration options you can either
supply values for the `api_url`, `api_key` and `secret_key` fields, or
for the `config` and `profile` fields. A combination of both is not
allowed and will not work.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Cosmic Provider
provider "cosmic" {
  api_url    = "${var.cosmic_api_url}"
  api_key    = "${var.cosmic_api_key}"
  secret_key = "${var.cosmic_secret_key}"
  zone       = "${var.cosmic_zone}"
}

# Create a web server
resource "cosmic_instance" "web" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `api_url` - (Optional) This is the Cosmic API URL. It can also be sourced
  from the `COSMIC_API_URL` environment variable.

* `api_key` - (Optional) This is the Cosmic API key. It can also be sourced
  from the `COSMIC_API_KEY` environment variable.

* `secret_key` - (Optional) This is the Cosmic secret key. It can also be
  sourced from the `COSMIC_SECRET_KEY` environment variable.

* `zone` - (Optional) This is the Cosmic zone It can also be
  sourced from the `COSMIC_ZONE` environment variable.

* `config` - (Optional) The path to a `CloudMonkey` config file. If set the API
  URL, key, secret and zone will be retrieved from this file. It can also be
  sourced from the `COSMIC_CONFIG` environment variable.

* `profile` - (Optional) Used together with the `config` option. Specifies which
  `CloudMonkey` profile in the config file to use. It can also be
  sourced from the `COSMIC_PROFILE` environment variable.

* `http_get_only` - (Optional) Some cloud providers only allow HTTP GET calls to
  their Cosmic API. If using such a provider, you need to set this to `true`
  in order for the provider to only make GET calls and no POST calls. It can also
  be sourced from the `COSMIC_HTTP_GET_ONLY` environment variable.

* `timeout` - (Optional) A value in seconds. This is the time allowed for Cosmic
  to complete each asynchronous job triggered. If unset, this can be sourced from the
  `COSMIC_TIMEOUT` environment variable. Otherwise, this will default to 300
  seconds.

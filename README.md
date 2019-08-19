Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/MissionCriticalCloud/terraform-provider-cosmic`

```sh
$ mkdir -p $GOPATH/src/github.com/MissionCriticalCloud; cd $GOPATH/src/github.com/MissionCriticalCloud
$ git clone git@github.com:MissionCriticalCloud/terraform-provider-cosmic
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/MissionCriticalCloud/terraform-provider-cosmic
$ make build
```

Using the provider
------------------
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-cosmic
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
$ make testacc
```

*Note:* Acceptance tests create real resources and require the following shell variables to be exported:

- `COSMIC_API_URL` -  Cosmic API URL to run tests against
- `COSMIC_API_KEY` - Valid API key
- `COSMIC_SECRET_KEY` - Valid API secret key
- `COSMIC_ZONE` - Cosmic zone to run tests against

Optionally the following need to be exported for certain tests:

- `COSMIC_DISK_OFFERING_1` -  An existing disk offering to test storage
- `COSMIC_DISK_OFFERING_2` -  A second existing disk offering to test provisioning storage
- `COSMIC_INSTANCE_ID` - An existing instance to test instance-related resources
- `COSMIC_SERVICE_OFFERING_1` - An existing service offering to test provisioning instances
- `COSMIC_SERVICE_OFFERING_2` - A second existing service offering to test provisioning instances
- `COSMIC_TEMPLATE` - An existing template to test provisioning instances
- `COSMIC_VPC_ID` -  An existing VPC ID to test provisioning VPC resources
- `COSMIC_VPC_NETWORK_OFFERING` - An existing VPC network offering to test provisioning VPC networks

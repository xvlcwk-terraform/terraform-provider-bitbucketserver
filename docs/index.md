# Bitbucket Server Provider

[Bitbucket Server](https://www.atlassian.com/software/bitbucket) is the self-hosted version of Bitbucket.
Whilst terraform provides a default bitbucket provider, this only works for _Bitbucket Cloud_ - this provider
unlocks the power of terraform to manage your self-hosted Bitbucket Server instance. 

## Installation

### Terraform 1.0+

The provider can be installed and managed automatically by Terraform. Sample `versions.tf` file :

```hcl
terraform {
  required_providers {
    bitbucketserver = {
      source  = "xvlcwk-terraform/bitbucketserver"
    }
  }
}
```

#### Install manually

If you don't want to use the one-liner above, you can download a binary for your system from the [release page](https://github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/releases),
then either place it at the root of your Terraform folder or in the Terraform plugin folder on your system.


## Configuration

The provider supports parameters to determine the bitbucket server and admin user/password to use.

```hcl
provider "bitbucketserver" {
  server   = "https://mybitbucket.example.com"
  username = "admin"
  password = "password"
}
```

Alternatively you can specify a http-access-token, but please be aware that those aren't able to create other access tokens and are overall limited in their scope.


```hcl
provider "bitbucketserver" {
  server   = "https://mybitbucket.example.com"
  token = "BBDC-ODc4NTMyODU1MTI3OnSOs8sHd7EZneuWx2ZiYyF/Xmuj"
}
```

### Authentication

The `username` and `password` specified should be of a user with sufficient privileges to perform the operations you are after.
Typically, this is a user with `SYS_ADMIN` global permissions.

### Environment Variables

You can also specify the provider configuration using the following env vars:

* `BITBUCKET_SERVER`
* `BITBUCKER_USERNAME`
* `BITBUCKET_PASSWORD`
* `BITBUCKET_TOKEN`

> Note: The hcl provider configuration takes precedence over the environment variables.

## Example - Creating a Project and Repository

Creating a project and repository is super simple with this provider:

```hcl
provider "bitbucketserver" {
  server   = "https://mybitbucket.example.com"
  username = "admin"
  password = "password"
}

resource "bitbucketserver_project" "test" {
  key         = "TEST"
  name        = "test-01"
  description = "Test project"
}

resource "bitbucketserver_repository" "test" {
  project     = bitbucketserver_project.test.key
  name        = "repo-01"
  description = "Test repository"
}
```

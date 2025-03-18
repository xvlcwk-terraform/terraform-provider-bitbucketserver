# Resource: bitbucketserver_user_access_token

Repository access tokens can be used  to authenticate using the Bitbucket Server REST API over Token auth. 

## Example Usage

```hcl
resource "bitbucketserver_repository_access_token" "test" {
  project = bitbucketserver_project.test.key
  repository = bitbucketserver_repository.test.slug
  name = "newLabelForTest"
  permissions = ["REPO_READ"]
}
```

## Argument Reference

* `project` - Required. Project for the repository.
* `repository` - Required. Slug for the repository.
* `name` - Required. Name of the access token.
* `permissions` - Required. List of permissions to grant the access token.
     * `REPO_READ`
     * `REPO_WRITE`
     * `REPO_ADMIN`

## Attribute Reference

* `token` - The generated access token. Only available if token was generated on Terraform resource creation, not import/update.

## Import

Currently not supported
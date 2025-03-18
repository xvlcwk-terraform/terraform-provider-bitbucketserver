# Resource: bitbucketserver_user_access_token

Repository access tokens can be used to authenticate using the Bitbucket Server REST API over Token auth.

## Example Usage

```hcl
resource "bitbucketserver_repository_deploy_key" "test" {
  project    = bitbucketserver_project.test.key
  repository = bitbucketserver_repository.test.slug
  label      = "newLabelForTest"
  key        = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQD1F2KTOi8ptpqkXdcARgYy27uWRCav1isJUB3Hz59MDM6CSAXa+HCgUNDV+6gXQ1eR24nr1Efb7AkEM8LmvMkNlQqpAsPEPxVlndA0KSRLXb3mWruJzkZtrJo1HhVDffuKnevPLFkB75wyBX3jn1FNtc2qfc2i80mu8TfviZnOPYnrRa8B2S9Q8IlUtXAyihQ1G3fn5r5nxrw3QQGrD3cckB9nCZpzAPn6hlDHVk6n5efoWAUxM5AhKln0OLsA1HwjGJN7/dbDPu1nSLuiJSAISSSg4E4SNdfnr3FOhTA79AKNzsor2/EXdbG+f+S4op3s3wtt05zwkHLXRSqKYoT31RFiV9d1XIav+dGTvXCgc6DNG6rbogE6ZbugDZXdcHOAoNs7IUDFbtI/HGKS7CStxAUEchoxM8HDYXmhYt1kUEpmP3g2ckILHGoGPEkOCYGPqx5HbDvAXAJVk3DdSOibCckR2FK2qEoCbMgnPUX84CqNJPHBZ24AaE8htE6TKr0= user@somewhere"
  permission = "REPO_READ"
}
```

## Argument Reference

* `project` - Required. Project for the repository.
* `repository` - Required. Slug for the repository.
* `label` - Required. Name of the deploy key token.
* `key` - Required. The public key to grant access to
* `permissions` - Required. List of permissions to grant the access token.
    * `REPO_READ`
    * `REPO_WRITE`
    * `REPO_ADMIN`
* expiry_days - Optional. Set if the key should expire. 0 means "does not expire" and is the default

## Attribute Reference

## Import

Currently not supported
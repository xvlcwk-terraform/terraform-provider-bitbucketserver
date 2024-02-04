package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketResourceRepositoryAccessToken(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-project-%v"
		}

		resource "bitbucketserver_repository" "test" {
			project = bitbucketserver_project.test.key
			name = "repo"
		}

		resource "bitbucketserver_repository_access_token" "test" {
			project = bitbucketserver_project.test.key
			repository = bitbucketserver_repository.test.slug
			name = "newLabelForTest"
			permissions = ["REPO_READ"]
		}
	`, projectKey, projectKey)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucketserver_repository_access_token.test", "id"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_access_token.test", "project", projectKey),
					resource.TestCheckResourceAttr("bitbucketserver_repository_access_token.test", "repository", "repo"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_access_token.test", "name", "newLabelForTest"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_access_token.test", "permissions.#", "1"),
					resource.TestCheckResourceAttr("bitbucketserver_repository_access_token.test", "permissions.0", "REPO_READ"),
					resource.TestCheckResourceAttrSet("bitbucketserver_repository_access_token.test", "token"),
				),
			},
		},
	})
}

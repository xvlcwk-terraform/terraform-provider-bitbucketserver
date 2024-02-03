package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketDataUser(t *testing.T) {
	config := `
		data "bitbucketserver_user" "test" {
			name = "admin"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_user.test", "name", "admin"),
					resource.TestCheckResourceAttr("data.bitbucketserver_user.test", "email_address", "admin@example.com"),
					resource.TestCheckResourceAttr("data.bitbucketserver_user.test", "display_name", "Admin"),
					// not sure, why admin has user_id 2 by default, but it seems to be the case since the bitbucket server update
					resource.TestCheckResourceAttr("data.bitbucketserver_user.test", "user_id", "2"),
				),
			},
		},
	})
}

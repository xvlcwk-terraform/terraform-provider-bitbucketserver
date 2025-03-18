package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketDataApplicationProperties(t *testing.T) {
	config := `
		data "bitbucketserver_application_properties" "main" {}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bitbucketserver_application_properties.main", "version", "9.4.2"),
					resource.TestCheckResourceAttr("data.bitbucketserver_application_properties.main", "build_number", "9004002"),
					resource.TestCheckResourceAttr("data.bitbucketserver_application_properties.main", "build_date", "1736803793767"),
					resource.TestCheckResourceAttr("data.bitbucketserver_application_properties.main", "display_name", "Bitbucket"),
				),
			},
		},
	})
}

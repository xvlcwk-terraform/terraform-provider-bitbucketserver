package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/client"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestAccBitbucketResourceGroup_basic(t *testing.T) {
	config := `
		resource "bitbucketserver_group" "test" {
			name = "test-group"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_group.test", "name", "test-group"),
				),
			},
		},
	})
}

func TestAccBitbucketResourceGroup_DisallowImport(t *testing.T) {
	resourceName := "duplicate_group"
	groupName := "duplicate-group"
	config := fmt.Sprintf(`
		resource "bitbucketserver_group" "%s" {
			name = "%s"
		}
	`, resourceName, groupName)

	createGroup(groupName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("bitbucketserver_group.%s", resourceName), "name", groupName),
				),
				ExpectError: regexp.MustCompile("API Error: 409"),
			},
		},
	})
}

func TestAccBitbucketResourceGroup_AllowImport(t *testing.T) {
	resourceName := "duplicate_group"
	groupName := "duplicate-group"
	config := fmt.Sprintf(`
		resource "bitbucketserver_group" "%s" {
			name = "%s"
			import_if_exists = true
		}
	`, resourceName, groupName)

	createGroup(groupName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("bitbucketserver_group.%s", resourceName), "name", groupName),
				),
			},
		},
	})
}

func createGroup(groupName string) {
	client := newBitbucketClient()
	client.Post(fmt.Sprintf("/rest/api/1.0/admin/groups?name=%s", url.QueryEscape(groupName)), nil)
}

func newBitbucketClient() *client.BitbucketClient {
	serverSanitized := os.Getenv("BITBUCKET_SERVER")
	if strings.HasSuffix(serverSanitized, "/") {
		serverSanitized = serverSanitized[0 : len(serverSanitized)-1]
	}

	return &client.BitbucketClient{
		Server:     serverSanitized,
		Username:   os.Getenv("BITBUCKET_USERNAME"),
		Password:   os.Getenv("BITBUCKET_PASSWORD"),
		HTTPClient: &http.Client{},
	}
}

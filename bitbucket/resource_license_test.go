package bitbucket

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBitbucketLicense(t *testing.T) {
	// Licenses can be found at: https://developer.atlassian.com/platform/marketplace/timebomb-licenses-for-testing-server-apps/
	testAccBitbucketLicenseConfig := `
		resource "bitbucketserver_license" "test" {
			license = "AAABrQ0ODAoPeNp9kVFvmzAQx9/9KSztLZIJZIu0RkJqA6yNViAK0G3d+uDApXgjNrKPbPn2dYG0a6fuwS8+393v//O7vAMa8yN1PerOFrPZYn5GL+OczlzvI0m6/RZ0uisMaON7LgmURF5iwvfgVy3XWpj6nGPDjRFcOqXaE4Pc1M61KEEayI8t9I+DNI6jTbC6uP73wd/FdafLmhsIOYL/yMDcOXM98p95Yyn60wp97PvW769OpFHMRfMWagb6AHoV+svLs5x9LW4+sM+3t1ds6XpfRkw7jwcgEbSPugOSdVtTatGiUHK4mUwmSZqzT+mGrTdpWAT5Kk1YkUW24AcaLFBFt0eKNdARlUayVBVo2mr1E0qk32vE9sdiOr1XzgvEaTN0MBg67hwaKioV0koY1GLbIdjJwlBUtOwMqr39KYfY1JZZclm+9jLEsmbEAZ4CBJvoIo9Ctvz2CP2GrRHe6irkL6l+S5JFiW8Pm7suSfU9l8LwXkwIB2hUaxPmYPAUm/Q2bP315w5MGXL95DmEZ839jFEE3SlNedvS6rTCkOjAm25YvOON3fMAVTj4nTAtAhRH4o+fI5MQ7xSh2mtA1bPJrq0WAgIVAIGperR8m2N0fl/GfUUJfQnd+T1aX02kk"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketLicenseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketLicenseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketLicenseExists("bitbucketserver_license.test"),
				),
			},
		},
	})
}

func testAccCheckBitbucketLicenseDestroy(s *terraform.State) error {
	_, ok := s.RootModule().Resources["bitbucketserver_license.test"]
	if !ok {
		return fmt.Errorf("not found %s", "bitbucketserver_license.test")
	}

	return nil
}

func testAccCheckBitbucketLicenseExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no license ID is set")
		}
		return nil
	}
}

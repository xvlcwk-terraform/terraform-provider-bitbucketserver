package bitbucket

import (
	"fmt"
	bitbucketTypes "github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/types"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBitbucketRepository_basic(t *testing.T) {
	var repo Repository

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_repository" "test_repo" {
			project = bitbucketserver_project.test.key
			name = "test-repo-for-repository-test"
			description = "My Repo"
			forkable = false
			public = false
		}

		resource "bitbucketserver_repository" "test_repo_defaults" {
			project = bitbucketserver_project.test.key
			name = "test-repo-for-repository-test_defaults"
			description = "My_Repo2"
		}
	`, rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	configModified := strings.ReplaceAll(config, "My Repo", "My Updated Repo")
	configModifiedBool := strings.ReplaceAll(config, "false", "true")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucketserver_repository.test_repo", &repo),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "slug", "test-repo-for-repository-test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "name", "test-repo-for-repository-test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "description", "My Repo"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "forkable", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "public", "false"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo_defaults", "forkable", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo_defaults", "public", "false"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucketserver_repository.test_repo", &repo),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "slug", "test-repo-for-repository-test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "name", "test-repo-for-repository-test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "description", "My Updated Repo"),
				),
			},
			{
				Config: configModifiedBool,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucketserver_repository.test_repo", &repo),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "slug", "test-repo-for-repository-test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "name", "test-repo-for-repository-test"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "forkable", "true"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_repo", "public", "true"),
				),
			},
		},
	})
}

func TestAccBitbucketRepository_namewithspaces(t *testing.T) {
	var repo Repository

	testAccBitbucketRepositoryConfig := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "test-repo-for-repository-test"
		}

		resource "bitbucketserver_repository" "test_repo" {
			project = bitbucketserver_project.test.key
			name = "Test Repo For Repository Test"
			slug = "test-repo-for-repository-test"
		}
	`, rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketRepositoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucketserver_repository.test_repo", &repo),
				),
			},
		},
	})
}

func TestAccBitbucketRepository_fork(t *testing.T) {
	var repo Repository

	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "Test-%v"
		}

		resource "bitbucketserver_repository" "test_repo" {
			project = bitbucketserver_project.test.key
			name = "test-repo-for-repository-test"
			description = "My Repo"
		}

		resource "bitbucketserver_repository" "test_fork" {
			project = bitbucketserver_repository.test_repo.project
			name = "My Fork"
			description = "My Repo Forked"
			fork_repository_project = bitbucketserver_repository.test_repo.project
			fork_repository_slug = bitbucketserver_repository.test_repo.slug
		}
	`, rand.New(rand.NewSource(time.Now().UnixNano())).Int(), rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	configModified := strings.ReplaceAll(config, "My Repo Forked", "My Updated Repo Forked")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucketserver_repository.test_fork", &repo),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_fork", "slug", "my-fork"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_fork", "name", "My Fork"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_fork", "description", "My Repo Forked"),
				),
			},
			{
				Config: configModified,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucketserver_repository.test_fork", &repo),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_fork", "slug", "my-fork"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_fork", "name", "My Fork"),
					resource.TestCheckResourceAttr("bitbucketserver_repository.test_fork", "description", "My Updated Repo Forked"),
				),
			},
		},
	})
}

func TestAccBitbucketRepository_gitlfs(t *testing.T) {
	var repo Repository

	key := fmt.Sprintf("%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	testAccBitbucketRepositoryConfig := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "TEST%v"
			name = "Test%v"
		}

		resource "bitbucketserver_repository" "test_repo" {
			project = bitbucketserver_project.test.key
			name = "test-repo-for-repository-test"
			enable_git_lfs = true
		}
	`, key, key)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketRepositoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucketserver_repository.test_repo", &repo),
				),
			},
		},
	})
}

func testAccCheckBitbucketRepositoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*bitbucketTypes.BitbucketServerProvider).BitbucketClient
	rs, ok := s.RootModule().Resources["bitbucketserver_repository.test_repo"]
	if !ok {
		return fmt.Errorf("not found %s", "bitbucketserver_repository.test_repo")
	}

	response, _ := client.Get(fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s", rs.Primary.Attributes["project"], rs.Primary.Attributes["slug"]))

	if response.StatusCode != 404 {
		return fmt.Errorf("repository still exists")
	}

	return nil
}

func testAccCheckBitbucketRepositoryExists(n string, repository *Repository) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no repository ID is set")
		}
		return nil
	}
}

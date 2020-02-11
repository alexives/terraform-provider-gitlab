package gitlab

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectPagesDomain_basic(t *testing.T) {
	var pagesDomain gitlab.PagesDomain
	rString := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectPagesDomainDestroy,
		Steps: []resource.TestStep{
			// Create a project and variable with default options
			{
				Config: testAccGitlabProjectPagesDomainConfig(rString, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPagesDomainExists("gitlab_project_variable.foo", &pagesDomain),
					testAccCheckGitlabProjectPagesDomainAttributes(&pagesDomain, &testAccGitlabProjectPagesDomainExpectedAttributes{
						Domain:         fmt.Sprintf("domain_%s", rString),
						AutoSslEnabled: true,
					}),
				),
			},
			// Update the project variable to toggle all the values to their inverse
			{
				Config: testAccGitlabProjectPagesDomainUpdateConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPagesDomainExists("gitlab_project_variable.foo", &pagesDomain),
					testAccCheckGitlabProjectPagesDomainAttributes(&pagesDomain, &testAccGitlabProjectPagesDomainExpectedAttributes{
						Domain: fmt.Sprintf("domain_%s", rString),
					}),
				),
			},
			// Update the project variable to toggle the options back
			{
				Config: testAccGitlabProjectPagesDomainConfig(rString, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPagesDomainExists("gitlab_project_variable.foo", &pagesDomain),
					testAccCheckGitlabProjectPagesDomainAttributes(&pagesDomain, &testAccGitlabProjectPagesDomainExpectedAttributes{
						Domain: fmt.Sprintf("domain_%s", rString),
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabProjectPagesDomainExists(n string, projectVariable *gitlab.PagesDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		domain := rs.Primary.Attributes["domain"]
		if domain == "" {
			return fmt.Errorf("No variable domain is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		gotPagesDomain, _, err := conn.PagesDomains.GetPagesDomain(repoName, domain)
		if err != nil {
			return err
		}
		*projectVariable = *gotPagesDomain
		return nil
	}
}

type testAccGitlabProjectPagesDomainExpectedAttributes struct {
	Domain           string
	URL              string
	Verified         bool
	VerificationCode string
	AutoSslEnabled   bool
	Certificate      struct {
		Expired    bool
		Expiration *time.Time
	}
}

func testAccCheckGitlabProjectPagesDomainAttributes(pagesDomain *gitlab.PagesDomain, want *testAccGitlabProjectPagesDomainExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pagesDomain.Domain != want.Domain {
			return fmt.Errorf("got domain %s; want %s", pagesDomain.Domain, want.Domain)
		}

		if pagesDomain.Verified != want.Verified {
			return fmt.Errorf("got verified %t; want %t", pagesDomain.Verified, want.Verified)
		}

		if pagesDomain.VerificationCode != want.VerificationCode {
			return fmt.Errorf("got verification code %s; want %s", pagesDomain.VerificationCode, want.VerificationCode)
		}

		return nil
	}
}

func testAccCheckGitlabProjectPagesDomainDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := conn.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				return fmt.Errorf("Repository still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabProjectPagesDomainConfig(rString string, autoSslEnabled bool) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%s"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_pages_domain" "foo" {
  project = "${gitlab_project.foo.id}"
	domain = "domain_%s"
	certificate = "cert_%s"
	key = "key_%s"
	auto_ssl_enabled = %t
}
	`, rString, rString, rString, rString, autoSslEnabled)
}

func testAccGitlabProjectPagesDomainUpdateConfig(rString string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%s"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_project_variable" "foo" {
  project = "${gitlab_project.foo.id}"
  key = "key_%s"
  value = "value-inverse-%s"
  protected = true
}
	`, rString, rString, rString)
}

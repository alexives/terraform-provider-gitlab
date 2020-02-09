package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectPagesDomain_basic(t *testing.T) {
	var projectVariable gitlab.ProjectVariable
	rString := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectPagesDomainDestroy,
		Steps: []resource.TestStep{
			// Create a project and variable with default options
			{
				Config: testAccGitlabProjectPagesDomainConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPagesDomainExists("gitlab_project_variable.foo", &projectVariable),
					testAccCheckGitlabProjectPagesDomainAttributes(&projectVariable, &testAccGitlabProjectPagesDomainExpectedAttributes{
						Key:   fmt.Sprintf("key_%s", rString),
						Value: fmt.Sprintf("value-%s", rString),
					}),
				),
			},
			// Update the project variable to toggle all the values to their inverse
			{
				Config: testAccGitlabProjectPagesDomainUpdateConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPagesDomainExists("gitlab_project_variable.foo", &projectVariable),
					testAccCheckGitlabProjectPagesDomainAttributes(&projectVariable, &testAccGitlabProjectPagesDomainExpectedAttributes{
						Key:       fmt.Sprintf("key_%s", rString),
						Value:     fmt.Sprintf("value-inverse-%s", rString),
						Protected: true,
					}),
				),
			},
			// Update the project variable to toggle the options back
			{
				Config: testAccGitlabProjectPagesDomainConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPagesDomainExists("gitlab_project_variable.foo", &projectVariable),
					testAccCheckGitlabProjectPagesDomainAttributes(&projectVariable, &testAccGitlabProjectPagesDomainExpectedAttributes{
						Key:       fmt.Sprintf("key_%s", rString),
						Value:     fmt.Sprintf("value-%s", rString),
						Protected: false,
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabProjectPagesDomainExists(n string, projectVariable *gitlab.ProjectVariable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		key := rs.Primary.Attributes["key"]
		if key == "" {
			return fmt.Errorf("No variable key is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		gotVariable, _, err := conn.ProjectVariables.GetVariable(repoName, key)
		if err != nil {
			return err
		}
		*projectVariable = *gotVariable
		return nil
	}
}

type testAccGitlabProjectPagesDomainExpectedAttributes struct {
	Key       string
	Value     string
	Protected bool
}

func testAccCheckGitlabProjectPagesDomainAttributes(variable *gitlab.ProjectVariable, want *testAccGitlabProjectPagesDomainExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if variable.Key != want.Key {
			return fmt.Errorf("got key %s; want %s", variable.Key, want.Key)
		}

		if variable.Value != want.Value {
			return fmt.Errorf("got value %s; value %s", variable.Value, want.Value)
		}

		if variable.Protected != want.Protected {
			return fmt.Errorf("got protected %t; want %t", variable.Protected, want.Protected)
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

func testAccGitlabProjectPagesDomainConfig(rString string) string {
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
  value = "value-%s"
  variable_type = "env_var"
}
	`, rString, rString, rString)
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

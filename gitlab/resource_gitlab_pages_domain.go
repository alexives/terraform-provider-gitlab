package gitlab

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectPagesDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectPagesDomainCreate,
		Read:   resourceGitlabProjectPagesDomainRead,
		Update: resourceGitlabProjectPagesDomainUpdate,
		Delete: resourceGitlabProjectPagesDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			// Required changes pending https://github.com/xanzy/go-gitlab/pull/769
			// "auto_ssl_enabled": {
			// 	Type:     schema.TypeBool,
			// 	Required: false,
			// 	Default:  true,
			// },
			"certificate": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: false,
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: false,
			},
		},
	}
}

func resourceGitlabProjectPagesDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	domain := d.Get("domain").(string)
	// autoSslEnabled := d.Get("auto_ssl_enabled").(string)
	certificate := d.Get("certificate").(string)
	key := d.Get("key").(string)

	options := gitlab.CreatePagesDomainOptions{
		Domain:      &domain,
		Key:         &key,
		Certificate: &certificate,
	}
	log.Printf("[DEBUG] create gitlab project variable %s/%s", project, key)

	_, _, err := client.PagesDomains.CreatePagesDomain(project, &options)
	if err != nil {
		return err
	}

	d.SetId(buildTwoPartID(&project, &domain))

	return resourceGitlabProjectPagesDomainRead(d, meta)
}

func resourceGitlabProjectPagesDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project, domain, err := parseTwoPartID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab project variable %s/%s", project, domain)

	v, _, err := client.PagesDomains.GetPagesDomain(project, domain)
	if err != nil {
		return err
	}

	d.Set("project", project)
	d.Set("domain", v.Domain)
	d.Set("url", v.URL)
	d.Set("verified", v.Verified)
	d.Set("verification_code", v.VerificationCode)
	// d.Set("auto_ssl_enabled", v.AutoSslEnabled)
	d.Set("protected", v.VerificationCode)
	d.Set("certificate", v.Certificate)

	return nil
}

func resourceGitlabProjectPagesDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	domain := d.Get("domain").(string)
	// autoSslEnabled := d.Get("auto_ssl_enabled").(string)
	certificate := d.Get("certificate").(string)
	key := d.Get("key").(string)

	options := &gitlab.UpdatePagesDomainOptions{
		// This  is spelled incorrectly in the go-gitlab client :(
		Cerificate: &certificate,
		Key:        &key,
		// AutoSslEnabled: &autoSslEnabled,
	}
	log.Printf("[DEBUG] update gitlab project variable %s/%s", project, domain)

	_, _, err := client.PagesDomains.UpdatePagesDomain(project, domain, options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectPagesDomainRead(d, meta)
}

func resourceGitlabProjectPagesDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	domain := d.Get("domain").(string)
	log.Printf("[DEBUG] Delete gitlab project variable %s/%s", project, domain)

	_, err := client.PagesDomains.DeletePagesDomain(project, domain)
	return err
}

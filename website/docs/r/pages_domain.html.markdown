---
layout: "gitlab"
page_title: "GitLab: gitlab_pages_domain"
sidebar_current: "docs-gitlab-resource-pages_domain"
description: |-
  Creates and manages pages domains for gitlab projects
---

# gitlab\_pages\_domain

This resource allows you to create and manage pages domains for your GitLab projects.
For further information on pages, consult the [GitLab Documentation](https://docs.gitlab.com/ee/user/project/pages/index.html)

This resource allows you to create and manage project clusters for your GitLab projects.
For further information on clusters, consult the [gitlab
documentation](https://docs.gitlab.com/ce/user/project/clusters/index.html).


## Example Usage

```hcl
resource "gitlab_project" "foo" {
  name = "foo-project"
}

resource "acme_certificate" "certificate" {
  ... your acme certificate provider
}

resource gitlab_pages_domain "bar" {
  project     = "${gitlab_project.foo.id}"
  domain      = "pages.example.com"
  key         = "${acme_certificate.private_key_pem}"
  certificate = "${acme_certificate.certificate_pem}"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The id of the project to add the cluster to.

* `domain` - (Required, string) The name of cluster.
<!-- Required changes pending https://github.com/xanzy/go-gitlab/pull/769
* `auto_ssl_enabled` - (Optiona, boolean) If you would like gitlab to automatically generate lets encrypt certificates.
 -->
* `key` - (Optional, string) The private key used for the SSL certificate.

* `certificate` - (Optional, string) The certificate used for SSL validation.

## Import

GitLab project clusters can be imported using an id made up of `projectid:domain`, e.g.

```
$ terraform import gitlab_project_cluster.bar 123:pages.example.com
```

# terraform_module_pinned_source

Disallow specifying a git or mercurial repository as a module source without pinning to a non-default version.

## Example

```hcl
module "unpinned" {
  source = "git://hashicorp.com/consul.git"
}

module "default git" {
  source = "git://hashicorp.com/consul.git?ref=master"
}

module "default mercurial" {
  source = "hg::http://hashicorp.com/consul.hg?rev=default"
}
```

```
$ tflint
template.tf
        WARNING:2 Module source "git://hashicorp.com/consul.git" is not pinned (terraform_module_pinned_source)
        WARNING:6 Module source "git://hashicorp.com/consul.git?ref=master" uses default ref "master" (terraform_module_pinned_source)
        WARNING:10 Module source "hg::http://hashicorp.com/consul.hg?rev=default" uses default rev "default" (terraform_module_pinned_source)

Result: 3 issues  (0 errors , 3 warnings , 0 notices)
```

## Why

Terraform allows you to checkout module definitions from source control. If you do not pin the version to checkout, the dependency you require may introduce major breaking changes without your awareness. To prevent this, always specify an explicit version to checkout.

## How To Fix

Specify a version pin.  For git repositories, it should not be "master". For Mercurial repositories, it should not be "default"

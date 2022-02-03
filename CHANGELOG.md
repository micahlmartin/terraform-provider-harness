# 0.1.13 (Unreleased)

ENHANCEMENTS:

* Made a bunch of interesting changes you should look at ([#2](https://github.com/hashicorp/terraform-provider-aws/issues/2))
* Some big changes here are happening. ([#1](https://github.com/hashicorp/terraform-provider-aws/issues/1))
* resource/aws_lb: Update `name` and `name_prefix` plan-time validation to exclude `"internal-"` ([#80](https://github.com/hashicorp/terraform-provider-aws/issues/80))

BUG FIXES:

* Fixed missing nextgen auth configuration in the provider. ([#76](https://github.com/hashicorp/terraform-provider-aws/issues/76))
* Fixes some bugs and other things. ([#79](https://github.com/hashicorp/terraform-provider-aws/issues/79))

# 0.1.12

BUG FIXES:

* Fixes some bugs. ([#74](https://github.com/hashicorp/terraform-provider-aws/issues/74))

# 0.1.11

BUG FIXES:

* Upgraded harness-go-sdk to v0.1.11 to fix authentication configuration bug [71](https://github.com/harness-io/terraform-provider-harness/issues/71)

## 0.1.10

BREAKING CHANGES:

* Anyone using either on-prem or compliance clusters (anything other than `app.harness.io`) will need to update their provider endpoint configuration. The `/api` should be dropped from the configuration since this is now implied. For example, if you were previously setting this:

```terraform
provider "harness" {
  endpoint = "https://my.domain.com/api"
}
```

It should now be:
```terraform
provider "harness" {
  endpoint = "https://my.domain.com"
}
```

ENHANCEMENTS:

* Upgraded harness-go-sdk to v0.1.10
* Richer debug logging support added
* Refactored and simplified client configuration setup

## 0.1.9

ENHANCEMENTS:

* Upgraded harness-go-sdk to v0.
* **New Resource:** `harness_delegate_ids`
* data-source/delegate: Add support for looking up delegates by `hostname`

BUG FIXES:

* resource/harness_environment: Variable override field `service_name` is now optional. This allows a variable override to apply to all services when being deployed to an environment.
* Fixes delegate not found panic [#64](https://github.com/harness-io/terraform-provider-harness/issues/64)

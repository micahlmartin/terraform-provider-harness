package provider

import (
	"context"
	"errors"

	"github.com/harness-io/harness-go-sdk/harness/api"
	"github.com/harness-io/harness-go-sdk/harness/api/cac"
	"github.com/harness-io/terraform-provider-harness/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudProviderGcp() *schema.Resource {

	providerSchema := map[string]*schema.Schema{
		"skip_validation": {
			Description:   "Skip validation of GCP configuration.",
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"secret_file_id"},
		},
		"delegate_selectors": {
			Description: "Delegate selectors to use for this provider.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ConflictsWith: []string{"secret_file_id"},
		},
		"secret_file_id": {
			Description:   "The id of the secret containing the GCP credentials",
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"delegate_selectors", "usage_scope"},
		},
	}

	helpers.MergeSchemas(commonCloudProviderSchema(), providerSchema)

	return &schema.Resource{
		Description:   "Resource for creating a GCP cloud provider",
		CreateContext: resourceCloudProviderGcpCreate,
		ReadContext:   resourceCloudProviderGcpRead,
		UpdateContext: resourceCloudProviderGcpUpdate,
		DeleteContext: resourceCloudProviderGcpDelete,

		Schema: providerSchema,
	}
}

func resourceCloudProviderGcpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	name := d.Get("name").(string)

	cp := &cac.GcpCloudProvider{}
	err := c.ConfigAsCode().GetCloudProviderByName(name, cp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cp.Id)
	d.Set("name", cp.Name)
	d.Set("skip_validation", cp.SkipValidation)
	d.Set("delegate_selectors", cp.DelegateSelectors)

	if cp.ServiceAccountKeyFileContent != nil {
		d.Set("secret_file_id", cp.ServiceAccountKeyFileContent.Name)
	}

	scope, err := flattenUsageRestrictions(c, cp.UsageRestrictions)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("usage_scope", scope)

	return nil
}

func resourceCloudProviderGcpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	input := cac.NewEntity(cac.ObjectTypes.GcpCloudProvider).(*cac.GcpCloudProvider)
	input.Name = d.Get("name").(string)
	input.SkipValidation = d.Get("skip_validation").(bool)

	if selectors := d.Get("delegate_selectors").([]interface{}); len(selectors) > 0 {
		input.UseDelegateSelectors = true
		input.DelegateSelectors = expandDelegateSelectors(selectors)
	}

	if secretId := d.Get("secret_file_id").(string); secretId != "" {
		input.ServiceAccountKeyFileContent = &cac.SecretRef{
			Name: secretId,
		}
	}

	restrictions, err := expandUsageRestrictions(c, d.Get("usage_scope").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}
	input.UsageRestrictions = restrictions

	cp, err := c.ConfigAsCode().UpsertGcpCloudProvider(input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cp.Id)

	return nil
}

func resourceCloudProviderGcpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	if d.HasChange("name") {
		return diag.FromErr(errors.New("name is immutable"))
	}

	cp := cac.NewEntity(cac.ObjectTypes.AzureCloudProvider).(*cac.GcpCloudProvider)
	cp.Name = d.Get("name").(string)
	cp.SkipValidation = d.Get("skip_validation").(bool)

	cp.DelegateSelectors = d.Get("delegate_selectors").([]string)
	cp.UseDelegateSelectors = len(cp.DelegateSelectors) > 0

	if secretId := d.Get("secret_file_id").(string); secretId != "" {
		cp.ServiceAccountKeyFileContent = &cac.SecretRef{
			Name: secretId,
		}
	}

	usageRestrictions, err := expandUsageRestrictions(c, d.Get("usage_scope").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}
	cp.UsageRestrictions = usageRestrictions

	_, err = c.ConfigAsCode().UpsertGcpCloudProvider(cp)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCloudProviderGcpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)

	id := d.Get("id").(string)
	err := c.CloudProviders().DeleteCloudProvider(id)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

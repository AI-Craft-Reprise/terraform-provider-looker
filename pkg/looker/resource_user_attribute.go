package looker

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUserAttribute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserAttributeCreate,
		ReadContext:   resourceUserAttributeRead,
		UpdateContext: resourceUserAttributeUpdate,
		DeleteContext: resourceUserAttributeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value_is_hidden": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_can_view": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_can_edit": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceUserAttributeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)
	userAttributeName := d.Get("name").(string)
	userAttributeLabel := d.Get("label").(string)
	userAttributeType := d.Get("type").(string)
	userAttributeDefaultValue := d.Get("default_value").(string)
	userAttributeValueIsHidden := d.Get("value_is_hidden").(bool)
	userAttributeUserCanView := d.Get("user_can_view").(bool)
	userAttributeUserCanEdit := d.Get("user_can_edit").(bool)

	writeUserAttribute := apiclient.WriteUserAttribute{
		Name:          userAttributeName,
		Label:         userAttributeLabel,
		Type:          userAttributeType,
		DefaultValue:  &userAttributeDefaultValue,
		ValueIsHidden: &userAttributeValueIsHidden,
		UserCanView:   &userAttributeUserCanView,
		UserCanEdit:   &userAttributeUserCanEdit,
	}

	userAttribute, err := client.CreateUserAttribute(writeUserAttribute, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	userAttributeID := *userAttribute.Id
	d.SetId(userAttributeID)

	return resourceUserAttributeRead(ctx, d, m)
}

func resourceUserAttributeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	userAttributeID := d.Id()

	userAttribute, err := client.UserAttribute(userAttributeID, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", userAttribute.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", userAttribute.Type); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("label", userAttribute.Label); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("default_value", userAttribute.DefaultValue); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value_is_hidden", userAttribute.ValueIsHidden); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_can_view", userAttribute.UserCanView); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_can_edit", userAttribute.UserCanEdit); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserAttributeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	userAttributeID := d.Id()

	userAttributeName := d.Get("name").(string)
	userAttributeType := d.Get("type").(string)
	userAttributeLabel := d.Get("label").(string)
	userAttributeDefaultValue := d.Get("default_value").(string)
	userAttributeValueIsHidden := d.Get("value_is_hidden").(bool)
	userAttributeUserCanView := d.Get("user_can_view").(bool)
	userAttributeUserCanEdit := d.Get("user_can_edit").(bool)

	writeUserAttribute := apiclient.WriteUserAttribute{
		Name:          userAttributeName,
		Label:         userAttributeLabel,
		Type:          userAttributeType,
		DefaultValue:  &userAttributeDefaultValue,
		ValueIsHidden: &userAttributeValueIsHidden,
		UserCanView:   &userAttributeUserCanView,
		UserCanEdit:   &userAttributeUserCanEdit,
	}

	_, err := client.UpdateUserAttribute(userAttributeID, writeUserAttribute, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserAttributeRead(ctx, d, m)
}

func resourceUserAttributeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	userAttributeID := d.Id()

	_, err := client.DeleteUserAttribute(userAttributeID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

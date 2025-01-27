package looker

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUserAttributeGroupValue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserAttributeGroupValueCreate,
		ReadContext:   resourceUserAttributeGroupValueRead,
		UpdateContext: resourceUserAttributeGroupValueUpdate,
		DeleteContext: resourceUserAttributeGroupValueDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_attribute_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserAttributeGroupValueCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	groupID := d.Get("group_id").(string)
	userAttributeID := d.Get("user_attribute_id").(string)
	value := d.Get("value").(string)

	body := apiclient.UserAttributeGroupValue{
		GroupId:         &groupID,
		UserAttributeId: &userAttributeID,
		Value:           &value,
	}
	userAttributeGroupValue, err := client.UpdateUserAttributeGroupValue(groupID, userAttributeID, body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	id := buildTwoPartID(userAttributeGroupValue.GroupId, userAttributeGroupValue.UserAttributeId)

	d.SetId(id)

	return resourceUserAttributeGroupValueRead(ctx, d, m)
}

func resourceUserAttributeGroupValueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	groupIDString, userAttributeIDString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	userAttributeGroupValues, err := client.AllUserAttributeGroupValues(userAttributeIDString, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var userAttributeGroupValue apiclient.UserAttributeGroupValue
	for _, groupValue := range userAttributeGroupValues {
		if *groupValue.GroupId == groupIDString {
			userAttributeGroupValue = groupValue
			break
		}
	}

	if err = d.Set("group_id", userAttributeGroupValue.GroupId); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_attribute_id", userAttributeGroupValue.UserAttributeId); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value", userAttributeGroupValue.Value); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserAttributeGroupValueUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	groupIDString, userAttributeIDString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	value := d.Get("value").(string)

	body := apiclient.UserAttributeGroupValue{
		GroupId:         &groupIDString,
		UserAttributeId: &userAttributeIDString,
		Value:           &value,
	}
	_, err = client.UpdateUserAttributeGroupValue(groupIDString, userAttributeIDString, body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserAttributeGroupValueRead(ctx, d, m)
}

func resourceUserAttributeGroupValueDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	groupIDString, userAttributeIDString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteUserAttributeGroupValue(groupIDString, userAttributeIDString, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

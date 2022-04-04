package looker

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permission_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"model_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	roleName := d.Get("name").(string)
	permissionSetID := d.Get("permission_set_id").(string)
	modelSetID := d.Get("model_set_id").(string)
	pSetID := permissionSetID

	mSetID := modelSetID

	writeRole := apiclient.WriteRole{
		Name:            &roleName,
		PermissionSetId: &pSetID,
		ModelSetId:      &mSetID,
	}

	role, err := client.CreateRole(writeRole, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := *role.Id
	d.SetId(roleID)

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	roleID := d.Id()

	role, err := client.Role(roleID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", role.Name); err != nil {
		return diag.FromErr(err)
	}
	pSetID := *role.PermissionSet.Id
	if err = d.Set("permission_set_id", pSetID); err != nil {
		return diag.FromErr(err)
	}
	mSetID := *role.ModelSet.Id
	if err = d.Set("model_set_id", mSetID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	roleID := d.Id()
	roleName := d.Get("name").(string)
	permissionSetID := d.Get("permission_set_id").(string)
	modelSetID := d.Get("model_set_id").(string)
	pSetID := permissionSetID

	mSetID := modelSetID

	writeRole := apiclient.WriteRole{
		Name:            &roleName,
		PermissionSetId: &pSetID,
		ModelSetId:      &mSetID,
	}
	_, err := client.UpdateRole(roleID, writeRole, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)
	roleID := d.Id()

	_, err := client.DeleteRole(roleID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

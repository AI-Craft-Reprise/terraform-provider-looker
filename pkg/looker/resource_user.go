package looker

import (
	"context"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)

	writeUser := apiclient.WriteUser{
		FirstName: &firstName,
		LastName:  &lastName,
	}

	// CreateUser sometimes returns 500 error
	var user apiclient.User
	err := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
		var err error

		user, err = client.CreateUser(writeUser, "", nil)
		if err != nil {
			if d.IsNewResource() && strings.Contains(err.Error(), "500") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	userID := *user.Id
	d.SetId(userID)

	writeCredentialsEmail := apiclient.WriteCredentialsEmail{
		Email: &email,
	}

	_, err = client.CreateUserCredentialsEmail(userID, writeCredentialsEmail, "", nil)
	if err != nil {
		if _, err = client.DeleteUser(userID, nil); err != nil {
			return diag.FromErr(err)
		}
		return diag.FromErr(err)
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	userID := d.Id()

	user, err := client.User(userID, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("first_name", user.FirstName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_name", user.LastName); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	userID := d.Id()

	if d.HasChanges("first_name", "last_name") {
		firstName := d.Get("first_name").(string)
		lastName := d.Get("last_name").(string)
		writeUser := apiclient.WriteUser{
			FirstName: &firstName,
			LastName:  &lastName,
		}
		_, err := client.UpdateUser(userID, writeUser, "", nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("email") {
		email := d.Get("email").(string)
		writeCredentialsEmail := apiclient.WriteCredentialsEmail{
			Email: &email,
		}
		_, err := client.UpdateUserCredentialsEmail(userID, writeCredentialsEmail, "", nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)
	userID := d.Id()

	_, err := client.DeleteUser(userID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

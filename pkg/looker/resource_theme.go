package looker

import (
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func resourceTheme() *schema.Resource {
	return &schema.Resource{
		Create: resourceThemeCreate,
		Read:   resourceThemeRead,
		Update: resourceThemeUpdate,
		Delete: resourceThemeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceThemeImport,
		},
		Description: "Futher documentation can be found here: https://docs.looker.com/reference/api-and-integration/api-reference/v3.1/theme",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of theme. Can only be alphanumeric and underscores",
				Required:    true,
			},
			"default": {
				Type:        schema.TypeBool,
				Description: "Only an active theme with no expiration (end_at not set) can be assigned as the default theme. As long as a theme has an active record with no expiration, it can be set as the default.",
				Optional:    true,
				Default:     false,
			},
			"background_color": {
				Type:        schema.TypeString,
				Description: "Default background color",
				Optional:    true,
				Default:     "#f6f8fa",
			},
			"base_font_size": {
				Type:        schema.TypeString,
				Description: "Base font size for scaling fonts",
				Optional:    true,
				Default:     "12px",
			},
			"color_collection_id": {
				Type:        schema.TypeString,
				Description: "Optional. ID of color collection to use with the theme. Use an empty string for none.",
				Optional:    true,
			},
			"font_color": {
				Type:        schema.TypeString,
				Description: "Default font color",
				Optional:    true,
				Default:     "#3e3f40",
			},
			"font_family": {
				Type:        schema.TypeString,
				Description: "Primary font family",
				Optional:    true,
				Default:     "Roboto, 'Noto Sans', 'Noto Sans JP', 'Noto Sans CJK KR', 'Noto Sans Arabic UI', 'Noto Sans Devanagari UI', 'Noto Sans Hebrew', 'Noto Sans Thai UI', Helvetica, Arial, sans-serif",
			},
			"font_source": {
				Type:        schema.TypeString,
				Description: "Source specification for font",
				Optional:    true,
			},
			"info_button_color": {
				Type:        schema.TypeString,
				Description: "Info button color",
				Optional:    true,
				Default:     "#0087e1",
			},
			"primary_button_color": {
				Type:        schema.TypeString,
				Description: "Primary button color",
				Optional:    true,
				Default:     "#1A73E8",
			},
			"show_filters_bar": {
				Type:        schema.TypeBool,
				Description: "Toggle to show filters. Defaults to true",
				Optional:    true,
				Default:     true,
			},
			"show_title": {
				Type:        schema.TypeBool,
				Description: "Togle to show the title. Defaults to true",
				Optional:    true,
				Default:     true,
			},
			"text_tile_text_color": {
				Type:        schema.TypeString,
				Description: "Text color for the text tiles",
				Optional:    true,
			},
			"tile_background_color": {
				Type:        schema.TypeString,
				Description: "Background color for tiles",
				Optional:    true,
				Default:     "#ffffff",
			},
			"tile_text_color": {
				Type:        schema.TypeString,
				Description: "text color for tiles",
				Optional:    true,
				Default:     "#3a4245",
			},
			"title_color": {
				Type:        schema.TypeString,
				Description: "Color for titles",
				Optional:    true,
				Default:     "#3a4245",
			},
			"warn_button_color": {
				Type:        schema.TypeString,
				Description: "Warning button color",
				Optional:    true,
				Default:     "#980c11",
			},
			"tile_title_alignment": {
				Type:        schema.TypeString,
				Description: "The text alignment of tile titles (New Dashboards)",
				Optional:    true,
				Default:     "center",
			},
			"tile_shadow": {
				Type:        schema.TypeBool,
				Description: "Toggles the tile shadow (New Dashboards)",
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceThemeCreate(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)
	themeName := d.Get("name").(string)

	themeDef := apiclient.WriteTheme{
		Name:     &themeName,
		Settings: getThemeSettings(d),
	}

	theme, err := client.CreateTheme(themeDef, nil)
	if err != nil {
		return err
	}

	if isDefaultTheme, ok := d.GetOk("default"); ok && isDefaultTheme.(bool) {
		if _, err := client.SetDefaultTheme(themeName, nil); err != nil {
			return err
		}
	}

	d.SetId(strconv.Itoa(int(*theme.Id)))

	return resourceThemeRead(d, m)
}

func resourceThemeRead(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	theme, err := client.Theme(d.Id(), "", nil)
	if err != nil {
		return err
	}

	if err = d.Set("name", theme.Name); err != nil {
		return err
	}

	if theme.Settings.BackgroundColor != nil {
		if err = d.Set("background_color", *theme.Settings.BackgroundColor); err != nil {
			return err
		}
	}

	if theme.Settings.BaseFontSize != nil {
		if err = d.Set("base_font_size", *theme.Settings.BaseFontSize); err != nil {
			return err
		}
	}

	if theme.Settings.ColorCollectionId != nil {
		if err = d.Set("color_collection_id", *theme.Settings.ColorCollectionId); err != nil {
			return err
		}
	}

	if theme.Settings.FontColor != nil {
		if err = d.Set("font_color", *theme.Settings.FontColor); err != nil {
			return err
		}
	}

	if theme.Settings.FontFamily != nil {
		if err = d.Set("font_family", *theme.Settings.FontFamily); err != nil {
			return err
		}
	}

	if theme.Settings.FontSource != nil {
		if err = d.Set("font_source", *theme.Settings.FontSource); err != nil {
			return err
		}
	}

	if theme.Settings.InfoButtonColor != nil {
		if err = d.Set("info_button_color", *theme.Settings.InfoButtonColor); err != nil {
			return err
		}
	}

	if theme.Settings.PrimaryButtonColor != nil {
		if err = d.Set("primary_button_color", *theme.Settings.PrimaryButtonColor); err != nil {
			return err
		}
	}

	if theme.Settings.ShowFiltersBar != nil {
		if err = d.Set("show_filters_bar", *theme.Settings.ShowFiltersBar); err != nil {
			return err
		}
	}

	if theme.Settings.ShowTitle != nil {
		if err = d.Set("show_title", *theme.Settings.ShowTitle); err != nil {
			return err
		}
	}

	if theme.Settings.TextTileTextColor != nil {
		if err = d.Set("text_tile_text_color", *theme.Settings.TextTileTextColor); err != nil {
			return err
		}
	}

	if theme.Settings.TileBackgroundColor != nil {
		if err = d.Set("tile_background_color", *theme.Settings.TileBackgroundColor); err != nil {
			return err
		}
	}

	if theme.Settings.TileTextColor != nil {
		if err = d.Set("tile_text_color", *theme.Settings.TileTextColor); err != nil {
			return err
		}
	}

	if theme.Settings.TitleColor != nil {
		if err = d.Set("title_color", *theme.Settings.TitleColor); err != nil {
			return err
		}
	}

	if theme.Settings.WarnButtonColor != nil {
		if err = d.Set("warn_button_color", *theme.Settings.WarnButtonColor); err != nil {
			return err
		}
	}

	if theme.Settings.TileTitleAlignment != nil {
		if err = d.Set("tile_title_alignment", *theme.Settings.TileTitleAlignment); err != nil {
			return err
		}
	}

	if theme.Settings.TileShadow != nil {
		if err = d.Set("tile_shadow", *theme.Settings.TileShadow); err != nil {
			return err
		}
	}

	defaultTheme, err := client.DefaultTheme(time.Now(), nil)
	if err != nil {
		return err
	}

	themeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	if err = d.Set("default", int(*defaultTheme.Id) == themeID); err != nil {
		return err
	}

	return nil
}

func resourceThemeUpdate(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	themeName := d.Get("name").(string)

	writeRole := apiclient.WriteTheme{
		Name:     &themeName,
		Settings: getThemeSettings(d),
	}

	_, err := client.UpdateTheme(d.Id(), writeRole, nil)
	if err != nil {
		return err
	}

	return resourceThemeRead(d, m)
}

func resourceThemeDelete(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	_, err := client.DeleteTheme(d.Id(), nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceThemeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceThemeRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func getThemeSettings(d *schema.ResourceData) *apiclient.ThemeSettings {
	themeSettings := &apiclient.ThemeSettings{}

	if backgroundColor, ok := d.GetOk("background_color"); ok {
		bColor := backgroundColor.(string)
		themeSettings.BackgroundColor = &bColor
	}

	if baseFontSize, ok := d.GetOk("base_font_size"); ok {
		bFontSize := baseFontSize.(string)
		themeSettings.BaseFontSize = &bFontSize
	}

	if colorCollectionID, ok := d.GetOk("color_collection_id"); ok {
		setting := colorCollectionID.(string)
		themeSettings.ColorCollectionId = &setting
	}

	if fontColor, ok := d.GetOk("font_color"); ok {
		setting := fontColor.(string)
		themeSettings.FontColor = &setting
	}

	if fontFamily, ok := d.GetOk("font_family"); ok {
		setting := fontFamily.(string)
		themeSettings.FontFamily = &setting
	}

	if fontSource, ok := d.GetOk("font_source"); ok {
		setting := fontSource.(string)
		themeSettings.FontSource = &setting
	}

	if infoButtonColor, ok := d.GetOk("info_button_color"); ok {
		setting := infoButtonColor.(string)
		themeSettings.InfoButtonColor = &setting
	}

	if primaryButtonColor, ok := d.GetOk("primary_button_color"); ok {
		setting := primaryButtonColor.(string)
		themeSettings.PrimaryButtonColor = &setting
	}

	if textTileTextColor, ok := d.GetOk("text_tile_text_color"); ok {
		setting := textTileTextColor.(string)
		themeSettings.TextTileTextColor = &setting
	}

	if tileBackgroundColor, ok := d.GetOk("tile_background_color"); ok {
		setting := tileBackgroundColor.(string)
		themeSettings.TileBackgroundColor = &setting
	}

	if tileTextColor, ok := d.GetOk("tile_text_color"); ok {
		setting := tileTextColor.(string)
		themeSettings.TileTextColor = &setting
	}

	if titleColor, ok := d.GetOk("title_color"); ok {
		setting := titleColor.(string)
		themeSettings.TitleColor = &setting
	}

	if warnButtonColor, ok := d.GetOk("warn_button_color"); ok {
		setting := warnButtonColor.(string)
		themeSettings.WarnButtonColor = &setting
	}

	if tileTitleAlignment, ok := d.GetOk("tile_title_alignment"); ok {
		setting := tileTitleAlignment.(string)
		themeSettings.TileTitleAlignment = &setting
	}

	tileShadow := d.Get("tile_shadow").(bool)
	themeSettings.TileShadow = &tileShadow

	showFiltersBar := d.Get("show_filters_bar").(bool)
	themeSettings.ShowFiltersBar = &showFiltersBar

	showTitle := d.Get("show_title").(bool)
	themeSettings.ShowTitle = &showTitle

	return themeSettings
}

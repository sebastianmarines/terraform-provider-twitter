package provider

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"k8s.io/apimachinery/pkg/util/rand"
)

func TestAccProfileResource(t *testing.T) {
	accName := rand.String(5)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Update profile
			{
				Config: testAccProfileResourceConfig(accName, "", "", ""),
				Check:  resource.TestCheckResourceAttr("twitter_profile.acc", "name", accName),
			},
			// Fail when setting `name` to null
			// Keep same URL
		},
	})
}

func testAccProfileResourceConfig(name string, url string, location string, description string) string {
	tmpl, err := template.New("test").Parse(`
resource "twitter_profile" "acc" {
	{{ if ne .name "" }} name = "{{ .screenName }}" {{ end }}
	{{ if ne .url "" }}url = "{{.url}}" {{ end }}
	{{ if ne .location "" }}location = "{{.location}}" {{ end }}
	{{ if ne .description "" }}description = "{{.description}}" {{ end }}
}`)

	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{
		"screenName":  name,
		"url":         url,
		"location":    location,
		"description": description,
	})

	if err != nil {
		panic(err)
	}

	return buf.String()
}

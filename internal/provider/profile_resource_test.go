package provider

import (
	"bytes"
	"regexp"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"k8s.io/apimachinery/pkg/util/rand"
)

func TestAccProfileResource(t *testing.T) {
	accName := rand.String(5)
	desc := rand.String(20)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Update profile
			{
				Config: testAccProfileResourceConfig(accName, "", "Goland", desc),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("twitter_profile.acc", "name", accName),
					resource.TestCheckResourceAttr("twitter_profile.acc", "url", ""),
					resource.TestCheckResourceAttr("twitter_profile.acc", "location", "Goland"),
					resource.TestCheckResourceAttr("twitter_profile.acc", "description", desc),
				),
			},
		},
	})
}

func TestAccProfileResourceValidators(t *testing.T) {
	accName := rand.String(5)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Valid URL
			{
				Config:      testAccProfileResourceConfig(accName, "invalid url", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("The URL is not a valid URL"),
			},
			// Non blank name
			{
				Config:      testAccProfileResourceConfig("", "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Name can't be blank"),
			},
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

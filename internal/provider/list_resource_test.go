package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccListResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccListResourceConfig("Terraform Provider", "public", "A terraform list"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("twitter_list.acc", "name", "Terraform Provider"),
					resource.TestCheckResourceAttr("twitter_list.acc", "mode", "public"),
					resource.TestCheckResourceAttr("twitter_list.acc", "description", "A terraform list"),
					resource.TestCheckResourceAttr("twitter_list.acc", "member_count", "0"),
				),
				Destroy: false,
			},
			{
				Config: testAccListResourceConfig("Terraform modifier", "public", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("twitter_list.acc", "name", "Terraform modifier"),
					resource.TestCheckResourceAttr("twitter_list.acc", "mode", "public"),
					resource.TestCheckResourceAttr("twitter_list.acc", "description", ""),
				),
			},
		},
	})
}

func testAccListResourceConfig(name string, mode string, description string) string {
	if name == "" {
		name = "null"
	} else {
		name = "\"" + name + "\""
	}

	if mode == "" {
		mode = "null"
	} else {
		mode = "\"" + mode + "\""
	}

	if description == "" {
		description = "null"
	} else {
		description = "\"" + description + "\""
	}

	return fmt.Sprintf(`
resource "twitter_list" "acc" {
  name = %[1]s
  mode = %[2]s
  description = %[3]s
}
`, name, mode, description)
}

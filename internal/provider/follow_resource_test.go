package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFollowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Follow a user
			{
				Config: testAccFollowResourceConfig("HashiCorp", -1),
				Check:  resource.TestCheckResourceAttr("twitter_follow.acc", "user_id", "290900886"),
			},
			// Test that following a private user fails
			{
				Config:      testAccFollowResourceConfig("Terraformpriva1", -1),
				ExpectError: regexp.MustCompile("Following private users is not supported"),
			},
		},
	})
}

func testAccFollowResourceConfig(screenName string, userId int64) string {
	var userIdString string
	if userId == -1 {
		userIdString = "null"
	} else {
		userIdString = strconv.FormatInt(userId, 10)
	}

	if screenName == "" {
		screenName = "null"
	}

	return fmt.Sprintf(`
resource "twitter_follow" "acc" {
  screen_name = %[1]q
}
`, screenName, userIdString)
}

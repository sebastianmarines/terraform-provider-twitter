package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	var userId int64 = 290900886
	screenName := "HashiCorp"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(userId, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.twitter_user.acc", "id", strconv.FormatInt(userId, 10)),
					resource.TestCheckResourceAttr("data.twitter_user.acc", "screen_name", screenName),
					resource.TestCheckResourceAttrSet("data.twitter_user.acc", "followers_count"),
					resource.TestCheckResourceAttrSet("data.twitter_user.acc", "profile_image_url"),
				),
				Destroy: true,
			},
			{
				Config: testAccUserDataSourceConfig(-1, screenName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.twitter_user.acc", "id", strconv.FormatInt(userId, 10)),
					resource.TestCheckResourceAttr("data.twitter_user.acc", "screen_name", screenName),
					resource.TestCheckResourceAttrSet("data.twitter_user.acc", "followers_count"),
					resource.TestCheckResourceAttrSet("data.twitter_user.acc", "profile_image_url"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig(id int64, screenName string) string {
	var userIdString string
	if id == -1 {
		userIdString = "null"
	} else {
		userIdString = strconv.FormatInt(id, 10)
	}

	if screenName == "" {
		screenName = "null"
	} else {
		screenName = "\"" + screenName + "\""
	}

	return fmt.Sprintf(`
data "twitter_user" "acc" {
  id          = %[1]s
  screen_name = %[2]s
}`, userIdString, screenName)

}

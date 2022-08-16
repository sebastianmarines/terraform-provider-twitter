package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTweetDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTweetDataSourceConfig(1559537820804399104),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.twitter_tweet.acc", "text", "Hello World!"),
				),
			},
			{
				Config:      testAccTweetDataSourceConfig(1),
				ExpectError: regexp.MustCompile("Unable to read tweet, got error: .+"),
			},
		},
	})
}

func testAccTweetDataSourceConfig(id int64) string {
	return fmt.Sprintf(`
data "twitter_tweet" "acc" {
  id = %[1]q
}`, strconv.FormatInt(id, 10))

}

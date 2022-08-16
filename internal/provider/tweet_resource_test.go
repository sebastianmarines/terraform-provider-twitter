package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"k8s.io/apimachinery/pkg/util/rand"
)

func TestAccTweetResource(t *testing.T) {
	tweetText := rand.String(5)
	modifiedText := rand.String(5)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTweetResourceConfig(tweetText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("twitter_tweet.acc", "text", tweetText),
					resource.TestCheckResourceAttrSet("twitter_tweet.acc", "id"),
					resource.TestCheckResourceAttrSet("twitter_tweet.acc", "user_id"),
					resource.TestCheckResourceAttrSet("twitter_tweet.acc", "favorite_count"),
				),
				Destroy: false,
			},
			{
				Config: testAccTweetResourceConfig(modifiedText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("twitter_tweet.acc", "text", modifiedText),
					resource.TestCheckResourceAttrSet("twitter_tweet.acc", "id"),
					resource.TestCheckResourceAttrSet("twitter_tweet.acc", "user_id"),
					resource.TestCheckResourceAttrSet("twitter_tweet.acc", "favorite_count"),
				),
			},
		},
	})
}

func testAccTweetResourceConfig(text string) string {
	return fmt.Sprintf(`
resource "twitter_tweet" "acc" {
  text = %[1]q
}`, text)

}

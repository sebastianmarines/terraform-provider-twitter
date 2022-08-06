package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"twitter": providerserver.NewProtocol6WithError(New("twitter")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TWITTER_API_KEY"); v == "" {
		t.Error("Missing Twitter API key")
	}
	if v := os.Getenv("TWITTER_API_SECRET_KEY"); v == "" {
		t.Error("Missing Twitter API secret key")
	}
	if v := os.Getenv("TWITTER_ACCESS_TOKEN"); v == "" {
		t.Error("Missing Twitter access token")
	}
	if v := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"); v == "" {
		t.Error("Missing Twitter access secret")
	}
}

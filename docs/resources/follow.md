---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "twitter_follow Resource - terraform-provider-twitter"
subcategory: ""
description: |-
  Allows the authenticating user to follow (friend) a user.
---

# twitter_follow (Resource)

Allows the authenticating user to follow (friend) a user.

## Example Usage

```terraform
resource "twitter_follow" "test" {
  screen_name = "HashiCorp"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `screen_name` (String) The screen name of the user being followed.
- `user_id` (Number) The ID of the user being followed.

### Read-Only

- `id` (Number) The ID of the user being followed.



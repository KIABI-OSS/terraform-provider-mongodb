---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mongodb Provider"
subcategory: ""
description: |-
  Create resources in MongoDB.
---

# mongodb Provider

Create resources in MongoDB.

## Example Usage

```terraform
provider "mongodb" {
  # the env variable MONGODB_URL can be used instead
  url = "mongodb://localhost:27017"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `url` (String) URL of the MongoDB instance to connect to.

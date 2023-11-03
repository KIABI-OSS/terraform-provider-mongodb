resource "mongodb_index" "example" {
  database   = "test"
  collection = "test"
  name       = "example"
  keys = [
    {
      "field" : "f1"
      "type" : "asc"
    }
  ]
}

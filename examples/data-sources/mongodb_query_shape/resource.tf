data "mongodb_query_shape" "example" {
  database   = "my_database"
  collection = "my_collection"
  command    = "find"
  filter     = jsonencode({ "status" : "active", "category" : "electronics" })
  sort       = jsonencode({ "created_at" : -1 })
  projection = jsonencode({ "_id" : 0, "name" : 1, "price" : 1 })
}
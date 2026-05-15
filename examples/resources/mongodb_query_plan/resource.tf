resource "mongodb_query_plan" "example" {
  query_hash = "0000000000000000000000000000000000000000000000000000000000000000"
  database   = "my_database"
  collection = "my_collection"
  allowed_indexes = [
    "idx_1",
    "idx_2"
  ]
  comment = "Force index usage for this query shape"
}

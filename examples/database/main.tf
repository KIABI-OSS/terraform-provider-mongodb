terraform {
  required_providers {
    mongodb = {
      source = "hashicorp.com/edu/mongodb"
    }
  }
}

provider "mongodb" {
  # the env variable MONGODB_URL can be used instead
  url = "mongodb://localhost:27017"
}


resource "mongodb_database" "db" {
  name = "some-database-name"
}

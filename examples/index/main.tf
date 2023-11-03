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

resource "mongodb_index" "test_single_field" {
  database   = "test"
  collection = "test"
  name       = "single_field"
  keys = [
    {
      "field" : "f1"
      "type" : "asc"
    }
  ]
}

resource "mongodb_index" "test_compound" {
  database   = "test"
  collection = "test"
  name       = "compound"
  keys = [
    {
      "field" : "f1"
      "type" : "asc"
    },
    {
      "field" : "f2"
      "type" : "desc"
    }
  ]
}

resource "mongodb_index" "test_wildcard" {
  database   = "test"
  collection = "test"
  name       = "wildcard"
  keys = [
    {
      "field" : "$**",
      "type" : "asc"
    }
  ]
  wildcard_projection = {
    "some_field" : 1
  }
}

resource "mongodb_index" "test_2dsphere" {
  database   = "test"
  collection = "test"
  name       = "2dsphere"
  keys = [
    {
      "field" : "sphere_location",
      "type" : "2dsphere"
    }
  ]
}

resource "mongodb_index" "test_2d" {
  database   = "test"
  collection = "test"
  name       = "2d"
  keys = [
    {
      "field" : "plan_location",
      "type" : "2d"
    }
  ]
}

resource "mongodb_index" "test_hashed" {
  database   = "test"
  collection = "test"
  name       = "hashed"
  keys = [
    {
      "field" : "f1",
      "type" : "hashed"
    }
  ]
}

resource "mongodb_index" "test_ttl" {
  database   = "test"
  collection = "test"
  name       = "ttl"
  keys = [
    {
      "field" : "timestamp",
      "type" : "asc"
    }
  ]
  expire_after_seconds = 3600
}

resource "mongodb_index" "test_sparse" {
  database   = "test"
  collection = "test"
  name       = "sparse"
  keys = [
    {
      "field" : "f1",
      "type" : "asc"
    }
  ]
  sparse = true
}

resource "mongodb_index" "test_unique" {
  database   = "test"
  collection = "test"
  name       = "unique"
  keys = [
    {
      "field" : "f1",
      "type" : "asc"
    }
  ]
  unique = true
}

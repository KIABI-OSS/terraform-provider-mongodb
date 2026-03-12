# Terraform Provider MongoDB

This provider can be used to create MongoDB resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install .
```

## Provider configuration

```terraform
provider "mongodb" {
  url = "mongodb://localhost:27017"
}
```

> The environment variable MONGODB_URL can be used instead.

## Available resources

### [Indexes](https://www.mongodb.com/docs/manual/indexes/)

The provider can be used to create indexes in a collection. The supported types of indexes are:

- Single Field Indexes
- Compound Indexes
- Multikey Indexes
- Geospatial Indexes
- Hashed Indexes
- Wildcard Indexes

The created indexes support the following properties

- Sparse Indexes
- TTL Indexes
- Unique
- Collations
- Background

You can find examples [here](examples/index/main.tf)

#### Import

All supported index types can now be imported using `terraform import <resource_path> <index_id>`.
Index id must use the format `<database>.<collection>.<index_name>`.

> This means that index id with database, collection or index containing `.` do NOT work.

## Known issues

### Index import and collation/wildcard projection

Even though index collations and wildcard projections are supported when creating index
they do NOT work with import. The `IndexSpecification` returned by the mongo driver when
fetching index data doesn't contain those fields so they cannot be inserted into the state.

Issue: https://github.com/KIABI-OSS/terraform-provider-mongodb/issues/4

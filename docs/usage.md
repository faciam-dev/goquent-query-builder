# Usage Guide

This document explains how to use **goquent-query-builder** to build SQL queries in Go.

## Installation

Add the module to your project using `go get`:

```shell
go get github.com/faciam-dev/goquent-query-builder
```

Import the packages you need in your code.

## Creating a builder

Create a query builder for your database. Example for MySQL:

```go
import (
    "github.com/faciam-dev/goquent-query-builder/api"
    "github.com/faciam-dev/goquent-query-builder/database/mysql"
)

qb := api.NewSelectQueryBuilder(mysql.NewMySQLQueryBuilder())
```

For PostgreSQL use `postgres.NewPostgreSQLQueryBuilder()` instead.

## Building queries

Builders expose a fluent API. For example, to construct a simple SELECT:

```go
query, values, err := qb.Table("users").
    Select("id", "name").
    Where("age", ">", 18).
    OrderBy("name", "ASC").
    Build()
```

`query` contains the SQL string and `values` provides the bound parameters.

There are also builders for INSERT, UPDATE and DELETE operations:

```go
api.NewInsertQueryBuilder(mysql.NewMySQLQueryBuilder()).
    Table("users").
    Insert(map[string]interface{}{"name": "John"}).
    Build()
```

See the [examples](../example) directory for complete programs.

## Running the examples

Each sub directory in `example` is a standalone Go module. Change into a folder
and execute:

```shell
go run .
```

This prints the generated SQL and parameter values.


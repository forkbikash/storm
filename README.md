# storm - A Golang ORM

⚠️ **Disclaimer: This package is currently under development and should not be used in production environments.**

`storm` is a lightweight and straightforward Object-Relational Mapping (ORM) library for Golang. It aims to provide a simple and intuitive interface for interacting with relational databases while abstracting away the complexities of SQL query construction and data mapping.

## Features to be implemented

- **Database Abstraction**: Storm supports multiple database engines, including PostgreSQL, MySQL, and SQLite.
- **Model Definitions**: Define your data models using Go structs with struct tags for mapping fields to database columns.
- **Query Builder**: A fluent query builder for constructing complex queries without writing raw SQL.
- **Relationships**: Support for modeling relationships between data models (one-to-one, one-to-many, many-to-many).
- **Migrations**: A built-in migration system for managing database schema changes.
- **Transactions**: Support for database transactions to ensure data consistency.
- **Connection Pooling**: Efficient connection pooling for better performance and resource management.

## Installation

```bash
go get github.com/forkbikash/storm
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/forkbikash/storm"
)

type User struct {
    ID        int    `storm:"id,pk"`
    Name      string `storm:"name"`
    Email     string `storm:"email"`
}

func main() {
    // Open a PostgreSQL connection
    db, err := storm.Open("postgres", "user=postgres password=mypassword host=localhost dbname=mydb sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create a new user
    user := &User{
        Name:  "John Doe",
        Email: "john.doe@example.com",
    }
    err = db.Create(user)
    if err != nil {
        log.Fatal(err)
    }

    // Find a user by ID
    foundUser, err := db.FindByID(user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(foundUser)

    // Update a user
    foundUser.Email = "john.doe.updated@example.com"
    err = db.Update(foundUser)
    if err != nil {
        log.Fatal(err)
    }

    // Delete a user
    err = db.Delete(foundUser)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Contributing

Contributions are welcome!

## License

Storm is released under the MIT License.

# Migrations

## Technology

We are using [Goose](https://github.com/pressly/goose)
To handle our DB migrations. Why? Goose is a simple migration tool that allows us to create and apply migrations to our DB.

- You can find more information about Goose [here](https://pressly.github.io/goose/)
## Installation

```bash 
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Create a new migration
To create a specific migration please go
to the DB you need. For example

```bash
/ - migrations
    / - postgres
    / - clickhouse
```

To create a new migration run
```bash
goose create migration_name sql
```

To create a new Go migration run
```bash
goose create migration_name go
```

---
written by @Xusk947 (e.g. Aziz)
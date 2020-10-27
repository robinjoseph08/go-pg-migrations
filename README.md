# go-pg-migrations

[![Version](https://img.shields.io/badge/version-v3.0.0-green.svg)](https://github.com/robinjoseph08/go-pg-migrations/releases)
[![GoDoc](https://godoc.org/github.com/robinjoseph08/go-pg-migrations?status.svg)](http://godoc.org/github.com/robinjoseph08/go-pg-migrations)
[![Build Status](https://travis-ci.org/robinjoseph08/go-pg-migrations.svg?branch=master)](https://travis-ci.org/robinjoseph08/go-pg-migrations)
[![Coverage Status](https://coveralls.io/repos/github/robinjoseph08/go-pg-migrations/badge.svg?branch=master)](https://coveralls.io/github/robinjoseph08/go-pg-migrations?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/robinjoseph08/go-pg-migrations)](https://goreportcard.com/report/github.com/robinjoseph08/go-pg-migrations)

A Go package to help write migrations with [`go-pg/pg`](https://github.com/go-pg/pg).

## Usage

### Installation

Because `go-pg` [now has Go modules
support](https://github.com/go-pg/pg#get-started), `go-pg-migrations` also has
modules support; it currently depends on v10 of `go-pg`. To install it, use the
following command in a project with a `go.mod`:

```sh
$ go get github.com/robinjoseph08/go-pg-migrations/v3
```

If you are not yet using Go modules, you can still use v1 of this package.

### Running

To see how this package is intended to be used, you can look at the [example
directory](/example). All you need to do is have a `main` package (e.g.
`example`); call `migrations.Run` with the directory you want the migration
files to be saved in (which will be the same directory of the main package, e.g.
`example`), an instance of `*pg.DB`, and `os.Args`; and log any potential errors
that could be returned.

Once this has been set up, then you can use the `create`, `migrate`, `rollback`,
`help` commands like so:

```
$ go run example/*.go create create_users_table
Creating example/20180812001528_create_users_table.go...

$ go run example/*.go migrate
Running batch 1 with 1 migration(s)...
Finished running "20180812001528_create_users_table"

$ go run example/*.go rollback
Rolling back batch 1 with 1 migration(s)...
Finished rolling back "20180812001528_create_users_table"

$ go run example/*.go help
Usage:
  go run example/*.go [command]

Commands:
  create   - create a new migration in example with the provided name
  migrate  - run any migrations that haven't been run yet
  rollback - roll back the previous run batch of migrations
  help     - print this help text

Examples:
  go run example/*.go create create_users_table
  go run example/*.go migrate
  go run example/*.go rollback
  go run example/*.go help
```

While this works when you have the Go toolchain installed, there might be a
scenario where you have to run migrations and you don't have the toolchain
available (e.g. in a `scratch` or `alpine` Docker image deployed to production).
In that case, you should compile another binary (in addition to your actual
application) and copy it into the final image. This will include all of your
migrations and allow you to run it by overriding the command when running the
Docker container.

This would look something like this:

```dockerfile
# Dockerfile
FROM golang:1.13.3 as build

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags '-w -s' -o ./bin/serve ./cmd/serve
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags '-w -s' -o ./bin/migrations ./cmd/migrations

FROM alpine:3.8

RUN apk --no-cache add ca-certificates
COPY --from=build /app/bin /bin

CMD ["serve"]
```

```sh
$ docker build -t service:latest .
$ docker run --rm service:latest migrations migrate
```

## Why?

While go-pg has [its own `migrations`
package](https://github.com/go-pg/migrations), it leaves a bit to be desired.
Some additional features that this package supports:

- Complete migration diffing to determine which migrations still need to be run.
  Since `go-pg/migrations` checks the current version of migrations and runs any
  migrations after that, some migrations can be missed, especially when several
  people are working on the same project.
- Timestamp-based prefixes to prevent two people creating a migration with the
  same version on two separate branches. If the current version is 3, and more
  than one person branches off and creates a new migration, all of them will be
  version 4.
- The ability to run migrations in a transaction on a case-by-case basis. Most
  of the time, running migrations within a transaction is desirable, so that if
  it errs out within the "up" function, the whole migration is reverted. But
  since some long-running migrations might have a statement with a relatively
  exclusive lock, you might opt out of running that specific migration within a
  transaction.
- A migration locking mechanism. This is to avoid two people (or an automated
  deployment system) attempting to run migrations at the same time against the
  same database, which could lead to undesired behavior.
- An expected workflow of how this package should be used within a project.
  While `go-pg/migrations` has some recommendations and examples, this package
  takes a more opinionated approach which makes it so you don't have to think
  about it as much, and there's less code for you to write and maintain.
- Batch-level rollbacks. When there are multiple migration files run during the
  same migration invocation, they are all grouped together into a "batch".
  During rollbacks, each batch gets rolled back together. This tends to be more
  desireable since this usually means the application is reverting back to a
  previous release, so the database should be in the state expected for that
  release.

Many of these features and expected behaviors come from using [Knex.js
migrations](https://knexjs.org/#Migrations) in production for many years. This
project is heavily inspired by Knex to provide a robust and safe migration
experience.

`go-pg` is a great and performant project, and hopefully, this makes it a little
better.

## Development

To develop on this project, you'll need to have Postgres running because the
tests depend on it.

If you have it running on your machine because it was installed through your
package manager (like `brew` or `apt-get`), you just need to run the following
to get it set up correctly:

```sh
make setup
```

If you don't have it on your laptop, you can run the following to start it
within Docker:

```sh
make postgres
```

That should start the container and keep it running while you develop. Once
you're done, you can `^C` out of it and it will stop the container.

To run the tests, you should run:

```sh
make test
```

To run the linter, you should run:

```sh
make lint
```

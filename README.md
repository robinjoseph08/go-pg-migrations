# migrations

A Go package to help write migrations with [`go-pg/pg`](https://github.com/go-pg/pg).

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

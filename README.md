# Migrate stats from the Go version

This is a utility that help you migrate from [`5HT2B/heartbeat`](https://github.com/5HT2B/heartbeat)
to [`5HT2B/heartbeat-riir`](https://github.com/5HT2B/heartbeat-riir). This copies over beats history
and statistics and changes devices to the new structure, emitting their ID, name and token.

## Usage

### Prerequisites

For best experience while using this tool, ensure that:

- Your old heartbeat server is stopped, or alternatively, no clients are using it/it is not publicly
  accessible.
- Your `ReJSON` server is ***not*** stopped, so that statistics may be migrated over.
  - The RDB file format is a mess; and there are multiple breaking changes between versions, and
    looking at `libredis` code to port it just raises my blood pressure; so I just issue commands to
    the running server instead. Not to mention, you may have disabled RDB-based persistence
    altogether, so that whole thing would have been an exercise in futility.
- Your new heartbeat server is up and running.
- The PostgreSQL server the new server is connected to already has the relevant schemas and tables
  defined.
- The environment that this tool will be run in has access to
  - Your `ReJSON` server
  - Your new heartbeat server
  - The PostgreSQL server your new heartbeat server is connected to.
  
  Configuring your networking for this is left as an exercise to the reader. There are far too many
  corner cases for me to cover here.

### Compilation

I don't really plan to build this in CI, so you will need a Go toolchain (>= 1.18.0) available, and
a `Makefile` is provided for convenience.

```console
$ pwd
/home/.../heartbeat/contrib/migrate_from_go

$ make
go build ./cmd/migrate_stats.go

$ ./migrate_stats --help
Usage of ./migrate_stats:
  -d, --database-dsn DSN        PostgreSQL database DSN (default "postgresql://postgres@localhost/postgres")
  -r, --redis-address address   Redis (ReJSON) address (default "127.0.0.1:6379")
  -p, --redis-password string   Redis password
  -s, --server-base-url URL     Server base URL (of the *new* server) (default "http://localhost:6060")
  -t, --token secret_key        Token to use for authentication to the server to add new devices. This is the value of the secret_key config variable that you set.
```

### Running

**Do not run this more than once**

New devices are created on every run, and you may end up with duplicates. All SQL migrations are run
in a transaction, so you will never end up in a partially migrated state.

## Troubleshooting

Any errors encountered will immediately trigger an abort and an error message will be printed.
Sometimes these errors may be cryptic (because I haven't actually bothered to write better ones and
just propagate the existing errors as is for the most part), but hopefully the only errors you may
encounter will be related to connectivity, in which case the error message is clear enough.

If you get error messages relating to (un)marshalling of JSON, check that your `ReJSON` server is
sending valid data and that your new server isn't rejecting your requests for adding new devices.

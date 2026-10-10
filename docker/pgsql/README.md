# PostgreSQL

Copy `.env.example` to `.env`, then set a strong `POSTGRES_PASSWORD`. The
`.env` file is gitignored. For an existing database volume, keep
`POSTGRES_USER`, `POSTGRES_DB`, and the password consistent with the database
already initialized there. Changing these variables does not change existing
database users or passwords.

Run `./start.sh` from this directory. The `pgsql_data` Docker volume holds
persistent data.

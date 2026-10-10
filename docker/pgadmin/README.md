# pgAdmin

Copy `.env.example` to `.env`, then set `PGADMIN_DEFAULT_EMAIL` and a strong
`PGADMIN_DEFAULT_PASSWORD`. The `.env` file is gitignored. Keep
`PGADMIN_LISTEN_PORT=8084` because `start.sh` publishes port 8084.

Run `./start.sh` from this directory. The `pgadmin_data` Docker volume holds
persistent data. For an existing volume, changing the default login values in
`.env` does not reset an already created account.

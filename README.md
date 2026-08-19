# Advanced Web Technologies Week 2 Lab Stage 1
- This Lab is based on Measuring a Synchronous API Contract
-- Abner Bobadilla & Christian Hope


## Lab Report

| Format | Link |
|---|---|
| Google Docs (view) | [Open][View report (Google Docs)](https://docs.google.com/document/d/1HETA4iRyaVHTDELxj4mEv5Yx82j1WhMMZvq4V_XUhSk/edit?usp=sharing) |
| PDF (download) | [Download](./lab1/CMPS4191%20•%20LABORATORY%201_AB_CH.pdf) |

## Project structure
 
```
cmd/api/            → application entrypoint and HTTP handlers (package main)
internal/data/       → database models (package data)
internal/validator/  → input validation helpers (package validator)
migrations/          → golang-migrate schema + seed data migrations
```
## Environment
 
The application reads its Postgres connection string from the `-db-dsn` flag. Set it as an environment variable before running anything:
 
```bash
export GATEKEEPER_DB_DSN="postgres://<user>:<password>@localhost:5432/gatekeeper?sslmode=disable"
```
## Database setup
 
Create the database, then apply all migrations:
 
```bash
createdb gatekeeper
migrate -path=./migrations -database="$GATEKEEPER_DB_DSN" up
```
### Note on Postgres version compatibility
 
This project was developed/tested in a **GitHub Codespaces** environment, which provisions **PostgreSQL 16** by default via `apt`. Some of the original migration SQL assumed native support for `uuidv7()` and `uuidv4()`, which are built-in functions in **PostgreSQL 18+** but do not exist in Postgres 16.
 
To keep the migrations runnable in this environment without requiring a manual Postgres upgrade, `migrations/000001_create_init_extensions_and_types.up.sql` (and the matching `.down.sql`) were extended to:
 
- Explicitly enable the `citext` and `pgcrypto` extensions (`CREATE EXTENSION IF NOT EXISTS ...`), rather than assuming they're already available on a fresh database.
- Define `uuidv7()` as a PL/pgSQL function that reproduces Postgres 18's time-ordered UUID generation using `pgcrypto`'s `gen_random_uuid()`.
- Define `uuidv4()` as a thin wrapper around `gen_random_uuid()`, reproducing Postgres 18's random UUID generator.
These additions are purely compatibility shims — they don't change the schema, data model, or any application behaviour. On a Postgres 18+ instance, `CREATE OR REPLACE FUNCTION` simply overrides the native functions with functionally equivalent versions, so the migrations remain safe to run on either version.

## Running the API
 
```bash
go run ./cmd/api -db-dsn="$GATEKEEPER_DB_DSN" -report-delay=0s
```
 
`-report-delay` sets an artificial delay in report generation, used to run the timed experiment described in the lab handout (`0s`, `3s`, `7s`, `12s`).
 
## Testing
 
```bash
curl --silent --show-error   --output response.json   --write-out 'status=%{http_code}\ntime=%{time_total}s\nbytes=%{size_download}\n'   --request POST   --header 'Content-Type: application/json'   --data @request.json   http://localhost:4000/v1/reports
```
 
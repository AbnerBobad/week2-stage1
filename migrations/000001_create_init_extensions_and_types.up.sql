-- Filename: 000001_create_init_extensions_and_types.up.sql

BEGIN;

CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION uuidv7() RETURNS uuid AS $$
DECLARE
  unix_ts_ms bytea;
  uuid_bytes bytea;
BEGIN
  unix_ts_ms := substring(int8send(floor(extract(epoch FROM clock_timestamp()) * 1000)::bigint) FROM 3 FOR 6);
  uuid_bytes := unix_ts_ms || substring(gen_random_uuid()::bytea FROM 7 FOR 10);
  uuid_bytes := set_byte(uuid_bytes, 6, (b'0111' || get_byte(uuid_bytes, 6)::bit(4))::bit(8)::int);
  uuid_bytes := set_byte(uuid_bytes, 8, (b'10' || get_byte(uuid_bytes, 8)::bit(6))::bit(8)::int);
  RETURN encode(uuid_bytes, 'hex')::uuid;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION uuidv4() RETURNS uuid AS $$
BEGIN
  RETURN gen_random_uuid();
END;
$$ LANGUAGE plpgsql;

CREATE TYPE consumer_status AS ENUM ('active', 'suspended', 'terminated');
CREATE TYPE key_status     AS ENUM ('active', 'rotating', 'revoked');
CREATE TYPE job_status     AS ENUM ('queued', 'processing', 'completed', 'failed', 'cancelled');
COMMIT;
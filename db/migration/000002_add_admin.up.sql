CREATE TABLE "admin_type" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "admin_type" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "admin" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "type_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "last_login" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

ALTER TABLE "admin" ADD FOREIGN KEY ("type_id") REFERENCES "admin_type" ("id");

CREATE TRIGGER set_timestamp_admin_type
BEFORE UPDATE ON "admin_type"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_admin
BEFORE UPDATE ON "admin"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
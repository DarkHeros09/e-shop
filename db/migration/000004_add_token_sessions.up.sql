CREATE TABLE "user_session" (
  "id" uuid PRIMARY KEY NOT NULL,
  "user_id" bigint NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expires_at" timestamptz NOT NULL
);

ALTER TABLE "user_session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
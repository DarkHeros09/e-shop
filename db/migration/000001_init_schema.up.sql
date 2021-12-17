CREATE TABLE "user" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "username" varchar NOT NULL,
  "password" text,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "telephone" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz
);

CREATE TABLE "user_address" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "user_id" bigint NOT NULL,
  "address_line1" varchar NOT NULL,
  "city" varchar NOT NULL,
  "country" varchar NOT NULL,
  "telephone" int NOT NULL
);

CREATE TABLE "user_payment" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "user_id" bigint NOT NULL,
  "payment_type" varchar,
  "provider" varchar,
  "account_no" int,
  "expiry" date
);

CREATE TABLE "shopping_session" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "user_id" bigint UNIQUE NOT NULL,
  "total" decimal NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz
);

CREATE TABLE "cart_item" (
  "id" bigserial UNIQUE PRIMARY KEY NOT NULL,
  "session_id" bigint UNIQUE NOT NULL,
  "product_id" bigint UNIQUE NOT NULL,
  "quantity" int NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz
);

CREATE TABLE "payment_detail" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "order_id" bigint NOT NULL,
  "amount" int NOT NULL,
  "provider" varchar,
  "status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz
);

CREATE TABLE "product" (
  "id" bigserial UNIQUE PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "description" text,
  "SKU" varchar,
  "category_id" bigint UNIQUE NOT NULL,
  "inventory_id" bigint UNIQUE NOT NULL,
  "price" decimal NOT NULL,
  "discount_id" bigint UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE TABLE "order_detail" (
  "id" bigserial UNIQUE PRIMARY KEY NOT NULL,
  "user_id" bigint UNIQUE NOT NULL,
  "total" decimal NOT NULL,
  "payment_id" bigint UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz
);

CREATE TABLE "order_items" (
  "id" bigserial UNIQUE PRIMARY KEY NOT NULL,
  "order_id" bigint UNIQUE NOT NULL,
  "product_id" bigint UNIQUE NOT NULL,
  "quantity" int,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz
);

CREATE TABLE "product_category" (
  "id" bigserial UNIQUE PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "description" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE TABLE "product_inventory" (
  "id" bigserial UNIQUE PRIMARY KEY NOT NULL,
  "quantity" int,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE TABLE "discount" (
  "id" bigserial UNIQUE PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "description" text,
  "discount_percent" decimal NOT NULL DEFAULT 0,
  "active" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz,
  "deleted_at" timestamptz
);

ALTER TABLE "user" ADD FOREIGN KEY ("id") REFERENCES "shopping_session" ("user_id");

ALTER TABLE "user_payment" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "user_address" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "cart_item" ADD FOREIGN KEY ("session_id") REFERENCES "shopping_session" ("id");

ALTER TABLE "payment_detail" ADD FOREIGN KEY ("id") REFERENCES "order_detail" ("payment_id");

ALTER TABLE "product" ADD FOREIGN KEY ("id") REFERENCES "order_items" ("product_id");

ALTER TABLE "product" ADD FOREIGN KEY ("id") REFERENCES "cart_item" ("product_id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "order_detail" ("id");

ALTER TABLE "order_detail" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "product" ADD FOREIGN KEY ("category_id") REFERENCES "product_category" ("id");

ALTER TABLE "product_inventory" ADD FOREIGN KEY ("id") REFERENCES "product" ("inventory_id");

ALTER TABLE "product" ADD FOREIGN KEY ("discount_id") REFERENCES "discount" ("id");

CREATE INDEX ON "user" ("username");

CREATE INDEX ON "user" ("first_name");

CREATE INDEX ON "user" ("last_name");

CREATE INDEX ON "user" ("first_name", "last_name");

CREATE INDEX ON "user" ("telephone");

CREATE INDEX ON "user_address" ("user_id");

CREATE INDEX ON "user_payment" ("user_id");

CREATE INDEX ON "shopping_session" ("user_id");

CREATE INDEX ON "cart_item" ("product_id");

CREATE INDEX ON "cart_item" ("session_id");

CREATE INDEX ON "payment_detail" ("order_id");

CREATE INDEX ON "product" ("category_id");

CREATE INDEX ON "product" ("inventory_id");

CREATE INDEX ON "order_items" ("order_id");

CREATE INDEX ON "order_items" ("product_id");

CREATE INDEX ON "product_category" ("name");

CREATE INDEX ON "discount" ("name");

COMMENT ON COLUMN "payment_detail"."amount" IS 'must be positive';

COMMENT ON COLUMN "product"."price" IS 'must be positive';

COMMENT ON COLUMN "order_detail"."total" IS 'must be positive';

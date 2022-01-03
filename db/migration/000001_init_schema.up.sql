CREATE TABLE "user" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "username" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "telephone" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_address" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "user_id" bigint UNIQUE NOT NULL,
  "address_line" varchar NOT NULL,
  "city" varchar NOT NULL,
  "telephone" int NOT NULL
);

CREATE TABLE "user_payment" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "user_id" bigint UNIQUE NOT NULL,
  "payment_type" varchar NOT NULL,
  "provider" varchar NOT NULL,
  "account_no" int NOT NULL,
  "expiry" date NOT NULL
);

CREATE TABLE "shopping_session" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "user_id" bigint UNIQUE NOT NULL,
  "total" decimal NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "cart_item" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "session_id" bigint UNIQUE NOT NULL,
  "product_id" bigint UNIQUE NOT NULL,
  "quantity" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "payment_detail" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "order_id" bigint NOT NULL DEFAULT 0,
  "amount" int NOT NULL,
  "provider" varchar NOT NULL,
  "status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "product" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "name" varchar UNIQUE NOT NULL,
  "description" text NOT NULL,
  "sku" varchar UNIQUE NOT NULL,
  "category_id" bigint UNIQUE NOT NULL,
  "inventory_id" bigint UNIQUE NOT NULL,
  "price" decimal NOT NULL,
  "active" boolean NOT NULL DEFAULT false,
  "discount_id" bigint UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order_detail" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "user_id" bigint UNIQUE NOT NULL,
  "total" decimal NOT NULL,
  "payment_id" bigint UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order_items" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "order_id" bigint UNIQUE NOT NULL,
  "product_id" bigint UNIQUE NOT NULL,
  "quantity" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "product_category" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "name" varchar UNIQUE NOT NULL,
  "description" text NOT NULL,
  "active" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "product_inventory" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "quantity" int NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "discount" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "name" varchar UNIQUE NOT NULL,
  "description" text NOT NULL,
  "discount_percent" decimal NOT NULL DEFAULT 0,
  "active" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "user_address" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "user_payment" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "shopping_session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "cart_item" ADD FOREIGN KEY ("session_id") REFERENCES "shopping_session" ("id");

ALTER TABLE "cart_item" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");

ALTER TABLE "product" ADD FOREIGN KEY ("category_id") REFERENCES "product_category" ("id");

ALTER TABLE "product" ADD FOREIGN KEY ("inventory_id") REFERENCES "product_inventory" ("id");

ALTER TABLE "product" ADD FOREIGN KEY ("discount_id") REFERENCES "discount" ("id");

ALTER TABLE "order_detail" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "order_detail" ADD FOREIGN KEY ("payment_id") REFERENCES "payment_detail" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "order_detail" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");

CREATE INDEX ON "user" ("username");

CREATE INDEX ON "user" ("telephone");

COMMENT ON COLUMN "shopping_session"."total" IS 'must be positive';

COMMENT ON COLUMN "cart_item"."quantity" IS 'must be positive';

COMMENT ON COLUMN "payment_detail"."order_id" IS 'default is 0';

COMMENT ON COLUMN "payment_detail"."amount" IS 'must be positive';

COMMENT ON COLUMN "product"."price" IS 'must be positive';

COMMENT ON COLUMN "product"."active" IS 'default is false';

COMMENT ON COLUMN "order_detail"."total" IS 'must be positive';

COMMENT ON COLUMN "order_items"."quantity" IS 'must be positive';

COMMENT ON COLUMN "product_category"."active" IS 'default is false';

COMMENT ON COLUMN "product_inventory"."quantity" IS 'must be positive';

COMMENT ON COLUMN "product_inventory"."active" IS 'default is true';

COMMENT ON COLUMN "discount"."discount_percent" IS 'default is 0';

COMMENT ON COLUMN "discount"."active" IS 'default is false';

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_user
BEFORE UPDATE ON "user"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_shopping_session
BEFORE UPDATE ON "shopping_session"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_cart_item
BEFORE UPDATE ON "cart_item"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_payment_detail
BEFORE UPDATE ON "payment_detail"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_product
BEFORE UPDATE ON "product"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_order_detail
BEFORE UPDATE ON "order_detail"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_order_items
BEFORE UPDATE ON "order_items"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_product_category
BEFORE UPDATE ON "product_category"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_product_inventory
BEFORE UPDATE ON "product_inventory"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_discount
BEFORE UPDATE ON "discount"
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
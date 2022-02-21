DROP trigger set_timestamp_admin_type on "admin_type";

DROP trigger set_timestamp_admin on "admin";

ALTER TABLE IF EXISTS "admin" DROP CONSTRAINT IF EXISTS "admin_type_id_fkey";

DROP TABLE IF EXISTS "admin";

DROP TABLE IF EXISTS "admin_type";
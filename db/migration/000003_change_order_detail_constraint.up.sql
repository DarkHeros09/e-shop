ALTER TABLE order_detail
DROP CONSTRAINT "order_detail_payment_id_fkey";

ALTER TABLE order_detail
ADD CONSTRAINT "order_detail_payment_id_fkey" 
FOREIGN KEY ("payment_id") REFERENCES "public"."payment_detail"("id") ON DELETE SET NULL;

ALTER TABLE order_detail
ALTER COLUMN payment_id
DROP NOT NULL;
-- Modify "lock_payment_orders" table
ALTER TABLE "lock_payment_orders" DROP COLUMN "order_id", DROP COLUMN "label", ADD COLUMN "gateway_id" character varying NOT NULL DEFAULT '';
-- Modify "payment_orders" table
ALTER TABLE "payment_orders" DROP COLUMN "label", ADD COLUMN "gateway_id" character varying NULL;

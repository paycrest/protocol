-- Modify "lock_payment_orders" table
ALTER TABLE "lock_payment_orders" ALTER COLUMN "order_percent" SET NOT NULL;
-- Modify "payment_orders" table
ALTER TABLE "payment_orders" ADD COLUMN "from_address" character varying NULL;

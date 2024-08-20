-- Modify "lock_payment_orders" table
ALTER TABLE "lock_payment_orders" ALTER COLUMN "gateway_id" DROP DEFAULT;
-- Create index "lock_payment_orders_gateway_id_key" to table: "lock_payment_orders"
CREATE UNIQUE INDEX "lock_payment_orders_gateway_id_key" ON "lock_payment_orders" ("gateway_id");

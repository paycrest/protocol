-- Drop index "lock_order_fulfillments_lock_payment_order_fulfillment_key" from table: "lock_order_fulfillments"
DROP INDEX "lock_order_fulfillments_lock_payment_order_fulfillment_key";
-- Rename a column from "lock_payment_order_fulfillment" to "lock_payment_order_fulfillments"
ALTER TABLE "lock_order_fulfillments" RENAME COLUMN "lock_payment_order_fulfillment" TO "lock_payment_order_fulfillments";
-- Modify "lock_order_fulfillments" table
ALTER TABLE "lock_order_fulfillments" DROP CONSTRAINT "lock_order_fulfillments_lock_payment_orders_fulfillment", ADD CONSTRAINT "lock_order_fulfillments_lock_payment_orders_fulfillments" FOREIGN KEY ("lock_payment_order_fulfillments") REFERENCES "lock_payment_orders" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;

-- Modify "payment_orders" table
ALTER TABLE "payment_orders" ALTER COLUMN "network_fee" DROP DEFAULT;
-- Modify "receive_addresses" table
ALTER TABLE "receive_addresses" ADD COLUMN "tx_hash" character varying NULL;

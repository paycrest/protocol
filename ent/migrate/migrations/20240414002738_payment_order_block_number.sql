-- Modify "networks" table
ALTER TABLE "networks" ALTER COLUMN "fee" DROP DEFAULT;
-- Modify "payment_orders" table
ALTER TABLE "payment_orders" ADD COLUMN "block_number" bigint NOT NULL DEFAULT 0;

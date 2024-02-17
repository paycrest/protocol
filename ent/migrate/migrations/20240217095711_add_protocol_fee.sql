-- Modify "payment_orders" table
ALTER TABLE "payment_orders" ALTER COLUMN "percent_settled" DROP DEFAULT, ADD COLUMN "protocol_fee" double precision NOT NULL DEFAULT 0.0;

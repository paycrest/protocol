-- Modify "payment_orders" table
ALTER TABLE "payment_orders" ADD COLUMN "percent_settled" double precision NOT NULL DEFAULT 0.0;

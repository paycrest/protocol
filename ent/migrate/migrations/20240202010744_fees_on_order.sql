-- Modify "payment_orders" table
ALTER TABLE "payment_orders" DROP COLUMN "network_fee_estimate", ADD COLUMN "network_fee" double precision NOT NULL DEFAULT 0.0, ADD COLUMN "fee_per_token_unit" double precision NOT NULL, ADD COLUMN "fee_address" character varying NULL;
-- Modify "sender_profiles" table
ALTER TABLE "sender_profiles" ADD COLUMN "is_partner" boolean NOT NULL DEFAULT false;

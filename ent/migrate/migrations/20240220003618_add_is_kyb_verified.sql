-- Modify "payment_orders" table
ALTER TABLE "payment_orders" ALTER COLUMN "protocol_fee" DROP DEFAULT;
-- Modify "provider_profiles" table
ALTER TABLE "provider_profiles" ALTER COLUMN "trading_name" DROP NOT NULL, ADD COLUMN "is_kyb_verified" boolean NOT NULL DEFAULT false;

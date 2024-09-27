-- Create "identity_verification_requests" table
CREATE TABLE "identity_verification_requests" ("id" uuid NOT NULL, "wallet_address" character varying NOT NULL, "wallet_signature" character varying NOT NULL, "platform" character varying NOT NULL, "platform_ref" character varying NOT NULL, "verification_url" character varying NOT NULL, "status" character varying NOT NULL DEFAULT 'pending', "fee_reclaimed" boolean NOT NULL DEFAULT false, "updated_at" timestamptz NOT NULL, "last_url_created_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "identity_verification_requests_wallet_address_key" to table: "identity_verification_requests"
CREATE UNIQUE INDEX "identity_verification_requests_wallet_address_key" ON "identity_verification_requests" ("wallet_address");
-- Create index "identity_verification_requests_wallet_signature_key" to table: "identity_verification_requests"
CREATE UNIQUE INDEX "identity_verification_requests_wallet_signature_key" ON "identity_verification_requests" ("wallet_signature");
-- Add pk ranges for ('identity_verification_requests') tables
INSERT INTO "ent_types" ("type") VALUES ('identity_verification_requests');

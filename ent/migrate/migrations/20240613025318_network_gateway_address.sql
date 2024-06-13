-- Modify "networks" table
ALTER TABLE "networks" ADD COLUMN "gateway_contract_address" character varying NOT NULL DEFAULT '';
-- Modify "transaction_logs" table
ALTER TABLE "transaction_logs" ALTER COLUMN "status" SET DEFAULT 'order_initiated';

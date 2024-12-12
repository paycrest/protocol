-- Rename a column from "fee_per_token_unit" to "fee_percent"
ALTER TABLE "payment_orders" RENAME COLUMN "fee_per_token_unit" TO "fee_percent";
-- Rename a column from "fee_per_token_unit" to "fee_percent"
ALTER TABLE "sender_order_tokens" RENAME COLUMN "fee_per_token_unit" TO "fee_percent";

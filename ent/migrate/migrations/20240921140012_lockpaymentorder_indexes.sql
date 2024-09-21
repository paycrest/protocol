-- Create index "lockpaymentorder_gateway_id_ra_65d1cd4f9b7a0ff4525b6f2bc506afdc" to table: "lock_payment_orders"
CREATE UNIQUE INDEX "lockpaymentorder_gateway_id_ra_65d1cd4f9b7a0ff4525b6f2bc506afdc" ON "lock_payment_orders" ("gateway_id", "rate", "tx_hash", "block_number", "institution", "account_identifier", "account_name", "memo", "token_lock_payment_orders");
-- Create index "transactionlog_gateway_id_status_network_tx_hash" to table: "transaction_logs"
CREATE UNIQUE INDEX "transactionlog_gateway_id_status_network_tx_hash" ON "transaction_logs" ("gateway_id", "status", "network", "tx_hash");

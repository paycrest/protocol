-- Drop the trigger from the payment_orders table
DROP TRIGGER IF EXISTS enforce_payment_order_amount ON payment_orders;

-- Drop the trigger function check_payment_order_amount
DROP FUNCTION IF EXISTS check_payment_order_amount;

-- Drop the calculate_total_amount function
DROP FUNCTION IF EXISTS calculate_total_amount;

-- First, ensure we have a function to calculate the total amount with fees
CREATE OR REPLACE FUNCTION calculate_total_amount(
    amount DOUBLE PRECISION,
    sender_fee DOUBLE PRECISION,
    network_fee DOUBLE PRECISION,
    protocol_fee DOUBLE PRECISION,
    token_decimals SMALLINT
) RETURNS DOUBLE PRECISION AS $$
BEGIN
    RETURN ROUND(amount + sender_fee + network_fee + protocol_fee, token_decimals);
END;
$$ LANGUAGE plpgsql;

-- Now, create a trigger function to enforce our check
CREATE OR REPLACE FUNCTION check_payment_order_amount() RETURNS TRIGGER AS $$
DECLARE
    total_amount DOUBLE PRECISION;
    token_decimals SMALLINT;
BEGIN
    -- Get the token decimals
    SELECT decimals INTO token_decimals FROM tokens WHERE id = NEW.token_payment_orders;

    -- Calculate the total amount with fees
    total_amount := calculate_total_amount(NEW.amount, NEW.sender_fee, NEW.network_fee, NEW.protocol_fee, token_decimals);

    -- Check if the amount_paid is within the valid range
    IF NEW.amount_paid < 0 OR NEW.amount_paid >= total_amount THEN
        RAISE EXCEPTION 'Duplicate payment order';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply the trigger to the payment_orders table
CREATE TRIGGER enforce_payment_order_amount
BEFORE INSERT OR UPDATE ON payment_orders
FOR EACH ROW EXECUTE FUNCTION check_payment_order_amount();
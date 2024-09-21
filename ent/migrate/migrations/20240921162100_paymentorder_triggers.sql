-- First, ensure we have a function to calculate the total amount with fees
CREATE OR REPLACE FUNCTION calculate_total_amount(
    amount DECIMAL,
    sender_fee DECIMAL,
    network_fee DECIMAL,
    protocol_fee DECIMAL,
    token_decimals INTEGER
) RETURNS DECIMAL AS $$
BEGIN
    RETURN ROUND(amount + sender_fee + network_fee + protocol_fee, token_decimals);
END;
$$ LANGUAGE plpgsql;

-- Now, create a trigger function to enforce our check
CREATE OR REPLACE FUNCTION check_payment_order_amount() RETURNS TRIGGER AS $$
DECLARE
    total_amount DECIMAL;
    token_decimals INTEGER;
BEGIN
    -- Get the token decimals
    SELECT decimals INTO token_decimals FROM tokens WHERE id = NEW.token_payment_orders;

    -- Calculate the total amount with fees
    total_amount := calculate_total_amount(NEW.amount, NEW.sender_fee, NEW.network_fee, NEW.protocol_fee, token_decimals);

    -- Check if the amount_paid is within the valid range
    IF NEW.amount_paid < 0 OR NEW.amount_paid >= total_amount THEN
        RAISE EXCEPTION 'Invalid amount_paid: must be greater than or equal to 0 and less than total order amount with fees';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply the trigger to the payment_orders table
CREATE TRIGGER enforce_payment_order_amount
BEFORE INSERT OR UPDATE ON payment_orders
FOR EACH ROW EXECUTE FUNCTION check_payment_order_amount();
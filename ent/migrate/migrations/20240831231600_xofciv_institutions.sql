-- First, check if the XOF-CIV fiat currency exists
SELECT EXISTS (
    SELECT 1 FROM "fiat_currencies"
    WHERE "code" = 'XOF-CIV'
);

-- If the XOF-CIV fiat currency exists, then add the institutions
DO $$
DECLARE
    fiat_currency_id UUID;
    last_bucket_id BIGINT;
BEGIN
    -- Get the ID of the XOF-CIV fiat currency
    SELECT "id" INTO fiat_currency_institutions
    FROM "fiat_currencies"
    WHERE "code" = 'XOF-CIV';

    -- Add institutions to the XOF-CIV fiat currency
    WITH institutions (code, name, type, updated_at, created_at) AS (
        VALUES
            ('MOMOCIPC', 'MTN Mobile Money', 'mobile_money', now(), now())
    )
    INSERT INTO "institutions" ("code", "name", "fiat_currency_institutions", "type", "updated_at", "created_at")
    SELECT "code", "name", fiat_currency_id, "type", "updated_at", "created_at"
    FROM institutions
    ON CONFLICT ("code") DO NOTHING;

    -- Get the last bucket ID
    SELECT COALESCE(MAX(id), 0) INTO last_bucket_id FROM provision_buckets;

    -- Add provision buckets to the UGX fiat currency
    INSERT INTO provision_buckets (id, min_amount, max_amount, created_at, fiat_currency_provision_buckets)
    VALUES 
        (last_bucket_id + 1, 0, 1000, now(), fiat_currency_id),
        (last_bucket_id + 2, 1001, 5000, now(), fiat_currency_id),
        (last_bucket_id + 3, 5001, 50000, now(), fiat_currency_id)
    ON CONFLICT (id) DO NOTHING;
END$$;
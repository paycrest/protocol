-- First, check if the ZAR fiat currency exists

SELECT EXISTS
    (SELECT 1
     FROM "fiat_currencies"
     WHERE "code" = 'ZAR' );

-- If the ZAR fiat currency exists, then add the institutions
DO $$
DECLARE
    fiat_currency_id UUID;
    last_bucket_id BIGINT;
BEGIN
    -- Get the ID of the ZAR fiat currency
    SELECT "id" INTO fiat_currency_id
    FROM "fiat_currencies"
    WHERE "code" = 'ZAR';

    -- Add institutions to the ZAR fiat currency
    WITH institutions (code, name, type, updated_at, created_at) AS (
        VALUES
            ('ABSAZAJJ', 'ABSA', 'bank', now(), now()),
            ('BATHZAJJ', 'Access Bank', 'bank', now(), now()),
            ('AFRCZAJJ', 'African Bank Limited', 'bank', now(), now()),
            ('BWLINANX', 'Bank Eindhoek Limited', 'bank', now(), now()),
            ('BIDBZAJJ', 'Bidvest Bank LTD', 'bank', now(), now()),
            ('CABLZAJJ', 'Capitec Bank', 'bank', now(), now()),
            ('DISCZAJJ', 'Discovery Bank Limited', 'bank', now(), now()),
            ('FIRNZAJJ', 'First National Bank', 'bank', now(), now()),
            ('MAMRZAJ1', 'Grindrod', 'bank', now(), now()),
            ('HSBCZAJJ', 'HSBC', 'bank', now(), now()),
            ('IVESZAJJ', 'Investec Bank', 'bank', now(), now()),
            ('LISAZAJJ', 'Mercantile Bank', 'bank', now(), now()),
            ('NEDSZAJJ', 'Nedbank Limited', 'bank', now(), now()),
            ('OMIGZAJ1', 'Old Mutual Investment Group', 'bank', now(), now()),
            ('SASFZAJJ', 'Sasfin Bank Limited', 'bank', now(), now()),
            ('SBZAZAJJ', 'Standard Bank of South Africa Limited', 'bank', now(), now()),
            ('MELNUS33', 'The Bank of New York Mellon', 'bank', now(), now()),
            ('CBZAZAJJ', 'Tyme Bank', 'bank', now(), now()),
            ('YOUBZAJJ', 'Ubank Limited', 'bank', now(), now())
    )
    INSERT INTO "institutions" ("code", "name", "fiat_currency_institutions", "type", "updated_at", "created_at")
    SELECT "code", "name", fiat_currency_id, "type", "updated_at", "created_at"
    FROM institutions
    ON CONFLICT ("code") DO NOTHING;

    -- Get the last bucket ID
    SELECT COALESCE(MAX(id), 0) INTO last_bucket_id FROM provision_buckets;

    -- Add provision buckets to the ZAR fiat currency
    INSERT INTO provision_buckets (id, min_amount, max_amount, created_at, fiat_currency_provision_buckets)
    VALUES
        (last_bucket_id + 1, 0, 1000, now(), fiat_currency_id),
        (last_bucket_id + 2, 1001, 5000, now(), fiat_currency_id),
        (last_bucket_id + 3, 5001, 50000, now(), fiat_currency_id)
    ON CONFLICT (id) DO NOTHING;
END$$;
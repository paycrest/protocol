-- First, check if the KES fiat currency exists
SELECT EXISTS (
    SELECT 1 FROM "fiat_currencies"
    WHERE "code" = 'KES'
);

-- If the KES fiat currency exists, then add the institutions
DO $$
DECLARE
    fiat_currency_id UUID;
    last_bucket_id BIGINT;
BEGIN
    -- Get the ID of the KES fiat currency
    SELECT "id" INTO fiat_currency_institutions
    FROM "fiat_currencies"
    WHERE "code" = 'KES';

    -- Add institutions to the KES fiat currency
    WITH institutions (code, name, type, updated_at, created_at) AS (
        VALUES
            ('KCBLKENX', 'Kenya Commercial Bank', 'bank', now(), now()),
            ('SCBLKENX', 'Standard Chartered Kenya', 'bank', now(), now()),
            ('CITIKENA', 'Citi Bank', 'bank', now(), now()),
            ('DTKEKENA', 'Diamond Trust Bank', 'bank', now(), now()),
            ('BARCKENX', 'ABSA Bank Kenya', 'bank', now(), now()),
            ('MIEKKENA', 'Middle East Bank', 'bank', now(), now()),
            ('EQBLKENA', 'Equity Bank', 'bank', now(), now()),
            ('BARBKENA', 'Bank of Baroda', 'bank', now(), now()),
            ('AFRIKENX', 'Bank of Africa', 'bank', now(), now()),
            ('SBICKENX', 'Stanbic Bank Kenya', 'bank', now(), now()),
            ('ABCLKENA', 'African Bank Corporation Limited', 'bank', now(), now()),
            ('ECOCKENA', 'Ecobank Transnational Inc.', 'bank', now(), now()),
            ('CRMFKENA', 'Caritas Microfinance Bank', 'bank', now(), now()),
            ('PAUTKENA', 'Paramount Bank', 'bank', now(), now()),
            ('CIFIKENA', 'Kingdom Bank Limited', 'bank', now(), now()),
            ('GTBIKENA', 'Guaranty Trust Holding Company PLC', 'bank', now(), now()),
            ('FABLKENA', 'Family Bank', 'bank', now(), now()),
            ('CBAFKENX', 'National Commercial Bank of Africa', 'bank', now(), now()),
            ('CONKKENA', 'Consolidated Bank Kenya', 'bank', now(), now()),
            ('HFCOKENA', 'Housing finance Company', 'bank', now(), now()),
            ('MYBKKENA', 'Commercial International Bank Kenya Limited', 'bank', now(), now()),
            ('GAFRKENA', 'Gulf African Bank', 'bank', now(), now()),
            ('FCBAKEPC', 'First Community Bank', 'bank', now(), now()),
            ('DUIBKENA', 'Dubai Islamic Bank', 'bank', now(), now()),
            ('KWMIKENX', 'Kenya Women Microfinance Bank', 'bank', now(), now()),
            ('FAUMKENA', 'Faulu Bank', 'bank', now(), now()),
            ('PRIEKENX', 'Prime Bank Limited', 'bank', now(), now()),
            ('CRBTKENA', 'Credit Bank Limited', 'bank', now(), now()),
            ('UNAIKEPC', 'Unaitas Sacco', 'bank', now(), now()),
            ('GUARKENA', 'Guardian Bank Limited', 'bank', now(), now()),
            ('IMBLKENA', 'Investments & Morgage Limited', 'bank', now(), now()),
            ('SIDNKENA', 'Sidian Bank', 'bank', now(), now()),
            ('UNAFKENA', 'United Bank for Africa', 'bank', now(), now()),
            ('KCOOKENA', 'Cooperative Bank of Kenya', 'bank', now(), now()),
            ('CHFIKENX', 'Choice Microfinance Bank Kenya Limited', 'bank', now(), now()),
            ('STIMKEPC', 'Stima SACCO', 'bank', now(), now()),
            ('NBKEKENX', 'National Bank of kenya', 'bank', now(), now()),
            ('VICMKENA', 'Victoria Commercial Bank', 'bank', now(), now()),
            ('MORBKENA', 'Oriental Commercial Bank Limited', 'bank', now(), now()),
            ('SBMKKENA', 'SBM Bank Kenya', 'bank', now(), now()),
            ('SAFAKEPC', 'SAFARICOM', 'mobile_money', now(), now()),
            ('AIRTKEPC', 'AIRTEL', 'mobile_money', now(), now())
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
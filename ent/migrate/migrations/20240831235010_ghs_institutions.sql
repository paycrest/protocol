-- First, check if the GHS fiat currency exists
SELECT EXISTS (
    SELECT 1 FROM "fiat_currencies"
    WHERE "code" = 'GHS'
);

-- If the GHS fiat currency exists, then add the institutions
DO $$
DECLARE
    fiat_currency_id UUID;
    last_bucket_id BIGINT;
BEGIN
    -- Get the ID of the GHS fiat currency
    SELECT "id" INTO fiat_currency_id
    FROM "fiat_currencies"
    WHERE "code" = 'GHS';

    -- Add institutions to the GHS fiat currency
    WITH institutions (code, name, type, updated_at, created_at) AS (
        VALUES
            ('HFCAGHAC', 'REPUBLIC BANK', 'bank', now(), now()),
            ('BAGHGHAC', 'BANK OF GHANA', 'bank', now(), now()),
            ('STBGGHAC', 'UBA GHANA', 'bank', now(), now()),
            ('FBLIGHAC', 'FIDELITY BANK GHANA', 'bank', now(), now()),
            ('CBGHGHAC', 'Consolidated Bank Ghana', 'bank', now(), now()),
            ('MBGHGHAC', 'Universal Merchant Bank', 'bank', now(), now()),
            ('AREXGHAC', 'APEX BANK', 'bank', now(), now()),
            ('OMBLGHAC', 'OMNI BANK', 'bank', now(), now()),
            ('NIBGGHAC', 'National Investment Bank Limited', 'bank', now(), now()),
            ('ADNTGHAC', 'Agricultural Development Bank', 'bank', now(), now()),
            ('ZEBLGHAC', 'ZENITH BANK GHANA', 'bank', now(), now()),
            ('ZEEPGHPC', 'ZEEPAY GHANA', 'bank', now(), now()),
            ('SBICGHAC', 'STANBIC BANK GHANA', 'bank', now(), now()),
            ('SCBLGHAC', 'STANDARD CHARTERED GHANA', 'bank', now(), now()),
            ('SISLGHPC', 'SERVICES INTEGRITY SAVINGS & LOANS', 'bank', now(), now()),
            ('BARCGHAC', 'ABSA BANK GHANA', 'bank', now(), now()),
            ('FAMCGHAC', 'First Atlantic Bank GHANA', 'bank', now(), now()),
            ('ABNGGHAC', 'ACCESS BANK GHANA', 'bank', now(), now()),
            ('FIRNGHAC', 'First National Bank', 'bank', now(), now()),
            ('GMONGHPC', 'G-MONEY', 'bank', now(), now()),
            ('GHCBGHAC', 'GCB Bank Limited', 'bank', now(), now()),
            ('GHLLGHPC', 'Ghl Bank Limited', 'bank', now(), now()),
            ('AMMAGHAC', 'Bank of Africa Ghana Limited', 'bank', now(), now()),
            ('PUBKGHAC', 'Prudential Bank Limited', 'bank', now(), now()),
            ('INCEGHAC', 'FBNBank Ghana Limited', 'bank', now(), now()),
            ('ECOCGHAC', 'ECOBANK GHANA', 'bank', now(), now()),
            ('ACCCGHAC', 'CAL BANK', 'bank', now(), now()),
            ('SSEBGHAC', 'Societe Generale Ghana Limited', 'bank', now(), now()),
            ('GTBIGHAC', 'GT BANK GHANA', 'bank', now(), now()),
            ('MOMOGHPC', 'MTN Mobile Money', 'mobile_money', now(), now()),
            ('VODAGHPC', 'Vodafone Cash', 'mobile_money', now(), now()),
            ('AIRTGHPC', 'AirtelTigo Money', 'mobile_money', now(), now())
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
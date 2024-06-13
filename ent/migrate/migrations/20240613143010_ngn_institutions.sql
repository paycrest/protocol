-- First, check if the NGN fiat currency exists
SELECT EXISTS (
    SELECT 1 FROM "fiat_currencies"
    WHERE "code" = 'NGN'
);

-- If the NGN fiat currency exists, then add the institutions
DO $$
DECLARE
    fiat_currency_institutions UUID;
BEGIN
    -- Get the ID of the NGN fiat currency
    SELECT "id" INTO fiat_currency_institutions
    FROM "fiat_currencies"
    WHERE "code" = 'NGN';

    -- Add institutions to the NGN fiat currency
    WITH institutions (code, name, type, updated_at, created_at) AS (
        VALUES
            ('ABNGNGLA', 'Access Bank', 'bank', now(), now()),
            ('DBLNNGLA', 'Diamond Bank', 'bank', now(), now()),
            ('FIDTNGLA', 'Fidelity Bank', 'bank', now(), now()),
            ('FCMBNGLA', 'FCMB', 'bank', now(), now()),
            ('FBNINGLA', 'First Bank Of Nigeria', 'bank', now(), now()),
            ('GTBINGLA', 'Guaranty Trust Bank', 'bank', now(), now()),
            ('PRDTNGLA', 'Polaris Bank', 'bank', now(), now()),
            ('UBNINGLA', 'Union Bank', 'bank', now(), now()),
            ('UNAFNGLA', 'United Bank for Africa', 'bank', now(), now()),
            ('CITINGLA', 'Citibank', 'bank', now(), now()),
            ('ECOCNGLA', 'Ecobank Bank', 'bank', now(), now()),
            ('HBCLNGLA', 'Heritage', 'bank', now(), now()),
            ('PLNINGLA', 'Keystone Bank', 'bank', now(), now()),
            ('SBICNGLA', 'Stanbic IBTC Bank', 'bank', now(), now()),
            ('SCBLNGLA', 'Standard Chartered Bank', 'bank', now(), now()),
            ('NAMENGLA', 'Sterling Bank', 'bank', now(), now()),
            ('ICITNGLA', 'Unity Bank', 'bank', now(), now()),
            ('SUTGNGLA', 'Suntrust Bank', 'bank', now(), now()),
            ('PROVNGLA', 'Providus Bank ', 'bank', now(), now()),
            ('KDHLNGLA', 'FBNQuest Merchant Bank', 'bank', now(), now()),
            ('GMBLNGLA', 'Greenwich Merchant Bank', 'bank', now(), now()),
            ('FSDHNGLA', 'FSDH Merchant Bank', 'bank', now(), now()),
            ('FIRNNGLA', 'Rand Merchant Bank', 'bank', now(), now()),
            ('JAIZNGLA', 'Jaiz Bank', 'bank', now(), now()),
            ('ZEIBNGLA', 'Zenith Bank', 'bank', now(), now()),
            ('WEMANGLA', 'Wema Bank', 'bank', now(), now()),
            ('KUDANGPC', 'Kuda Microfinance Bank', 'bank', now(), now()),
            ('OPAYNGPC', 'OPay', 'bank', now(), now()),
            ('MONINGPC', 'Moniepoint Microfinance Bank', 'bank', now(), now()),
            ('PALMNGPC', 'PalmPay', 'bank', now(), now()),
            ('SAHVNGPC', 'Safehaven Microfinance Bank', 'bank', now(), now()),
            ('PAYTNGPC', 'Paystack-Titan MFB', 'bank', now(), now())
    )
    INSERT INTO "institutions" ("code", "name", "fiat_currency_institutions", "type", "updated_at", "created_at")
    SELECT "code", "name", fiat_currency_institutions, "type", "updated_at", "created_at"
    FROM institutions
    ON CONFLICT ("code") DO NOTHING;
END$$;
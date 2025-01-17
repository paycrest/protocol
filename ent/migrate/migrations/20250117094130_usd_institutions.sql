-- First, check if the USD fiat currency exists

SELECT EXISTS
    ( SELECT 1
     FROM "fiat_currencies"
     WHERE "code" = 'USD' );

-- If the USD fiat currency exists, then add the institutions
DO $$
DECLARE
    fiat_currency_id UUID;
    last_bucket_id BIGINT;
BEGIN
    -- Get the ID of the USD fiat currency
    SELECT "id" INTO fiat_currency_id
    FROM "fiat_currencies"
    WHERE "code" = 'USD';

    -- Add institutions to the USD fiat currency
    WITH institutions (code, name, type, updated_at, created_at) AS (
        VALUES
            ('ALLYUS31', 'Ally Bank', 'bank', now(), now()),
            ('BOFAUS3N', 'Bank of America', 'bank', now(), now()),
            ('SFRUUS33', 'Bank-Funds Federal Credit Union', 'bank', now(), now()),
            ('BBSRUS31', 'BB&T', 'bank', now(), now()),
            ('TCBKUSD1', 'Chase Bank', 'bank', now(), now()),
            ('CHFGUS44', 'Choice Financial Group', 'bank', now(), now()),
            ('CITIUS33', 'Citibank', 'bank', now(), now()),
            ('CTZIUS33', 'Citizens Bank', 'bank', now(), now()),
            ('CLNOUS66', 'Column NA-Brex', 'bank', now(), now()),
            ('CMFGUS33', 'Community Federal Savings Bank', 'bank', now(), now()),
            ('EVOVUS42', 'Evolve Bank And Trust', 'bank', now(), now()),
            ('FHLBUS44', 'Federal Home Loan Bank', 'bank', now(), now()),
            ('FTBCUS3C', 'Fifth Third Bank', 'bank', now(), now()),
            ('FBERUS31', 'First Bank of Berne', 'bank', now(), now()),
            ('FCBIUS42', 'First Citizens Bank', 'bank', now(), now()),
            ('FRNAUS44', 'First National Bankers Bank', 'bank', now(), now()),
            ('MRMDUS33', 'HSBC', 'bank', now(), now()),
            ('JSBAUS41', 'Jefferson Bank', 'bank', now(), now()),
            ('CHASUS33', 'JPMorgan Chase Bank', 'bank', now(), now()),
            ('LIBYUS32', 'Liberty Bank', 'bank', now(), now()),
            ('MANTUS3A', 'M&T Bank', 'bank', now(), now()),
            ('MARQUS41', 'Marquette Bank', 'bank', now(), now()),
            ('MCBEUS33', 'Metropolitan Commercial Bank', 'bank', now(), now()),
            ('PNCCUS33', 'PNC Bank', 'bank', now(), now()),
            ('HOSVUS33', 'Premier Bank', 'bank', now(), now()),
            ('PROYUS44', 'Prosperity Bank', 'bank', now(), now()),
            ('PRNDUS33', 'The Provident Bank', 'bank', now(), now()),
            ('UPNBUS44', 'Regions Bank', 'bank', now(), now()),
            ('RTCOUS33', 'Rockland Trust Company', 'bank', now(), now()),
            ('SVBKUS6S', 'Silicon Valley Bank', 'bank', now(), now()),
            ('SFIBUS44', 'Star Financial Bank', 'bank', now(), now()),
            ('SNTRUS3A', 'Suntrust Bank Atlanta Los Alamit', 'bank', now(), now()),
            ('SUWWUS6I', 'Sunwest Bank', 'bank', now(), now()),
            ('SFGLUS31', 'Susquehanna Financial Group LLLP', 'bank', now(), now()),
            ('SYNEUS44', 'Synergy Bank', 'bank', now(), now()),
            ('FICOUS44', 'Synovus Bank', 'bank', now(), now()),
            ('NRTHUS33', 'TD Bank', 'bank', now(), now()),
            ('TENTUS44', 'Texas National Bank', 'bank', now(), now()),
            ('IRVTUS3N', 'The Bank of New York Mellon', 'bank', now(), now()),
            ('FMSBUS3A', 'The Farmers And Merchant State Bank', 'bank', now(), now()),
            ('TIMNUS55', 'Timberline Bank', 'bank', now(), now()),
            ('TRWIUS31', 'Transferwise Inc', 'bank', now(), now()),
            ('BRBTUS33', 'Truist Bank', 'bank', now(), now()),
            ('UMKCUS44', 'UMB Bank', 'bank', now(), now()),
            ('USBKUS44', 'US Bank National Association', 'bank', now(), now()),
            ('UFSBUS44', 'USAA Federal Saving Bank', 'bank', now(), now()),
            ('MBNYUS33', 'Valley National Bank', 'bank', now(), now()),
            ('WFBIUS6S', 'Wells Fargo Bank, N.A.', 'bank', now(), now()),
            ('PNBPUS3W', 'Wells Fargo Clearing Services, LLC', 'bank', now(), now())

    )
    INSERT INTO "institutions" ("code", "name", "fiat_currency_institutions", "type", "updated_at", "created_at")
    SELECT "code", "name", fiat_currency_id, "type", "updated_at", "created_at"
    FROM institutions
    ON CONFLICT ("code") DO NOTHING;

    -- Get the last bucket ID
    SELECT COALESCE(MAX(id), 0) INTO last_bucket_id FROM provision_buckets;

    -- Add provision buckets to the USD fiat currency
    INSERT INTO provision_buckets (id, min_amount, max_amount, created_at, fiat_currency_provision_buckets)
    VALUES
        (last_bucket_id + 1, 0, 1000, now(), fiat_currency_id),
        (last_bucket_id + 2, 1001, 5000, now(), fiat_currency_id),
        (last_bucket_id + 3, 5001, 50000, now(), fiat_currency_id)
    ON CONFLICT (id) DO NOTHING;
END$$;
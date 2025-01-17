-- First, check if the BRL fiat currency exists

SELECT EXISTS
    ( SELECT 1
     FROM "fiat_currencies"
     WHERE "code" = 'BRL' );

-- If the BRL fiat currency exists, then add the institutions
DO $$
DECLARE
    fiat_currency_id UUID;
    last_bucket_id BIGINT;
BEGIN
    -- Get the ID of the BRL fiat currency
    SELECT "id" INTO fiat_currency_id
    FROM "fiat_currencies"
    WHERE "code" = 'BRL';

    -- Add institutions to the BRL fiat currency
    WITH institutions (code, name, type, updated_at, created_at) AS (
        VALUES
            ('ABCBBRSP', 'Banco ABC Brasil S.A.', 'bank', now(), now()),
            ('AUFABRSP', 'Banco Alfa S.A.', 'bank', now(), now()),
            ('BMBCBRSP', 'Banco Bmg S.A.', 'bank', now(), now()),
            ('BBONBRSP', 'Banco Bonsucesso S.A.', 'bank', now(), now()),
            ('BBDEBRSP', 'Banco Bradesco S.A.', 'bank', now(), now()),
            ('CSIXBRSP', 'Banco C6 S.A.', 'bank', now(), now()),
            ('BCBZBRSP', 'Banco Cargill S.A.', 'bank', now(), now()),
            ('CITIBRBR', 'Banco Citibank S.A.', 'bank', now(), now()),
            ('CPBNBRDF', 'Banco Cooperativo do Brasil S.A. - BANCOOB', 'bank', now(), now()),
            ('BCSIBRRS', 'Banco Cooperativo Sicredi S.A.', 'bank', now(), now()),
            ('AMABBRAB', 'Banco da Amazônia S.A.', 'bank', now(), now()),
            ('DAYCBRSP', 'Banco Daycoval S.A.', 'bank', now(), now()),
            ('BRASBRRJ', 'Banco do Brasil S.A.', 'bank', now(), now()),
            ('SEEBBRR1', 'Banco do Estado de Sergipe S.A.', 'bank', now(), now()),
            ('BEPABRAB', 'Banco do Estado do Para S.A. - BANPARA', 'bank', now(), now()),
            ('BRGSBRRS', 'Banco do Estado do Rio Grande do Sul S.A.', 'bank', now(), now()),
            ('BNBRBRCF', 'Banco do Nordeste', 'bank', now(), now()),
            ('BNIIBRSP', 'Banco Inbursa S.A.', 'bank', now(), now()),
            ('ITEMBRSP', 'Banco Intermedium', 'bank', now(), now()),
            ('BMBRBRRB', 'Banco Mercantil do Brasil S.A.', 'bank', now(), now()),
            ('MODABRRJ', 'Banco Modal S.A.', 'bank', now(), now()),
            ('NACDBRRJ', 'Banco Nacional de Desenvolvimento Economico e Social', 'bank', now(), now()),
            ('BORIBRSP', 'Banco Original', 'bank', now(), now()),
            ('BJBSBRSP', 'Banco Original do Agronegócio S.A.', 'bank', now(), now()),
            ('BPNMBRSP', 'Banco Panamericano S.A.', 'bank', now(), now()),
            ('SAFRBRSP', 'Banco Safra S.A.', 'bank', now(), now()),
            ('BSCHBRSP', 'Banco Santander Brasil S.A.', 'bank', now(), now()),
            ('BSBSBRSP', 'Banco Sofisa S.A.', 'bank', now(), now()),
            ('BAVOBRSP', 'Banco Votorantim S.A.', 'bank', now(), now()),
            ('BEESBRRJ', 'Banestes', 'bank', now(), now()),
            ('BNYMBRRJ', 'Bny Mellon Banco S.A.', 'bank', now(), now()),
            ('BRBSBRDF', 'BRB - Banco de Brasília S.A.', 'bank', now(), now()),
            ('CEFXBRSP', 'Caixa Economica Federal - CEF', 'bank', now(), now()),
            ('CITIBRSP', 'Citibank N.A.', 'bank', now(), now()),
            ('HEDGBRS1', 'Credit Suisse Hedging-Griffo Corretora de Valores S.A.', 'bank', now(), now()),
            ('GOLDBRSP', 'Goldman Sachs do Brasil Banco Múltiplo S.A.', 'bank', now(), now()),
            ('ITAUBRSP', 'Itaú Unibanco S.A.', 'bank', now(), now()),
            ('NUPGBRPC', 'Nu Pagamentos (Nubank)', 'bank', now(), now()),
            ('PSITBRPC', 'PagSeguro Internet S.A.', 'bank', now(), now()),
            ('BCIIBRSP', 'State Street Brazil S.A.', 'bank', now(), now()),
            ('LICTBRS1', 'UBS Brazil Corretora de Câmbio, Títulos e Valores Mobiliários S.A.', 'bank', now(), now())

    )
    INSERT INTO "institutions" ("code", "name", "fiat_currency_institutions", "type", "updated_at", "created_at")
    SELECT "code", "name", fiat_currency_id, "type", "updated_at", "created_at"
    FROM institutions
    ON CONFLICT ("code") DO NOTHING;

    -- Get the last bucket ID
    SELECT COALESCE(MAX(id), 0) INTO last_bucket_id FROM provision_buckets;

    -- Add provision buckets to the BRL fiat currency
    INSERT INTO provision_buckets (id, min_amount, max_amount, created_at, fiat_currency_provision_buckets)
    VALUES
        (last_bucket_id + 1, 0, 1000, now(), fiat_currency_id),
        (last_bucket_id + 2, 1001, 5000, now(), fiat_currency_id),
        (last_bucket_id + 3, 5001, 50000, now(), fiat_currency_id)
    ON CONFLICT (id) DO NOTHING;
END$$;

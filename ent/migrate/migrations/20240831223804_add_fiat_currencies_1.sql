INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'KES', 'KES', 2, 'KSh', 'Kenyan Shilling', 129.00, true, now(), now()
);

INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'GHS', 'Cedi', 2, 'GH¢', 'Ghana Cedi', 15.65, false, now(), now()
);

INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'TZS', 'TZS', 2, 'TSh', 'Tanzanian Shilling', 3716.44, false, now(), now()
);

INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'UGX', 'UGX', 2, 'USh', 'Ugandan Shilling', 3716.44, false, now(), now()
);

INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'XOF-BEN', 'Céfa Benin', 2, 'CFA', 'West African CFA franc', 599.50, false, now(), now()
);

INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'XOF-CIV', 'Côte d''Ivoire', 2, 'CFA', 'West African CFA franc', 599.50, false, now(), now()
);

INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'BRL', 'BRL', 2, 'R$', 'Brazilian Real', 6.222336, false, now(), now()
);

INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'USD', 'USD', 2, '$', 'United States Dollar', 1.027908, false, now(), now()
);

INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'ZAR', 'ZAR', 2, 'R', 'South African Rand', 19.330148, false, now(), now()
);

INSERT INTO fiat_currencies (
    id, code, short_name, decimals, symbol, name, market_rate, is_enabled, created_at, updated_at
) VALUES (
    gen_random_uuid(), 'INR', 'INR', 2, '₹', 'Indian Rupee', 89.031731, false, now(), now()
);
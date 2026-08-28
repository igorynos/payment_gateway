CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice TEXT NOT NULL,
    status VARCHAR(20) NOT NULL
        CHECK (status IN ('NEW', 'PROCESSING', 'COMPLETED','FAILED','CANCELED')),
    amount NUMERIC(18, 2),
    currency VARCHAR(3) NOT NULL
        CHECK (currency IN ('RUB', 'USD', 'EUR')),
    provider VARCHAR(20) NOT NULL,
    provider_payment_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
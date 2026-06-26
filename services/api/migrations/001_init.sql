CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    display_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'EUR',
    balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    parent_id UUID REFERENCES categories(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS categorization_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pattern TEXT NOT NULL,
    match_type TEXT NOT NULL DEFAULT 'contains',
    category_id UUID NOT NULL REFERENCES categories(id),
    priority INT NOT NULL DEFAULT 100,
    confidence_score NUMERIC(5,2) NOT NULL DEFAULT 0.90,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    normalized_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS import_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'validated',
    detected_rows INT NOT NULL DEFAULT 0,
    valid_rows INT NOT NULL DEFAULT 0,
    invalid_rows INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    merchant_id UUID REFERENCES merchants(id),
    category_id UUID REFERENCES categories(id),
    import_batch_id UUID REFERENCES import_batches(id),
    booked_at DATE NOT NULL,
    label TEXT NOT NULL,
    raw_label TEXT NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'EUR',
    confidence_score NUMERIC(5,2),
    source_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS import_rows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    import_batch_id UUID NOT NULL REFERENCES import_batches(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    valid BOOLEAN NOT NULL,
    raw_data JSONB NOT NULL DEFAULT '{}',
    errors JSONB NOT NULL DEFAULT '[]',
    transaction_id UUID REFERENCES transactions(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    impact_amount NUMERIC(18,2),
    confidence_score NUMERIC(5,2) NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action TEXT NOT NULL,
    entity_type TEXT,
    entity_id UUID,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO users (email, display_name)
VALUES ('demo@noversia.com', 'Demo Noversia')
ON CONFLICT (email) DO NOTHING;

INSERT INTO accounts (user_id, name, type, currency, balance)
SELECT id, 'Compte courant demo', 'checking', 'EUR', 0
FROM users
WHERE email = 'demo@noversia.com'
AND NOT EXISTS (SELECT 1 FROM accounts WHERE name = 'Compte courant demo');

INSERT INTO categories (name) VALUES
('Courses'), ('Revenus'), ('Abonnements'), ('Transport'), ('Autres')
ON CONFLICT (name) DO NOTHING;

INSERT INTO categorization_rules (pattern, category_id, priority, confidence_score)
SELECT 'CARREFOUR', id, 10, 0.95 FROM categories WHERE name = 'Courses'
ON CONFLICT DO NOTHING;

INSERT INTO categorization_rules (pattern, category_id, priority, confidence_score)
SELECT 'NETFLIX', id, 10, 0.95 FROM categories WHERE name = 'Abonnements'
ON CONFLICT DO NOTHING;

INSERT INTO categorization_rules (pattern, category_id, priority, confidence_score)
SELECT 'SALAIRE', id, 10, 0.99 FROM categories WHERE name = 'Revenus'
ON CONFLICT DO NOTHING;

INSERT INTO categorization_rules (pattern, category_id, priority, confidence_score)
SELECT 'TOTAL', id, 20, 0.90 FROM categories WHERE name = 'Transport'
ON CONFLICT DO NOTHING;

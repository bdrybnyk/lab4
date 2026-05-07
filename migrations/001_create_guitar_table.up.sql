CREATE TABLE IF NOT EXISTS guitar (
    id UUID PRIMARY KEY,
    brand TEXT NOT NULL,
    model TEXT NOT NULL,
    strings INT NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL
);
CREATE TABLE IF NOT EXISTS cameras(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    mac_address text NOT NULL,
    model_no text NOT NULL,
    site_name text NOT NULL,
    username text NOT NULL DEFAULT 'root',
    password text NOT NULL DEFAULT 'pass'
);

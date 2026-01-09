CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS products (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  type text NOT NULL,
  flavor text NOT NULL,
  size varchar(10) NOT NULL,
  price integer NOT NULL,
  stock_qty integer NOT NULL,
  manufactured_date date NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (length(btrim(name)) > 0),
  CHECK (length(btrim(type)) > 0),
  CHECK (length(btrim(flavor)) > 0),
  CHECK (size IN ('Small', 'Medium', 'Large')),
  CHECK (price >= 0),
  CHECK (stock_qty >= 0)
);

CREATE INDEX IF NOT EXISTS products_manufactured_date_idx
  ON products (manufactured_date);

CREATE INDEX IF NOT EXISTS products_type_idx ON products (type);
CREATE INDEX IF NOT EXISTS products_flavor_idx ON products (flavor);
CREATE INDEX IF NOT EXISTS products_size_idx ON products (size);

DROP TRIGGER IF EXISTS products_set_updated_at ON products;
CREATE TRIGGER products_set_updated_at
BEFORE UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS customers (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  points integer NOT NULL DEFAULT 0,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (length(btrim(name)) > 0),
  CHECK (points >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS customers_lower_name_key
  ON customers (lower(btrim(name)));

DROP TRIGGER IF EXISTS customers_set_updated_at ON customers;
CREATE TRIGGER customers_set_updated_at
BEFORE UPDATE ON customers
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS transactions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id uuid NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
  product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  qty integer NOT NULL,
  unit_price integer NOT NULL,
  total_price integer NOT NULL,
  points_earned integer NOT NULL,
  transaction_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (qty > 0),
  CHECK (unit_price >= 0),
  CHECK (total_price >= 0),
  CHECK (points_earned >= 0)
);

CREATE INDEX IF NOT EXISTS transactions_transaction_at_idx
  ON transactions (transaction_at);

CREATE INDEX IF NOT EXISTS transactions_customer_id_idx ON transactions (customer_id);
CREATE INDEX IF NOT EXISTS transactions_product_id_idx ON transactions (product_id);
CREATE INDEX IF NOT EXISTS transactions_product_time_idx ON transactions (product_id, transaction_at);

CREATE TABLE IF NOT EXISTS redemptions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id uuid NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
  product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  qty integer NOT NULL,
  points_spent integer NOT NULL,
  redeem_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (qty > 0),
  CHECK (points_spent >= 0)
);

CREATE INDEX IF NOT EXISTS redemptions_redeem_at_idx ON redemptions (redeem_at);
CREATE INDEX IF NOT EXISTS redemptions_customer_id_idx ON redemptions (customer_id);
CREATE INDEX IF NOT EXISTS redemptions_product_id_idx ON redemptions (product_id);

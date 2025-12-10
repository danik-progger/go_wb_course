-- Create enum type for category
CREATE TYPE transaction_type AS ENUM ('income', 'expense');

-- Create the sales table
CREATE TABLE sales (
  id       BIGSERIAL PRIMARY KEY,
  type transaction_type NOT NULL,
  date     DATE NOT NULL,
  amount   NUMERIC NOT NULL CHECK (amount >= 0),
  category TEXT
);

-- Optional: Add an index on date for common queries
CREATE INDEX idx_sales_date ON sales (date);

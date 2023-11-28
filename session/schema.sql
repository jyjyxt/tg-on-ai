CREATE TABLE IF NOT EXISTS perpetuals (
  symbol                 TEXT NOT NULL,
  base_asset             TEXT NOT NULL,
  quote_asset            TEXT NOT NULL,
  categories             TEXT NOT NULL,
  source                 TEXT NOT NULL,
  mark_price             REAL NOT NULL,
  last_funding_rate      REAL NOT NULL,
  PRIMARY KEY(symbol)
);

CREATE INDEX IF NOT EXISTS perpetuals_funding_rate ON perpetuals(last_funding_rate);

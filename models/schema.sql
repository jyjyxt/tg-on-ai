CREATE TABLE IF NOT EXISTS perpetuals (
  symbol                 TEXT NOT NULL,
  base_asset             TEXT NOT NULL,
  quote_asset            TEXT NOT NULL,
  categories             TEXT NOT NULL,
  PRIMARY KEY(symbol)
);

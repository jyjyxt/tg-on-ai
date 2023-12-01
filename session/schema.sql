CREATE TABLE IF NOT EXISTS perpetuals (
  symbol                 TEXT NOT NULL,
  base_asset             TEXT NOT NULL,
  quote_asset            TEXT NOT NULL,
  categories             TEXT NOT NULL,
  source                 TEXT NOT NULL,
  mark_price             REAL NOT NULL,
  last_funding_rate      REAL NOT NULL,
  open_interest_value    REAL NOT NULL,
  updated_at             INTEGER NOT NULL,

  PRIMARY KEY(symbol)
);

CREATE INDEX IF NOT EXISTS perpetuals_funding_rate ON perpetuals(last_funding_rate);


CREATE TABLE IF NOT EXISTS candles (
  symbol                 TEXT NOT NULL,
  open                   REAL    NOT NULL, -- 开盘价
  high                   REAL    NOT NULL, -- 最高价
  low                    REAL    NOT NULL, -- 最低价
  close                  REAL    NOT NULL, -- 收盘价
  volume                 REAL    NOT NULL, -- 成交量
  open_time              INTEGER NOT NULL, -- 开盘时间
  close_time             INTEGER NOT NULL, -- 收盘时间

  PRIMARY KEY(symbol, open_time)
);


CREATE TABLE IF NOT EXISTS strategies (
  symbol                 TEXT NOT NULL,
  name                   TEXT NOT NULL,
  action                 INTEGER NOT NULL,
  score_x                REAL NOT NULL,
  score_y                REAL NOT NULL,
  open_time              INTEGER NOT NULL,
  PRIMARY KEY(symbol, name)
);

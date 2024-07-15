export interface Trend {
  symbol: string;
  high: number;
  low: number;
  now: number;
  up: number;
  down?: number;
}

export interface Perpetual {
  symbol: string;
  base_asset: string;
  quote_asset: string;
  categories: string;
  source: string;

  mark_price: number;
  last_funding_rate: number;
  sum_open_interest_value: number;

  updated_at: number;
  coingecko: string;

  trend?: Trend;
}

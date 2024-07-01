import { Card, List } from "flowbite-react";
import BigNumber from 'bignumber.js'
import { Perpetual, Trend } from '@/http/types';
import { formatNumber } from '@/utils/number';

interface prop {
  p: Perpetual
}

const fields: string[][] = [
  ['Price', 'now'],
  ['High', 'high'],
  ['Low', 'low'],
]

const Index = ({ p }: prop) => {
  const t = p.trend as Trend

  return (
    <Card className="w-96 flex-grow">
      <div className="flex items-center justify-between">
        <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          <span className="text-green-600 dark:text-green-300">{t.up}%</span> / <span className="text-red-600 dark:text-red-300">{t.down * -1}%</span>
        </h5>
        <span> {p.symbol} </span>
      </div>
      <List unstyled className="w-full">
        { fields && fields.map((f) => {
          return (
            <List.Item key={f[0]} className="sm:py-1 flex">
              <div className="min-w-0 flex-1"> 
                <p className="truncate text-sm font-medium text-gray-900 dark:text-white">{f[0]}</p>
              </div>
              <div className="inline-flex items-center text-base font-semibold text-gray-900 dark:text-white">${t[f[1] as keyof Trend]}</div>
            </List.Item>
          )
        })}
        <List.Item className="sm:py-1 flex">
          <div className="min-w-0 flex-1"> 
            <p className="truncate text-sm font-medium text-gray-900 dark:text-white">Funding Rate</p>
          </div>
          <div className="inline-flex items-center text-base font-semibold text-gray-900 dark:text-white">{BigNumber(p.last_funding_rate).times(100).toFormat()}%</div>
        </List.Item>
        <List.Item className="sm:py-1 flex">
          <div className="min-w-0 flex-1"> 
            <p className="truncate text-sm font-medium text-gray-900 dark:text-white">Interest Value</p>
          </div>
          <div className="inline-flex items-center text-base font-semibold text-gray-900 dark:text-white">${formatNumber(Math.floor(p.sum_open_interest_value))}</div>
        </List.Item>
      </List>
    </Card>
  )
}

export default Index

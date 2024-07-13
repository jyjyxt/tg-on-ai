import { Card, Badge, List } from "flowbite-react";
import BigNumber from 'bignumber.js'
import { Perpetual, Trend } from '@/apis/types';
import { formatNumber } from '@/utils/number';
import { formatDateFromNow } from '@/utils/date'

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
    <Card className="w-80 flex-grow">
      <div className="flex items-center justify-between">
        <h5 className="text-xl font-bold tracking-tight text-gray-900 dark:text-white">
          { t.down && <><span className="text-green-600 dark:text-green-300">{t.up}%</span> / <span className="text-red-600 dark:text-red-300">{t.down * -1}%</span></> }
          { !t.down && <span className={t.up > 0 ? 'text-green-600 dark:text-green-300' : 'text-red-600 dark:text-red-300'}>{t.up} Days</span> }
        </h5>
        <span className="text-sm"> {p.symbol} </span>
      </div>
      <div className="flex-1">
        <div className="flex flex-wrap gap-2">
          { fields && fields.map((f) => {
            return (
              <Badge key={f[0]} color="info">{f[0]}: ${t[f[1] as keyof Trend]}</Badge>
            )
          })}
          <Badge color="indigo">Funding Rate: {BigNumber(p.last_funding_rate).times(100).toFormat()}%</Badge>
          <Badge color="purple">Interest Value: {formatNumber(Math.floor(p.sum_open_interest_value))}</Badge>
        </div>
      </div>
      <div className="text-gray-500 dark:text-gray-400 text-sm">
        {formatDateFromNow(new Date(p.updated_at))}
      </div>
    </Card>
  )
}

export default Index

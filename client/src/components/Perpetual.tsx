import { Card } from "flowbite-react";
import { Perpetual, Trend } from '@/http/types'

interface prop {
  p: Perpetual
}

const Index = ({ p }: prop) => {
  const t = p.trend as Trend

  return (
    <Card className="max-w-sm">
      <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
        {t.high} / {t.low}
      </h5>
      <div> {p.symbol} </div>
    </Card>
  )
}

export default Index

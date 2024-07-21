import { Perpetual } from '@/apis/types'
import { initClient } from '@/apis/request'
import Perp from '@/components/Perpetual'
import DefaultLayout from '@/components/Layout'

const lowUp = (a: Perpetual, b: Perpetual) => b.trend!.up - a.trend!.up
const lowDown = (a: Perpetual, b: Perpetual) => a.trend!.up - b.trend!.up
const highUp = (a: Perpetual, b: Perpetual) => b.trend!.down! - a.trend!.down!
const highDown = (a: Perpetual, b: Perpetual) => a.trend!.down! - b.trend!.down!
const rateUp = (a: Perpetual, b: Perpetual) => b.last_funding_rate - a.last_funding_rate
const rateDown = (a: Perpetual, b: Perpetual) => a.last_funding_rate - b.last_funding_rate
const valueUp = (a: Perpetual, b: Perpetual) => b.sum_open_interest_value - a.sum_open_interest_value
const valueDown = (a: Perpetual, b: Perpetual) => a.sum_open_interest_value - b.sum_open_interest_value

const Index = async ({ params }: { params: { slug: string } }) => {
  const client = initClient()

  const arr = params.slug.split('-')

  const s: Perpetual[] = await client.perpetuals(arr[0]);
  let sort = lowUp
  if (params.slug.includes('low-down')) {
    sort = lowDown
  } else if (params.slug.includes('asc-down')) {
    sort = lowDown
  } else if (params.slug.includes('desc-down')) {
    sort = lowDown
  } else if (params.slug.includes('high-up')) {
    sort = highUp
  } else if (params.slug.includes('high-down')) {
    sort = highDown
  } else if (params.slug.includes('rate-up')) {
    sort = rateUp
  } else if (params.slug.includes('rate-down')) {
    sort = rateDown
  } else if (params.slug.includes('value-up')) {
    sort = valueUp
  } else if (params.slug.includes('value-down')) {
    sort = valueDown
  }
  const perps = s.filter((p: Perpetual) => {
    if (!p.trend) {
      return
    }
    if (params.slug.includes('-asc-')) {
      return p.trend.up > 0
    }
    if (params.slug.includes('-desc-')) {
      return p.trend.up < 0
    }
    return true
  }).sort(sort)

  return (
    <DefaultLayout>
      <main className="flex flex-wrap gap-2">
        {perps && perps.map((p: Perpetual, i: number) => {
          return (
            <Perp key={p.symbol} p={p} idx={i} days={params.slug.includes('dayspath') || params.slug.includes('dayschildpath')} />
          )})}
      </main>
    </DefaultLayout>
  );
}

export default Index

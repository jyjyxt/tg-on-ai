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

const Index = async ({ params }: { params: { slug: string } }) => {
  const client = initClient()

  const arr = params.slug.split('-')

  const s: Perpetual[] = await client.perpetuals(arr[0]);
  let sort = lowUp
  if (params.slug.includes('low-down')) {
    sort = lowDown
  } else if (params.slug.includes('high-up')) {
    sort = highUp
  } else if (params.slug.includes('high-down')) {
    sort = highDown
  } else if (params.slug.includes('rate-up')) {
    sort = rateUp
  } else if (params.slug.includes('rate-down')) {
    sort = rateDown
  }
  const perps = s.filter((p: Perpetual) => p.trend != null).sort(sort)

  return (
    <DefaultLayout>
      <main className="flex flex-wrap gap-2">
        {perps && perps.map((p: Perpetual) => {
          return <Perp key={p.symbol} p={p} />
        })}
      </main>
    </DefaultLayout>
  );
}

export default Index

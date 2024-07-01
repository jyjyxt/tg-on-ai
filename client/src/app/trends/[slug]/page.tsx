import Link from 'next/link';
import { Button } from "flowbite-react";
import { LiaArrowCircleDownSolid, LiaArrowCircleUpSolid } from "react-icons/lia";
import { Perpetual } from '@/http/types'
import { initClient } from '@/http/request'
import Perp from '@/components/Perpetual'
import Switcher from '@/components/Switcher'

const lowUp = (a: Perpetual, b: Perpetual) => b.trend!.up - a.trend!.up
const lowDown = (a: Perpetual, b: Perpetual) => a.trend!.up - b.trend!.up
const highUp = (a: Perpetual, b: Perpetual) => b.trend!.down! - a.trend!.down!
const highDown = (a: Perpetual, b: Perpetual) => a.trend!.down! - b.trend!.down!

interface Period {
  name: string;
  link: string;
  color: string;
}

const periods: Period[][] = [
  [
    {
      name: '3DaysLow',
      link: 'days3-low-up',
      color: 'blue',
    },
    {
      name: '3DaysHigh',
      link: 'days3-high-up',
      color: 'blue',
    },
  ],
  [
    {
      name: '7DaysLow',
      link: 'days7-low-up',
      color: 'success',
    },
    {
      name: '7DaysHigh',
      link: 'days7-high-up',
      color: 'success',
    },
  ],
  [
    {
      name: '15DaysLow',
      link: 'days15-low-up',
      color: 'purple',
    },
    {
      name: '15DaysHigh',
      link: 'days15-high-up',
      color: 'purple',
    },
  ],
  [
    {
      name: '30DaysLow',
      link: 'days30-low-up',
      color: 'warning',
    },
    {
      name: '30DaysHigh',
      link: 'days30-high-up',
      color: 'warning',
    },
  ],
]

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
  }
  const perps = s.filter((p: Perpetual) => p.trend != null).sort(sort)

  return (
    <main className="p-2">
      <div className="flex flex-wrap gap-2 mb-4">
        <Button as={Link} href="/trends/dayspath-low-up">Days</Button>
        { periods.map((p, i) => {
          return <Button.Group key={i}>
            {p.map((pp) => {
              return <Button key={pp.link} as={Link} href={`/trends/${pp.link}`} color={pp.color}>{pp.name}</Button>
            })}
          </Button.Group>
        })}
      </div>
      <div className="mb-4">
        <Button.Group>
          <Button as={Link} href={params.slug.replaceAll('down', 'up')} color="gray">
            <LiaArrowCircleUpSolid className="mr-3 h-5 w-5" /> UP
          </Button>
          <Button as={Link} href={params.slug.replaceAll('up', 'down')} color="gray">
            <LiaArrowCircleDownSolid className="mr-3 h-5 w-5" /> Down
          </Button>
        </Button.Group>
      </div>
      <div className="flex flex-wrap gap-2">
        {perps && perps.map((p: Perpetual) => {
          return <Perp key={p.symbol} p={p} />
        })}
      </div>
    </main>
  );
}

export default Index

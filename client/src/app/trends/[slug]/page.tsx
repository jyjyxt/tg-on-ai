import { Button } from "flowbite-react";
import { HiAdjustments, HiCloudDownload, HiUserCircle } from "react-icons/hi";
import { Perpetual } from '@/http/types'
import { initClient } from '@/http/request'
import Perp from '@/components/Perpetual'
import Switcher from '@/components/Switcher'

const lowUp = (a: Perpetual, b: Perpetual) => b.trend!.up - a.trend!.up
const lowDown = (a: Perpetual, b: Perpetual) => a.trend!.up - b.trend!.up
const highUp = (a: Perpetual, b: Perpetual) => b.trend!.down - a.trend!.down
const highDown = (a: Perpetual, b: Perpetual) => a.trend!.down - b.trend!.down

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
      <Switcher />
      <div> {params.slug} </div>
      <div>
        <Button.Group outline>
          <Button gradientDuoTone="cyanToBlue">
            <HiUserCircle className="mr-3 h-4 w-4" />
            UP
          </Button>
          <Button gradientDuoTone="cyanToBlue">
            <HiAdjustments className="mr-3 h-4 w-4" />
            Down
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

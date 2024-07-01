import { Button } from "flowbite-react";
import { HiAdjustments, HiCloudDownload, HiUserCircle } from "react-icons/hi";
import { Perpetual } from '@/http/types'
import { initClient } from '@/http/request'
import Perp from '@/components/Perpetual'
import Switcher from '@/components/Switcher'

const Index = async () => {
  const client = initClient()
  const s: Perpetual[] = await client.perpetuals('days30');
  const up = (a: Perpetual, b: Perpetual) => b.trend!.up - a.trend!.up
  const perps = s.filter((p: Perpetual) => p.trend != null).sort(up)

  return (
    <main className="p-2">
      <Switcher />
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

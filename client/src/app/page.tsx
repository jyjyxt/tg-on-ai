import { Button } from "flowbite-react";
import { HiAdjustments, HiCloudDownload, HiUserCircle } from "react-icons/hi";
import { Perpetual } from '@/http/types'
import { initClient } from '@/http/request'
import Perp from '@/components/Perpetual'
import Switcher from '@/components/Switcher'
import Header from '@/components/Header'

const Index = async () => {
  const client = initClient()
  const s: Perpetual[] = await client.perpetuals('days3');
  const up = (a: Perpetual, b: Perpetual) => b.trend!.up - a.trend!.up
  const perps = s.filter((p: Perpetual) => p.trend != null).sort(up)

  return (
    <main className="p-2">
      <Header slug="days3-low-up" />
      <div className="flex flex-wrap gap-2">
        {perps && perps.map((p: Perpetual) => {
          return <Perp key={p.symbol} p={p} />
        })}
      </div>
    </main>
  );
}

export default Index

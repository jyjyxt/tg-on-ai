import { Button } from "flowbite-react";
import { HiAdjustments, HiCloudDownload, HiUserCircle } from "react-icons/hi";
import { Perpetual } from '@/apis/types'
import { initClient } from '@/apis/request'
import Perp from '@/components/Perpetual'
import Switcher from '@/components/Switcher'
import DefaultLayout from '@/components/Layout'

const Index = async () => {
  const client = initClient()
  const s: Perpetual[] = await client.perpetuals('days3');
  const up = (a: Perpetual, b: Perpetual) => b.trend!.up - a.trend!.up
  const perps = s.filter((p: Perpetual) => p.trend != null).sort(up)

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

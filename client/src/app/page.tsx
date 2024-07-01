import { Perpetual } from '@/http/types'
import { initClient } from '@/http/request'
import Perp from '@/components/Perpetual'
import Switcher from '@/components/Switcher'

const Index = async () => {
  const client = initClient()
  const s: Perpetual[] = await client.perpetuals('days30');

  return (
    <main className="p-2">
      <Switcher />
      <div className="flex flex-wrap gap-2">
        {s && s.filter((p: Perpetual) => p.trend != null).sort((a: Perpetual, b: Perpetual) => b.trend!.up - a.trend!.up).map((p: Perpetual) => {
          return <Perp key={p.symbol} p={p} />
        })}
      </div>
    </main>
  );
}

export default Index

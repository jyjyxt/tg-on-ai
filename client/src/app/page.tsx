import { Perpetual } from '@/http/types'
import { initClient } from '@/http/request'
import Perp from '@/components/Perpetual'

const Index = async () => {
  const client = initClient()
  const s: Perpetual[] = await client.perpetuals('dayspath');

  return (
    <main className="p-2">
      {s && s.map((p: Perpetual) => {
        return <Perp key={p.symbol} p={p} />
      })}
    </main>
  );
}

export default Index

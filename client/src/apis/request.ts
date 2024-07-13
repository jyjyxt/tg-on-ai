import { Perpetual } from './types'

const base = 'https://trading-api.gmgn.top';

interface Resp {
  data: any
}

export const initClient = () => {
  const ins = (url: string) => {
    return fetch(`${base}${url}`, { cache: 'no-store' }).then((res) => res.json()).then((data: Resp) => data.data)
  }

  return {
    perpetuals: (trend: string): Promise<Perpetual[]> => ins(`/perpetuals?trend=${trend}`),
  }
}

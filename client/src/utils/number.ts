import BigNumber from 'bignumber.js'

export function formatNumber(n: BigNumber.Value): string {
  return new BigNumber(n).toFormat({ groupSeparator: ',', groupSize: 3 })
}

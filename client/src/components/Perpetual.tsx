import { Card, Badge, List } from "flowbite-react";
import BigNumber from 'bignumber.js'
import { Perpetual, Trend } from '@/apis/types';
import { formatNumber } from '@/utils/number';
import { formatDateFromNow } from '@/utils/date'
import { SiBinance } from "react-icons/si";

interface prop {
  p: Perpetual
  idx: number
  days?: boolean
}

const fields: string[][] = [
  ['Price', 'now'],
  ['High', 'high'],
  ['Low', 'low'],
]

const Index = ({ p, idx, days }: prop) => {
  const t = p.trend as Trend

  return (
    <Card className="w-80 flex-grow">
      <div className="flex items-center justify-between">
        <h5 className="text-xl font-bold tracking-tight text-gray-900 dark:text-white">
          { !days && <><span className="text-green-600 dark:text-green-300">{t.up}%</span> / <span className="text-red-600 dark:text-red-300">{(t.down || 0) * -1}%</span></> }
          { days && <span className={t.up > 0 ? 'text-green-600 dark:text-green-300' : 'text-red-600 dark:text-red-300'}>{t.up} Days</span> }
        </h5>
        <span className="text-sm"> {p.symbol} / No. {idx} </span>
      </div>
      <div className="flex-1">
        <div className="flex flex-wrap gap-2">
          { fields && fields.map((f) => {
            return (
              <Badge key={f[0]} color="info">{f[0]}: ${t[f[1] as keyof Trend]}</Badge>
            )
          })}
          <Badge color="indigo">Funding Rate: {BigNumber(p.last_funding_rate).times(100).toFormat()}%</Badge>
          <Badge color="purple">Interest Value: {formatNumber(Math.floor(p.sum_open_interest_value))}</Badge>
          {p.categories && <Badge color="pink">{p.categories}</Badge>}
        </div>
      </div>
      <div className="text-gray-500 dark:text-gray-400 text-sm flex items-center">
        <div className="flex-1 flex">
          <a href={`https://www.binance.com/futures/${p.symbol}`} target="_blank" rel="noopener noreferrer">
            <SiBinance fill="#F0B90B" />
          </a>
          { p.coingecko && (
            <a href={`https://www.coingecko.com/coins/${p.coingecko}`} target="_blank" rel="noopener noreferrer">
              <CGC />
            </a>
          )}
        </div>
        {formatDateFromNow(new Date(p.updated_at))}
      </div>
    </Card>
  )
}

const CGC = () => (
  <svg width="16" height="16" viewBox="0 0 96 96" fill="none" xmlns="http://www.w3.org/2000/svg">
    <g clip-path="url(#clip0_238_5139)">
      <path d="M95.3201 47.7848C95.437 74.2931 74.1955 95.8818 47.8757 95.9995C21.5525 96.1172 0.120677 74.7236 0.00047986 48.2152C-0.116378 21.7035 21.1251 0.118179 47.445 0.000483373C73.7681 -0.117212 95.2 21.2731 95.3168 47.7848H95.3201Z" fill="#8DC63F"/>
      <path d="M91.7276 47.7979C91.8378 72.3122 72.1955 92.2733 47.859 92.3843C23.519 92.4953 3.69988 72.7124 3.58969 48.1981C3.47951 23.6839 23.1217 3.72273 47.4616 3.61177C71.7982 3.50416 91.6174 23.2837 91.7276 47.7979Z" fill="#F9E988"/>
      <path d="M48.4633 6.48732C51.2411 5.93919 54.1259 6.01317 56.9372 6.47051C59.7484 6.94465 62.5196 7.81896 65.0271 9.19767C67.5379 10.5898 69.7114 12.4763 71.875 14.2485C74.0319 16.0408 76.1887 17.8298 78.2354 19.797C80.3021 21.744 82.1552 23.9533 83.6576 26.3946C85.1868 28.8158 86.519 31.3849 87.4606 34.1054C89.3002 39.5496 89.9379 45.4176 89.0231 51.0131H88.7393C87.8178 45.4646 86.5658 40.1818 84.6125 35.1512C83.661 32.6359 82.6326 30.1407 81.3539 27.7599C80.0484 25.3993 78.6528 23.0689 77.0668 20.8293C75.4642 18.6099 73.5611 16.5688 71.3174 14.9849C69.0804 13.3809 66.5062 12.3385 64.0188 11.3498C61.5247 10.3477 59.0406 9.38262 56.453 8.63274C53.8688 7.86267 51.2211 7.31791 48.4633 6.76979V6.48396V6.48732Z" fill="white"/>
      <path d="M70.0087 32.1144C66.8101 31.183 63.498 29.8581 60.1392 28.5231C59.9455 27.6756 59.201 26.6198 57.6918 25.3251C55.4982 23.4083 51.3781 23.4588 47.819 24.3062C43.8892 23.3747 40.0062 23.0418 36.2801 23.943C5.81008 32.4003 23.0851 53.0239 11.8967 73.7618C13.4893 77.1615 30.6475 97.0083 55.4749 91.6817C55.4749 91.6817 46.9843 71.1321 66.1457 61.2659C81.6879 53.266 92.9163 38.4095 70.0054 32.1111L70.0087 32.1144Z" fill="#8BC53F"/>
      <path d="M73.7681 45.6293C73.7681 46.6549 72.9468 47.4922 71.9284 47.4956C70.9101 47.5023 70.0787 46.6717 70.0721 45.6427C70.0654 44.6171 70.8901 43.7798 71.9117 43.7764C72.9301 43.7697 73.7614 44.6003 73.7681 45.6259V45.6293Z" fill="white"/>
      <path d="M47.819 24.3101C50.0393 24.4681 58.0658 27.1079 60.1358 28.5236C58.4197 23.4795 52.6001 22.8103 47.819 24.3101Z" fill="#009345"/>
      <path d="M49.9291 37.0607C49.9291 41.8022 46.1128 45.6424 41.4085 45.6424C36.7041 45.6424 32.8878 41.8022 32.8878 37.0607C32.8878 32.3193 36.7041 28.4824 41.4085 28.4824C46.1128 28.4824 49.9291 32.3227 49.9291 37.0607Z" fill="white"/>
      <path d="M47.4049 37.1416C47.4049 40.474 44.7205 43.1776 41.4117 43.1776C38.103 43.1776 35.4186 40.4774 35.4186 37.1416C35.4186 33.8057 38.103 31.1055 41.4117 31.1055C44.7205 31.1055 47.4049 33.8091 47.4049 37.1416Z" fill="#58595B"/>
      <path d="M80.6726 49.4091C73.7713 54.3086 65.9151 58.0244 54.7802 58.0244C49.5683 58.0244 48.5099 52.4457 45.0643 55.1796C43.2847 56.5919 37.0144 59.7495 32.0362 59.5108C27.0146 59.2687 18.9982 56.3296 16.7445 45.6328C15.853 56.3296 15.3989 64.2119 11.4091 73.2441C19.3521 86.0527 38.2865 95.9324 55.4747 91.6853C53.6283 78.6951 64.9001 65.9739 71.2505 59.4637C73.6545 56.9988 78.262 52.9736 80.6726 49.4091Z" fill="#8BC53F"/>
      <path d="M80.4024 49.7321C78.2588 51.6993 75.708 53.1587 73.1104 54.4432C70.4861 55.6908 67.7516 56.7063 64.927 57.4428C62.1124 58.1759 59.1709 58.7273 56.1927 58.4583C53.2679 58.1994 50.1761 57.1637 48.2029 54.9207L48.2964 54.8131C50.7304 56.3936 53.5049 56.9485 56.2795 57.0292C59.054 57.1032 61.882 56.8947 64.6699 56.3264C67.4511 55.748 70.1823 54.8838 72.8233 53.7875C75.4576 52.6913 78.0619 51.4235 80.3089 49.6211L80.399 49.7287L80.4024 49.7321Z" fill="#58595B"/>
    </g>
    <defs>
      <clipPath id="clip0_238_5139">
        <rect width="96" height="96" fill="white"/>
      </clipPath>
    </defs>
  </svg>
)

export default Index

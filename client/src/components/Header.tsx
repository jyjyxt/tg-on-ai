import Link from 'next/link';
import { Button } from "flowbite-react";
import { LiaArrowCircleDownSolid, LiaArrowCircleUpSolid } from "react-icons/lia";
import Switcher from '@/components/Switcher'

interface Period {
  name: string;
  link: string;
  color: string;
}

const periods: Period[][] = [
  [
    {
      name: 'Days',
      link: 'days-low-up',
      color: 'info',
    },
    {
      name: 'Rate',
      link: 'days3-rate-up',
      color: 'info',
    }
  ],
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

interface Props {
  slug: string
}

const Index = ({ slug }: Props) => {
  return (
    <>
      <div className="flex flex-wrap gap-2 mb-4">
        { periods.map((p, i) => {
          return <Button.Group key={i}>
            {p.map((pp) => {
              return <Button key={pp.link} as={Link} href={`/trends/${pp.link}`} color={pp.color}>{pp.name}</Button>
            })}
          </Button.Group>
        })}
      </div>
      <div className="mb-4">
        <Button.Group className="mr-4">
          <Button as={Link} href={`/trends/${slug.replaceAll('down', 'up')}`} color="gray">
            <LiaArrowCircleUpSolid className="mr-3 h-5 w-5" /> UP
          </Button>
          <Button as={Link} href={`/trends/${slug.replaceAll('up', 'down')}`} color="gray">
            <LiaArrowCircleDownSolid className="mr-3 h-5 w-5" /> Down
          </Button>
        </Button.Group>
        <Switcher />
      </div>
    </>
  )
}

export default Index

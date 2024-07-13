import Link from 'next/link';
import { usePathname } from "next/navigation";
import { Button } from "flowbite-react";
import { LiaArrowCircleDownSolid, LiaArrowCircleUpSolid } from "react-icons/lia";
import Switcher from '@/components/Switcher'

interface Props {
  slug: string
  sidebarOpen: boolean;
  setSidebarOpen: (arg: boolean) => void;
}

const Index = ({ slug, sidebarOpen, setSidebarOpen }: Props) => {
  const pathname = usePathname();

  return (
    <div className="mb-4">
      <button className="inline-flex items-center p-2 mt-2 ms-3 mr-4 text-sm text-gray-500 rounded-lg sm:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
        onClick={(e) => {
          e.stopPropagation();
          setSidebarOpen(!sidebarOpen);
        }}
      >
        <svg className="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
          <path stroke="currentColor" strokeLinecap="round" strokeWidth="2" d="M5 7h14M5 12h14M5 17h10"/>
        </svg>
      </button>
      <Button.Group className="mr-4">
        <Button as={Link} href={`/trends/${slug.replaceAll('down', 'up')}`} color={`${pathname.includes('up') ? 'success' : 'gray'}`}>
          <LiaArrowCircleUpSolid className="mr-3 h-5 w-5" /> UP
        </Button>
        <Button as={Link} href={`/trends/${slug.replaceAll('up', 'down')}`} color={`${pathname.includes('down') ? 'success' : 'gray'}`}>
          <LiaArrowCircleDownSolid className="mr-3 h-5 w-5" /> Down
        </Button>
      </Button.Group>
      <Switcher />
    </div>
  )
}

export default Index

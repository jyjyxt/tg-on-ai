"use client";

import React, { useEffect, useRef, useState } from "react";
import Link from 'next/link';
import { usePathname } from "next/navigation";
import { Sidebar } from "flowbite-react";
import { HiArrowSmRight, HiChartPie, HiInbox, HiShoppingBag, HiTable, HiUser, HiViewBoards } from "react-icons/hi";

interface Props {
  sidebarOpen: boolean;
  setSidebarOpen: (arg: boolean) => void;
}

interface Period {
  name: string;
  link: string;
  color: string;
}

const periods: Period[][] = [
  [
    {
      name: 'Days Up',
      link: 'dayspath-asc-up',
      color: 'info',
    },
    {
      name: 'Days Down',
      link: 'dayspath-desc-down',
      color: 'info',
    },
    {
      name: 'Days Child Up',
      link: 'dayschildpath-asc-up',
      color: 'info',
    },
    {
      name: 'Days Child Down',
      link: 'dayschildpath-desc-down',
      color: 'info',
    },
    {
      name: 'Runding Rate',
      link: 'days3-rate-up',
      color: 'info',
    },
    {
      name: 'Interest Value',
      link: 'days3-value-up',
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

const Index = ({ sidebarOpen, setSidebarOpen }: Props) => {
  const sidebar = useRef<any>(null);
  const pathname = usePathname();

  // close on click outside
  useEffect(() => {
    const clickHandler = ({ target }: MouseEvent) => {
      if (!sidebar.current) return;
      if (
        !sidebarOpen ||
        sidebar.current.contains(target)
      )
        return;
      setSidebarOpen(false);
    };
    document.addEventListener("click", clickHandler);
    return () => document.removeEventListener("click", clickHandler);
  });

  // close if the esc key is pressed
  useEffect(() => {
    const keyHandler = ({ key }: KeyboardEvent) => {
      if (!sidebarOpen || key !== "Escape") return;
      setSidebarOpen(false);
    };
    document.addEventListener("keydown", keyHandler);
    return () => document.removeEventListener("keydown", keyHandler);
  });

  return (
    <aside ref={sidebar} className={`fixed top-0 left-0 z-40 w-64 h-screen transition-transform ${sidebarOpen ? "translate-x-0" : "-translate-x-full"} sm:translate-x-0`} aria-label="Sidebar">
      <Sidebar>
        <Sidebar.Items>
          { periods.map((p, i) => {
            return <Sidebar.ItemGroup key={i}>
              {p.map((pp) => {
                return <Sidebar.Item active={pathname.replaceAll('down', 'up').includes(pp.link)} key={pp.link} as={Link} href={`/trends/${pp.link}`} color={pp.color}>{pp.name}</Sidebar.Item>
              })}
            </Sidebar.ItemGroup>
          })}
        </Sidebar.Items>
      </Sidebar>
    </aside>
  );
}

export default Index

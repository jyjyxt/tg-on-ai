"use client";
import React, { useState, ReactNode } from "react"
import { usePathname } from "next/navigation";
import SideBar from '@/components/Aside'
import Header from '@/components/Header'

export default function DefaultLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  let pathname = usePathname();
  if (pathname == '/') {
    pathname = 'days3-low-up' 
  }

  return (
    <div className="flex h-screen overflow-hidden">
      <SideBar sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />
      <div className="p-2 sm:ml-64 overflow-y-auto overflow-x-h">
        <Header slug={pathname.replaceAll('/trends/', '')} sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen}  />
        {children}
      </div>
    </div>
  );
}

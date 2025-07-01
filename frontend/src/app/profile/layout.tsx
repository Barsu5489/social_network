'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useUser } from '@/contexts/user-context';
import { Header } from "@/components/layout/header";
import { Sidebar } from "@/components/layout/sidebar";
import { PageLoader } from '@/components/page-loader';

export default function ProfileLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user, isLoading } = useUser();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !user) {
      router.replace('/');
    }
  }, [user, isLoading, router]);

  if (isLoading) {
    return <PageLoader />;
  }

  if (user) {
    return (
      <div className="grid min-h-screen w-full lg:grid-cols-[280px_1fr]">
        <Sidebar />
        <div className="flex flex-col bg-background/80 backdrop-blur-sm border-r border-l">
          <Header />
          {children}
        </div>
      </div>
    );
  }

  return <PageLoader />;
}

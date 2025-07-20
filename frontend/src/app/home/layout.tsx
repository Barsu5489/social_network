'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useUser } from '@/contexts/user-context';
import { Header } from "@/components/layout/header";
import { Sidebar } from "@/components/layout/sidebar";
import { PageLoader } from '@/components/page-loader';

export default function HomeLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user, isLoading } = useUser();
  const router = useRouter();

  useEffect(() => {
    // Only redirect if we're done loading and there's no user
    if (!isLoading && !user) {
      console.log('No user found, redirecting to login')
      router.replace('/');
    }
  }, [user, isLoading, router]);

  // Show loader while checking authentication
  if (isLoading) {
    return <PageLoader />;
  }

  // Show loader while redirecting
  if (!user) {
    return <PageLoader />;
  }

  // User is authenticated, show the layout
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

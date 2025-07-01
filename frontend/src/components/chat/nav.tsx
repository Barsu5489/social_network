
'use client';

import Link from 'next/link';
import { LucideIcon } from 'lucide-react';

import { cn } from '@/lib/utils';
import { buttonVariants } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { usePathname } from 'next/navigation';
import { Skeleton } from '../ui/skeleton';

interface NavProps {
  isCollapsed: boolean;
  links: {
    title: string;
    label?: string;
    icon: LucideIcon;
    variant: 'default' | 'ghost';
    href: string;
  }[];
  title: string;
  isLoading: boolean;
}

export function Nav({ links, isCollapsed, title, isLoading }: NavProps) {
  const pathname = usePathname();

  if (isCollapsed) {
    return (
        <div
            data-collapsed={isCollapsed}
            className="group flex flex-col gap-4 py-2 data-[collapsed=true]:py-2"
        >
            <nav className="grid gap-1 px-2 group-[[data-collapsed=true]]:justify-center group-[[data-collapsed=true]]:px-2">
                 {links.map((link, index) => (
                    <Tooltip key={index} delayDuration={0}>
                    <TooltipTrigger asChild>
                        <Link
                        href={link.href}
                        className={cn(
                            buttonVariants({ variant: 'ghost', size: 'icon' }),
                            'h-9 w-9',
                             pathname === link.href && 'bg-muted text-primary'
                        )}
                        >
                        <link.icon className="h-4 w-4" />
                        <span className="sr-only">{link.title}</span>
                        </Link>
                    </TooltipTrigger>
                    <TooltipContent side="right" className="flex items-center gap-4">
                        {link.title}
                        {link.label && (
                        <span className="ml-auto text-muted-foreground">
                            {link.label}
                        </span>
                        )}
                    </TooltipContent>
                    </Tooltip>
                ))}
            </nav>
        </div>
    )
  }

  return (
    <div
      data-collapsed={isCollapsed}
      className="group flex flex-col gap-4 py-2 data-[collapsed=true]:py-2"
    >
      <h2 className="px-4 text-lg font-semibold tracking-tight">{title}</h2>
      <nav className="grid gap-1 px-2 group-[[data-collapsed=true]]:justify-center group-[[data-collapsed=true]]:px-2">
        {isLoading ? (
            Array.from({length: 3}).map((_, i) => (
                <div key={i} className="flex items-center gap-2 px-2 py-1.5">
                    <Skeleton className="h-4 w-4" />
                    <Skeleton className="h-4 w-24" />
                </div>
            ))
        ) : links.length > 0 ? (
          links.map((link, index) => (
            <Link
              key={index}
              href={link.href}
              className={cn(
                buttonVariants({ variant: 'ghost', size: 'sm' }),
                pathname === link.href && 'bg-muted text-primary',
                'justify-start'
              )}
            >
              <link.icon className="mr-2 h-4 w-4" />
              {link.title}
              {link.label && (
                <span
                  className={cn(
                    'ml-auto',
                    pathname === link.href && 'text-background dark:text-white'
                  )}
                >
                  {link.label}
                </span>
              )}
            </Link>
          ))
        ) : (
          <p className="px-2 text-sm text-muted-foreground">No chats found.</p>
        )}
      </nav>
    </div>
  );
}


'use client';

import * as React from 'react';
import { cn } from '@/lib/utils';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue, } from '@/components/ui/select';
import { Skeleton } from '../ui/skeleton';

interface AccountSwitcherProps {
  isCollapsed: boolean;
  user: { id: string; email: string; first_name: string; last_name: string } | null;
  isLoading: boolean;
}

export function AccountSwitcher({ isCollapsed, user, isLoading }: AccountSwitcherProps) {
  const [selectedAccount, setSelectedAccount] = React.useState<string>(user?.id || '');

  React.useEffect(() => {
    if(user?.id) {
        setSelectedAccount(user.id);
    }
  }, [user]);

  if (isLoading) {
    return (
        <div className="flex items-center gap-2 p-2">
            <Skeleton className="h-8 w-8 rounded-full" />
            {!isCollapsed && <Skeleton className="h-6 w-32" />}
        </div>
    )
  }

  return (
    <Select value={selectedAccount} onValueChange={setSelectedAccount}>
      <SelectTrigger
        className={cn(
          'flex items-center gap-2 [&>span]:line-clamp-1 [&>span]:flex [&>span]:w-full [&>span]:items-center [&>span]:gap-1 [&>span]:truncate [&_svg]:h-4 [&_svg]:w-4 [&_svg]:shrink-0',
          isCollapsed &&
            'flex h-9 w-9 shrink-0 items-center justify-center p-0 [&>span]:w-auto [&>svg]:hidden',
          'border-none shadow-none focus:ring-0'
        )}
        aria-label="Select account"
      >
        <SelectValue placeholder="Select an account">
          {user && (
            <div className="flex items-center gap-3">
              <div className="h-8 w-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
                {user.first_name?.[0] || '?'}{user.last_name?.[0] || '?'}
              </div>
              {!isCollapsed && (
                <span className="ml-2">
                  {user.first_name || 'Unknown'} {user.last_name || 'User'}
                </span>
              )}
            </div>
          )}
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        {user && (
          <SelectItem value={user.id}>
            <div className="flex items-center gap-3">
              <div className="h-6 w-6 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold text-xs">
                {user.first_name?.[0] || '?'}{user.last_name?.[0] || '?'}
              </div>
              <div className="flex flex-col">
                <span>{user.first_name || 'Unknown'} {user.last_name || 'User'}</span>
                <span className="text-xs text-muted-foreground">{user.email}</span>
              </div>
            </div>
          </SelectItem>
        )}
      </SelectContent>
    </Select>
  );
}

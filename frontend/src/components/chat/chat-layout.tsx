'use client';

import * as React from 'react';
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';
import { Separator } from '@/components/ui/separator';
import { TooltipProvider } from '@/components/ui/tooltip';
import { cn } from '@/lib/utils';
import { API_BASE_URL } from '@/lib/config';
import { useToast } from '@/hooks/use-toast';
import { AccountSwitcher } from './account-switcher';
import { Nav } from './nav';
import { MessageCircle, User, Users } from 'lucide-react';
import { useUser } from '@/contexts/user-context';
import { Skeleton } from '../ui/skeleton';

interface ChatLayoutProps {
  defaultLayout: number[] | undefined;
  navCollapsedSize: number;
  children?: React.ReactNode;
}

interface Chat {
    id: string;
    type: 'direct' | 'group';
    name: string;
    avatar_url: string;
}

export function ChatLayout({ defaultLayout = [265, 440, 655], navCollapsedSize, children }: ChatLayoutProps) {
    const { user, isLoading: isUserLoading } = useUser();
    const [isCollapsed, setIsCollapsed] = React.useState(false);
    const [chats, setChats] = React.useState<Chat[]>([]);
    const [isLoading, setIsLoading] = React.useState(true);
    const { toast } = useToast();

    React.useEffect(() => {
        const fetchChats = async () => {
            setIsLoading(true);
             try {
                const response = await fetch(`${API_BASE_URL}/api/chats`, {
                    credentials: 'include'
                });
                if (response.ok) {
                    const data = await response.json();
                    setChats(data.chats || []);
                } else {
                    toast({ variant: 'destructive', title: 'Failed to load chats.' });
                }
            } catch (error) {
                toast({ variant: 'destructive', title: 'Network error loading chats.' });
            } finally {
                setIsLoading(false);
            }
        }
        fetchChats();
    }, [toast]);

    const directChats = chats.filter(chat => chat.type === 'direct');
    const groupChats = chats.filter(chat => chat.type === 'group');

    return (
      <TooltipProvider delayDuration={0}>
        <ResizablePanelGroup
          direction="horizontal"
          onLayout={(sizes: number[]) => {
            document.cookie = `react-resizable-panels:layout=${JSON.stringify(sizes)}`;
          }}
          className="h-full max-h-screen items-stretch"
        >
          <ResizablePanel
            defaultSize={defaultLayout[0]}
            collapsedSize={navCollapsedSize}
            collapsible={true}
            minSize={15}
            maxSize={20}
            onCollapse={() => setIsCollapsed(true)}
            onExpand={() => setIsCollapsed(false)}
            className={cn(isCollapsed && 'min-w-[50px] transition-all duration-300 ease-in-out')}
          >
            <AccountSwitcher
              isCollapsed={isCollapsed}
              user={user}
              isLoading={isUserLoading}
            />
            <Separator />
            <Nav
              isCollapsed={isCollapsed}
              title="Direct Messages"
              links={directChats.map(c => ({
                  title: c.name,
                  label: '',
                  icon: User,
                  variant: 'ghost',
                  href: `/chat/${c.id}`
              }))}
              isLoading={isLoading}
            />
            <Separator />
             <Nav
              isCollapsed={isCollapsed}
              title="Group Chats"
              links={groupChats.map(c => ({
                  title: c.name,
                  label: '',
                  icon: Users,
                  variant: 'ghost',
                  href: `/chat/${c.id}`
              }))}
              isLoading={isLoading}
            />
          </ResizablePanel>
          <ResizableHandle withHandle />
          <ResizablePanel defaultSize={defaultLayout[2]} minSize={30}>
            {children ? children : (
                <div className="flex h-full items-center justify-center p-6">
                    <div className="flex flex-col items-center gap-2 text-center">
                        <MessageCircle className="h-16 w-16 text-muted-foreground" />
                        <h2 className="text-2xl font-bold">Select a chat</h2>
                        <p className="text-muted-foreground">
                            Choose one of your conversations to start messaging.
                        </p>
                    </div>
                </div>
            )}
          </ResizablePanel>
        </ResizablePanelGroup>
      </TooltipProvider>
    );
}

// Add child components for ChatLayout
AccountSwitcher.displayName = "AccountSwitcher";
Nav.displayName = "Nav";

// These are placeholders to be filled out with actual chat components
// For now, this file just sets up the main layout
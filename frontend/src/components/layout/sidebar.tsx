'use client';

import { usePathname, useRouter } from 'next/navigation';
import Link from 'next/link';
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuPortal, DropdownMenuSeparator, DropdownMenuSub, DropdownMenuSubContent, DropdownMenuSubTrigger, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Bell, Calendar, Home, Settings, User, Users, Palette, LogOut, MessageSquare } from "lucide-react";
import { Logo } from "../logo";
import { cn } from '@/lib/utils';
import { useTheme } from 'next-themes';
import { API_BASE_URL } from '@/lib/config';
import { useUser } from '@/contexts/user-context';


const NavLink = ({ href, icon: Icon, children, badgeCount }: { href: string, icon: React.ElementType, children: React.ReactNode, badgeCount?: number }) => {
    const pathname = usePathname();
    const isActive = pathname.startsWith(href);

    return (
        <Link
            href={href}
            className={cn(
                'flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary',
                isActive ? 'bg-primary text-primary-foreground' : 'text-muted-foreground'
            )}
        >
            <Icon className="h-5 w-5" />
            {children}
            {badgeCount && badgeCount > 0 && (
                 <span className="ml-auto flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground text-xs">
                    {badgeCount}
                </span>
            )}
        </Link>
    );
};

export function Sidebar() {
    const { setTheme } = useTheme();
    const router = useRouter();
    const { user, setUser, isLoading } = useUser();

    const handleLogout = async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/api/logout`, {
                method: 'POST',
                credentials: 'include',
            });
            
            if (!response.ok) {
                console.error('Logout failed:', response.statusText);
            }
        } catch (error) {
            console.error('Logout error:', error);
        } finally {
            // Always clear user data regardless of API response
            setUser(null); // Clear user from context and localStorage
            
            // Clear any remaining session data
            localStorage.clear();
            sessionStorage.clear();
            
            // Force redirect to login page
            window.location.href = '/';
        }
    };

    if (isLoading) {
        return (
            <div className="hidden border-r bg-card lg:block">
                 <div className="flex h-full max-h-screen flex-col gap-2">
                    {/* Skeleton loader */}
                </div>
            </div>
        );
    }
    
    return (
        <div className="hidden border-r bg-card lg:block">
            <div className="flex h-full max-h-screen flex-col gap-2">
                <div className="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
                    <Logo />
                </div>
                <div className="flex-1">
                    <nav className="grid items-start px-2 text-sm font-medium lg:px-4 gap-1">
                        <NavLink href="/home" icon={Home}>Home</NavLink>
                        <NavLink href="/groups" icon={Users}>Groups</NavLink>
                        <NavLink href="/chat" icon={MessageSquare}>Chat</NavLink>
                        {/* <NavLink href="#" icon={Calendar}>Events</NavLink> */}
                        {/* Notifications are in header now */}
                    </nav>
                </div>
                <div className="mt-auto p-4 border-t">
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="ghost" className="flex items-center gap-3 w-full justify-start h-auto p-2">
                                <Avatar className="h-9 w-9">
                                    <AvatarImage src={user?.avatar_url || `https://i.pravatar.cc/40?u=${user?.id}`} alt={user?.first_name || 'User'} data-ai-hint="woman portrait" />
                                    <AvatarFallback>{user?.first_name?.[0]}{user?.last_name?.[0]}</AvatarFallback>
                                </Avatar>
                                <div className="text-left">
                                    <p className="font-semibold text-sm text-card-foreground">{user?.first_name} {user?.last_name}</p>
                                    <p className="text-xs text-muted-foreground">{user?.email}</p>
                                </div>
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent className="w-56" align="start" forceMount>
                            <DropdownMenuLabel className="font-normal">
                                <div className="flex flex-col space-y-1">
                                    <p className="text-sm font-medium leading-none">{user?.first_name} {user?.last_name}</p>
                                    <p className="text-xs leading-none text-muted-foreground">
                                        {user?.email}
                                    </p>
                                </div>
                            </DropdownMenuLabel>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem onClick={() => router.push(`/profile/${user?.id}`)}>
                                <User className="mr-2 h-4 w-4" />
                                <span>Profile</span>
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                                <Settings className="mr-2 h-4 w-4" />
                                <span>Settings</span>
                            </DropdownMenuItem>
                            <DropdownMenuSub>
                                <DropdownMenuSubTrigger>
                                    <Palette className="mr-2 h-4 w-4" />
                                    <span>Theme</span>
                                </DropdownMenuSubTrigger>
                                <DropdownMenuPortal>
                                    <DropdownMenuSubContent>
                                        <DropdownMenuItem onClick={() => setTheme('light')}>Light</DropdownMenuItem>
                                        <DropdownMenuItem onClick={() => setTheme('dark')}>Dark</DropdownMenuItem>
                                        <DropdownMenuItem onClick={() => setTheme('system')}>System</DropdownMenuItem>
                                    </DropdownMenuSubContent>
                                </DropdownMenuPortal>
                            </DropdownMenuSub>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem onClick={handleLogout} className="text-red-600 focus:text-red-600">
                                <LogOut className="mr-2 h-4 w-4" />
                                <span>Log out</span>
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </div>
        </div>
    );
}

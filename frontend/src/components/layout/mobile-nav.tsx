'use client';

import { usePathname } from 'next/navigation';
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Bell, Calendar, Home, Menu, Users, MessageSquare } from "lucide-react";
import Link from 'next/link';
import { Logo } from "../logo";
import { cn } from '@/lib/utils';

export function MobileNav() {
    const pathname = usePathname();

    const links = [
        { href: '/home', icon: Home, label: 'Home' },
        { href: '/groups', icon: Users, label: 'Groups' },
        { href: '/chat', icon: MessageSquare, label: 'Chat' },
        // { href: '#', icon: Calendar, label: 'Events' },
        // { href: '#', icon: Bell, label: 'Notifications' },
    ];


    return (
        <Sheet>
            <SheetTrigger asChild>
                <Button variant="outline" size="icon" className="shrink-0 lg:hidden">
                    <Menu className="h-5 w-5" />
                    <span className="sr-only">Toggle navigation menu</span>
                </Button>
            </SheetTrigger>
            <SheetContent side="left" className="flex flex-col">
                <nav className="grid gap-2 text-lg font-medium">
                    <div className="flex items-center gap-2 text-lg font-semibold mb-4">
                        <Logo />
                    </div>
                    {links.map((link) => (
                        <Link 
                            key={link.href}
                            href={link.href} 
                            className={cn("mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 hover:text-foreground", pathname.startsWith(link.href) ? "bg-muted text-foreground" : "text-muted-foreground")}
                        >
                            <link.icon className="h-5 w-5" />
                            {link.label}
                        </Link>
                    ))}
                </nav>
            </SheetContent>
        </Sheet>
    )
}

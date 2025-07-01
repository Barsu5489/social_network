'use client';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Button } from "@/components/ui/button"
import { Bell, UserPlus, Users } from "lucide-react"
import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";

const notifications = [
    {
        id: '1',
        type: 'follow_request',
        user: { name: 'John Doe', avatar: 'https://i.pravatar.cc/40?u=john' },
        message: 'sent you a follow request.',
        time: '5m ago',
        href: '#'
    },
    {
        id: '2',
        type: 'group_invite',
        user: { name: 'Jane Smith', avatar: 'https://i.pravatar.cc/40?u=jane' },
        message: 'invited you to join Design Enthusiasts.',
        time: '1h ago',
        href: '#'
    },
];

const icons: { [key: string]: React.ElementType } = {
    follow_request: UserPlus,
    group_invite: Users
}

export function NotificationBell() {
    const hasUnread = notifications.length > 0;

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                 <Button variant="ghost" size="icon" className="relative">
                    <Bell className="h-5 w-5" />
                    {hasUnread && (
                        <span className="absolute top-1 right-1 flex h-2 w-2">
                            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                            <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
                        </span>
                    )}
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-80" align="end">
                <DropdownMenuLabel>Notifications</DropdownMenuLabel>
                <DropdownMenuSeparator />
                {notifications.length > 0 ? (
                    notifications.map((notif) => {
                        const Icon = icons[notif.type] || Bell;
                        return (
                             <DropdownMenuItem key={notif.id} asChild>
                                <Link href={notif.href} className="flex items-start gap-3">
                                    <Avatar className="h-8 w-8 mt-1">
                                        <AvatarImage src={notif.user.avatar} />
                                        <AvatarFallback>{notif.user.name[0]}</AvatarFallback>
                                    </Avatar>
                                    <div className="text-sm">
                                        <p>
                                            <span className="font-semibold">{notif.user.name}</span>
                                            {' '}
                                            {notif.message}
                                        </p>
                                        <p className="text-xs text-muted-foreground">{notif.time}</p>
                                    </div>
                                </Link>
                            </DropdownMenuItem>
                        )
                    })
                ) : (
                    <DropdownMenuItem disabled>No new notifications</DropdownMenuItem>
                )}
            </DropdownMenuContent>
        </DropdownMenu>
    )
}

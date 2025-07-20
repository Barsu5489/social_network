'use client'

import { useState, useEffect } from "react"
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
import Link from "next/link"
import { API_BASE_URL } from '@/lib/config';
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar"

type Notification = {
  id: string
  type: string
  user: { name: string; avatar: string }
  message: string
  time: string
  href: string
}

const icons: { [key: string]: React.ElementType } = {
  follow_request: UserPlus,
  group_invite: Users
}

export function NotificationBell() {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading] = useState(true)

  const markAsRead = async (id: string) => {
    console.log('DEBUG: Marking notification as read:', id)
    try {
      const res = await fetch(`${API_BASE_URL}/api/notifications/${id}`, {
        method: 'PUT',
        credentials: 'include',
      });
      console.log('DEBUG: Mark as read response status:', res.status)
      
      if (!res.ok) {
        throw new Error('Failed to mark notification as read');
      }
      
      console.log('SUCCESS: Notification marked as read:', id)
      // Update the notifications state to remove the read notification
      setNotifications((prevNotifications) =>
        (prevNotifications || []).filter((notif) => notif.id !== id)
      );
    } catch (err) {
      console.error('ERROR: Error marking notification as read:', err);
    }
  };

  useEffect(() => {
    const fetchNotifications = async () => {
      console.log('DEBUG: Fetching notifications...')
      try {
        const res = await fetch(`${API_BASE_URL}/api/notifications`, {
          credentials: 'include'
        })
        console.log('DEBUG: Notifications fetch response status:', res.status)
        
        if (!res.ok) {
          console.error('ERROR: Notifications fetch failed with status:', res.status);
          throw new Error('Failed to fetch notifications')
        }
        const data = await res.json()
        console.log('DEBUG: Notifications data received:', data)
        console.log('DEBUG: Number of notifications:', Array.isArray(data) ? data.length : 'Not an array')
        
        setNotifications(Array.isArray(data) ? data : [])
      } catch (err) {
        console.error('ERROR: Error fetching notifications:', err)
        setNotifications([]) // Ensure notifications is always an array
      } finally {
        setLoading(false)
        console.log('DEBUG: Notifications fetch completed')
      }
    }

    fetchNotifications()
  }, [])

  const hasUnread = notifications && notifications.length > 0

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
        {loading ? (
          <DropdownMenuItem disabled>Loading...</DropdownMenuItem>
        ) : notifications && notifications.length > 0 ? (
          notifications.map((notif) => {
            const Icon = icons[notif.type] || Bell
            return (
              <DropdownMenuItem key={notif.id} asChild onClick={() => markAsRead(notif.id)}>
                <Link href={notif.href || '#'} className="flex items-start gap-3 w-full">
                  <Avatar className="h-8 w-8 mt-1">
                    <AvatarImage src={notif.user.avatar} />
                    <AvatarFallback>{notif.user.name[0]}</AvatarFallback>
                  </Avatar>
                  <div className="text-sm">
                    <p>
                      <span className="font-semibold">{notif.user.name}</span> {notif.message}
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

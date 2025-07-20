'use client'

import { useState, useEffect, useRef } from "react"
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
  const ws = useRef<WebSocket | null>(null)

  const markAsRead = async (id: string) => {
    console.log('DEBUG: Marking notification as read:', id)
    try {
      const res = await fetch(`${API_BASE_URL}/api/notifications/mark-read`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          notification_id: id
        }),
        credentials: 'include',
      });
      console.log('DEBUG: Mark as read response status:', res.status)
      
      if (!res.ok) {
        const errorText = await res.text();
        console.error('ERROR: Mark as read failed:', errorText);
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
      
      // Handle both array and null responses
      let notificationsArray = []
      if (Array.isArray(data)) {
        notificationsArray = data
      } else if (data && Array.isArray(data.notifications)) {
        notificationsArray = data.notifications
      } else if (data === null) {
        notificationsArray = []
      }
      
      console.log('DEBUG: Number of notifications:', notificationsArray.length)
      setNotifications(notificationsArray)
    } catch (err) {
      console.error('ERROR: Error fetching notifications:', err)
      setNotifications([])
    } finally {
      setLoading(false)
      console.log('DEBUG: Notifications fetch completed')
    }
  }

  // Setup WebSocket for real-time notifications
  useEffect(() => {
    console.log('Setting up WebSocket for notifications...')
    ws.current = new WebSocket(`ws://localhost:3000/ws`)
    
    ws.current.onopen = () => {
      console.log('Notification WebSocket connected')
    }

    ws.current.onmessage = (event) => {
      try {
        const messageData = JSON.parse(event.data)
        console.log('DEBUG: WebSocket notification received:', messageData)
        
        if (messageData.type === 'notification') {
          console.log('DEBUG: Processing real-time notification')
          // Refresh notifications when we receive a new one
          fetchNotifications()
        }
      } catch (error) {
        console.error('ERROR: Failed to parse WebSocket notification:', error)
      }
    }

    ws.current.onerror = (error) => {
      console.error('Notification WebSocket error:', error)
    }

    ws.current.onclose = () => {
      console.log('Notification WebSocket disconnected')
    }

    return () => {
      if (ws.current) {
        ws.current.close()
      }
    }
  }, [])

  useEffect(() => {
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

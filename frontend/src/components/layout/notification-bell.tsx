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
import { Bell, UserPlus, Users, Calendar, MessageSquare, Heart, UserCheck } from "lucide-react"
import Link from "next/link"
import { API_BASE_URL } from '@/lib/config';
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar"
import { useWebSocket } from "@/contexts/websocket-context"

type Notification = {
  id: string
  type: string
  message: string
  link: string
  is_read: boolean
  created_at: number
  reference_id: string
  actor_nickname?: string
  actor_avatar?: string
}

const icons: { [key: string]: React.ElementType } = {
  follow_request: UserPlus,
  new_follower: UserCheck,
  group_invite: Users,
  group_join_request: Users,
  event_created: Calendar,
  new_message: MessageSquare,
  new_like: Heart,
  new_comment: MessageSquare,
  group_join_response: Users,
  group_invitation_response: Users
}

export function NotificationBell() {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading] = useState(true)
  const { ws } = useWebSocket()

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
      
      if (!res.ok) {
        const errorText = await res.text();
        console.error('ERROR: Mark as read failed:', errorText);
        throw new Error('Failed to mark notification as read');
      }
      
      console.log('SUCCESS: Notification marked as read:', id)
      // Remove the notification from the list immediately
      setNotifications((prevNotifications) =>
        (prevNotifications || []).filter((notif) => notif.id !== id)
      );
    } catch (err) {
      console.error('ERROR: Error marking notification as read:', err)
    }
  }

  const handleNotificationClick = (notif: Notification) => {
    // Mark as read when clicked
    markAsRead(notif.id)
  }

  const fetchNotifications = async () => {
    console.log('DEBUG: Fetching notifications...')
    try {
      const res = await fetch(`${API_BASE_URL}/api/notifications`, {
        credentials: 'include'
      })
      
      if (!res.ok) {
        console.error('ERROR: Notifications fetch failed with status:', res.status);
        throw new Error('Failed to fetch notifications')
      }
      const data = await res.json()
      console.log('DEBUG: Notifications data received:', data)
      
      // Handle the array response directly
      const notificationsArray = Array.isArray(data) ? data : []
      
      console.log('DEBUG: Number of notifications:', notificationsArray.length)
      setNotifications(notificationsArray)
    } catch (err) {
      console.error('ERROR: Error fetching notifications:', err)
      setNotifications([])
    } finally {
      setLoading(false)
    }
  }

  // Listen for WebSocket notifications
  useEffect(() => {
    if (!ws) return

    const handleMessage = (event: MessageEvent) => {
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

    ws.addEventListener('message', handleMessage)

    return () => {
      ws.removeEventListener('message', handleMessage)
    }
  }, [ws])

  useEffect(() => {
    fetchNotifications()
  }, [])

  const hasUnread = notifications && notifications.length > 0

  const formatTime = (timestamp: number) => {
    const date = new Date(timestamp * 1000) // Convert Unix timestamp to milliseconds
    const now = new Date()
    const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60))
    
    if (diffInMinutes < 1) return 'Just now'
    if (diffInMinutes < 60) return `${diffInMinutes}m ago`
    if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}h ago`
    return `${Math.floor(diffInMinutes / 1440)}d ago`
  }

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
              <DropdownMenuItem key={notif.id} asChild>
                <Link 
                  href={notif.link || '#'} 
                  className="flex items-start gap-3 w-full"
                  onClick={() => handleNotificationClick(notif)}
                >
                  <Avatar className="h-8 w-8 mt-1">
                    <AvatarImage src={notif.actor_avatar} />
                    <AvatarFallback>
                      {notif.actor_nickname ? notif.actor_nickname[0] : <Icon className="h-4 w-4" />}
                    </AvatarFallback>
                  </Avatar>
                  <div className="text-sm">
                    <p>
                      <span className="font-semibold">{notif.actor_nickname || 'Someone'}</span> {notif.message}
                    </p>
                    <p className="text-xs text-muted-foreground">{formatTime(notif.created_at)}</p>
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

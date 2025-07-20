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

  console.log('DEBUG: NotificationBell - WebSocket connection status:', ws ? 'Connected' : 'Not connected')
  console.log('DEBUG: NotificationBell - WebSocket readyState:', ws?.readyState)

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

  const handleNotificationClick = async (notif: Notification, e: React.MouseEvent) => {
    e.preventDefault()
    
    // Mark as read first
    await markAsRead(notif.id)
    
    // Then navigate to the link
    if (notif.link && notif.link !== '#') {
      window.location.href = notif.link
    }
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
    if (!ws) {
      console.log('DEBUG: NotificationBell - No WebSocket connection available')
      return
    }
    console.log('DEBUG: NotificationBell - Setting up WebSocket listener')

    const handleMessage = (event: MessageEvent) => {
      console.log('DEBUG: NotificationBell - Raw WebSocket message received:', event.data)
      try {
        const messageData = JSON.parse(event.data)
        console.log('DEBUG: NotificationBell - Parsed WebSocket message:', messageData)
        
        if (messageData.type === 'notification') {
          console.log('DEBUG: NotificationBell - Processing notification type message')
          console.log('DEBUG: NotificationBell - Full messageData:', JSON.stringify(messageData, null, 2))
          
          // The notification data is in messageData.data
          const notificationPayload = messageData.data
          console.log('DEBUG: NotificationBell - Notification payload:', notificationPayload)
          
          if (notificationPayload && notificationPayload.notification) {
            const newNotification = {
              id: notificationPayload.notification.id,
              type: notificationPayload.notification.type,
              message: formatNotificationMessage(notificationPayload.notification.type, notificationPayload.data),
              link: formatNotificationLink(notificationPayload.notification.type, notificationPayload.notification.reference_id),
              is_read: false,
              created_at: Math.floor(new Date(notificationPayload.notification.created_at).getTime() / 1000),
              reference_id: notificationPayload.notification.reference_id,
              actor_nickname: notificationPayload.data?.actor_nickname,
              actor_avatar: notificationPayload.data?.actor_avatar
            }
            
            console.log('DEBUG: NotificationBell - Adding new notification to state:', newNotification)
            setNotifications(prev => {
              console.log('DEBUG: NotificationBell - Previous notifications:', prev)
              const updated = [newNotification, ...(prev || [])]
              console.log('DEBUG: NotificationBell - Updated notifications:', updated)
              return updated
            })
          } else {
            console.log('DEBUG: NotificationBell - Invalid notification payload structure:', notificationPayload)
          }
        } else {
          console.log('DEBUG: NotificationBell - Non-notification message type:', messageData.type)
        }
      } catch (error) {
        console.error('ERROR: NotificationBell - Failed to parse WebSocket message:', error)
      }
    }

    ws.addEventListener('message', handleMessage)
    console.log('DEBUG: NotificationBell - WebSocket listener added')

    return () => {
      console.log('DEBUG: NotificationBell - Removing WebSocket listener')
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

  // Helper functions to format notifications
  const formatNotificationMessage = (type: string, data: any) => {
    const actorName = data?.actor_nickname || 'Someone'
    switch (type) {
      case 'new_follower':
        return 'started following you'
      case 'follow_request':
        return 'sent you a follow request'
      case 'new_like':
        return 'liked your post'
      case 'new_comment':
        return 'commented on your post'
      case 'new_message':
        return 'sent you a message'
      case 'group_invite':
        return 'invited you to join a group'
      case 'group_join_request':
        return 'requested to join your group'
      case 'event_created':
        return 'created a new event in your group'
      default:
        return 'sent you a notification'
    }
  }

  const formatNotificationLink = (type: string, referenceId: string) => {
    switch (type) {
      case 'new_follower':
      case 'follow_request':
        return `/profile/${referenceId}`
      case 'new_like':
      case 'new_comment':
        return `/post/${referenceId}`
      case 'new_message':
        return `/chat/${referenceId}`
      case 'group_invite':
      case 'group_join_request':
        return `/groups/${referenceId}`
      case 'event_created':
        return `/events/${referenceId}`
      default:
        return '#'
    }
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
                <div 
                  className="flex items-start gap-3 w-full cursor-pointer p-2 hover:bg-accent"
                  onClick={(e) => handleNotificationClick(notif, e)}
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
                </div>
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

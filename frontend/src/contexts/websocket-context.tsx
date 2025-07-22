'use client'

import { createContext, useContext, useEffect, useState, ReactNode, useCallback, useRef } from 'react'
import { useUser } from './user-context'

interface WebSocketContextType {
  ws: WebSocket | null
  sendMessage: (message: any) => void
  isConnected: boolean
}

const WebSocketContext = createContext<WebSocketContextType>({
  ws: null,
  sendMessage: () => {},
  isConnected: false
})

export function WebSocketProvider({ children }: { children: ReactNode }) {
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const { user } = useUser()
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>()
  const reconnectAttemptsRef = useRef(0)
  const maxReconnectAttempts = 5

  const connectWebSocket = useCallback(() => {
    if (!user?.id) {
      console.log('DEBUG: No user ID, cannot connect WebSocket')
      return null
    }

    // Clear any existing reconnection timeout
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }

    console.log('Setting up global WebSocket connection for user:', user.id)
    
    const websocket = new WebSocket(`ws://localhost:3000/ws?user_id=${user.id}`)
    
    websocket.onopen = () => {
      console.log('Global WebSocket connected for user:', user.id)
      setIsConnected(true)
      reconnectAttemptsRef.current = 0 // Reset reconnect attempts on successful connection
    }
    
    websocket.onclose = (event) => {
      console.log('Global WebSocket disconnected:', event.code, event.reason)
      setIsConnected(false)
      
      // Only attempt to reconnect if we have a user and haven't exceeded max attempts
      if (user?.id && reconnectAttemptsRef.current < maxReconnectAttempts) {
        reconnectAttemptsRef.current++
        const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000) // Exponential backoff, max 30s
        
        console.log(`Attempting to reconnect WebSocket in ${delay}ms... (attempt ${reconnectAttemptsRef.current}/${maxReconnectAttempts})`)
        
        reconnectTimeoutRef.current = setTimeout(() => {
          if (user?.id) {
            const newWs = connectWebSocket()
            setWs(newWs)
          }
        }, delay)
      } else if (reconnectAttemptsRef.current >= maxReconnectAttempts) {
        console.log('Max reconnection attempts reached, giving up')
      }
    }
    
    websocket.onerror = (error) => {
      console.error('Global WebSocket error:', error)
      setIsConnected(false)
    }
    
    websocket.onmessage = (event) => {
      console.log('Global WebSocket message received:', event.data)
    }
    
    return websocket
  }, [user?.id])

  // Connect WebSocket when user changes
  useEffect(() => {
    if (user?.id && !ws) {
      const websocket = connectWebSocket()
      setWs(websocket)
    } else if (!user?.id && ws) {
      // Clean up WebSocket when user logs out
      console.log('User logged out, closing WebSocket')
      ws.close(1000, 'User logged out')
      setWs(null)
      setIsConnected(false)
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
    }

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
    }
  }, [user?.id, connectWebSocket])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (ws) {
        ws.close(1000, 'Component unmounting')
      }
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
    }
  }, [])

  const sendMessage = useCallback(
    (message: any) => {
      if (ws && isConnected) {
        ws.send(JSON.stringify(message))
      } else {
        console.warn('WebSocket is not connected. Cannot send message.')
      }
    },
    [ws, isConnected]
  )

  const value = {
    ws,
    sendMessage,
    isConnected,
  }

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>
}

export const useWebSocket = () => useContext(WebSocketContext)

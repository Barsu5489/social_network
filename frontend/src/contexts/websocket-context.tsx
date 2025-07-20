'use client'

import { createContext, useContext, useEffect, useState, ReactNode, useCallback } from 'react'
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

  const connectWebSocket = useCallback(() => {
    if (!user?.id) return

    console.log('Setting up global WebSocket connection for user:', user.id)
    const websocket = new WebSocket(`ws://localhost:3000/ws`)
    
    websocket.onopen = () => {
      console.log('Global WebSocket connected for user:', user.id)
      setWs(websocket)
      setIsConnected(true)
    }

    websocket.onerror = (error) => {
      console.error('Global WebSocket error:', error)
      console.error('WebSocket readyState:', websocket.readyState)
      setIsConnected(false)
    }

    websocket.onclose = (event) => {
      console.log('Global WebSocket disconnected. Code:', event.code, 'Reason:', event.reason)
      setWs(null)
      setIsConnected(false)
      
      // Only reconnect if user is still logged in and it wasn't a normal close
      if (user?.id && event.code !== 1000) {
        setTimeout(() => {
          console.log('Attempting to reconnect WebSocket...')
          connectWebSocket()
        }, 3000)
      }
    }

    websocket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        console.log('DEBUG: Global WebSocket message:', data)
      } catch (error) {
        console.error('ERROR: Failed to parse WebSocket message:', error)
      }
    }

    return websocket
  }, [user?.id])

  useEffect(() => {
    if (!user?.id) {
      setWs(null)
      setIsConnected(false)
      return
    }

    const websocket = connectWebSocket()

    return () => {
      if (websocket) {
        websocket.close()
      }
    }
  }, [user?.id, connectWebSocket])

  const sendMessage = (message: any) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message))
    }
  }

  return (
    <WebSocketContext.Provider value={{ ws, sendMessage, isConnected }}>
      {children}
    </WebSocketContext.Provider>
  )
}

export const useWebSocket = () => useContext(WebSocketContext)

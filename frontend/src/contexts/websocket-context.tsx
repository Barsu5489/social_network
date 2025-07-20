'use client'

import { createContext, useContext, useEffect, useState, ReactNode } from 'react'
import { useUser } from './user-context'

interface WebSocketContextType {
  ws: WebSocket | null
  sendMessage: (message: any) => void
}

const WebSocketContext = createContext<WebSocketContextType>({
  ws: null,
  sendMessage: () => {}
})

export function WebSocketProvider({ children }: { children: ReactNode }) {
  const [ws, setWs] = useState<WebSocket | null>(null)
  const { user } = useUser()

  useEffect(() => {
    if (!user?.id) return

    console.log('Setting up global WebSocket connection...')
    const websocket = new WebSocket(`ws://localhost:3000/ws`)
    
    websocket.onopen = () => {
      console.log('Global WebSocket connected')
      setWs(websocket)
    }

    websocket.onerror = (error) => {
      console.error('Global WebSocket error:', error)
    }

    websocket.onclose = () => {
      console.log('Global WebSocket disconnected')
      setWs(null)
    }

    websocket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        console.log('DEBUG: Global WebSocket message:', data)
        
        // Let components handle their own message types
        // This is just for logging and debugging
      } catch (error) {
        console.error('ERROR: Failed to parse WebSocket message:', error)
      }
    }

    return () => {
      if (websocket) {
        websocket.close()
      }
    }
  }, [user?.id])

  const sendMessage = (message: any) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message))
    }
  }

  return (
    <WebSocketContext.Provider value={{ ws, sendMessage }}>
      {children}
    </WebSocketContext.Provider>
  )
}

export const useWebSocket = () => useContext(WebSocketContext)

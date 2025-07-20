'use client'

import { createContext, useContext, useEffect, useRef, ReactNode } from 'react'
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
  const ws = useRef<WebSocket | null>(null)
  const { user } = useUser()

  useEffect(() => {
    if (!user?.id) return

    console.log('Setting up global WebSocket connection...')
    ws.current = new WebSocket(`ws://localhost:3000/ws`)
    
    ws.current.onopen = () => {
      console.log('Global WebSocket connected')
    }

    ws.current.onerror = (error) => {
      console.error('Global WebSocket error:', error)
    }

    ws.current.onclose = () => {
      console.log('Global WebSocket disconnected')
    }

    ws.current.onmessage = (event) => {
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
      if (ws.current) {
        ws.current.close()
      }
    }
  }, [user?.id])

  const sendMessage = (message: any) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(message))
    }
  }

  return (
    <WebSocketContext.Provider value={{ ws: ws.current, sendMessage }}>
      {children}
    </WebSocketContext.Provider>
  )
}

export const useWebSocket = () => useContext(WebSocketContext)

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

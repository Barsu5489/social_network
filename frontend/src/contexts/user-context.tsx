
'use client';
import React, { createContext, useContext, useState, useEffect, ReactNode, useCallback, useMemo } from 'react';
import { API_BASE_URL } from '@/lib/config';

// This interface matches the 'data' object in the successful login response
// and the shape of the user object in the profile response.
interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  nickname?: string;
  avatar_url?: string;
  is_private?: boolean;
}

interface UserContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  isLoading: boolean;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUserState] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [hasInitialized, setHasInitialized] = useState(false)

  const setUser = useCallback((userData: User | null) => {
    setUserState(userData)
    if (userData) {
      localStorage.setItem('user', JSON.stringify(userData))
    } else {
      localStorage.removeItem('user')
    }
  }, [])

  useEffect(() => {
    const initializeUser = async () => {
      if (hasInitialized) return
      
      setIsLoading(true)
      
      // First check localStorage
      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        try {
          const userData = JSON.parse(storedUser)
          setUserState(userData)
        } catch (error) {
          console.error('Error parsing stored user:', error)
          localStorage.removeItem('user')
        }
      }

      // Then verify with server
      try {
        const response = await fetch(`${API_BASE_URL}/api/profile`, {
          credentials: 'include',
        })
        
        if (response.ok) {
          const userData = await response.json()
          console.log('User verified with server:', userData.id)
          setUser(userData)
        } else {
          console.log('No valid session, clearing user')
          setUser(null)
        }
      } catch (error) {
        console.error('Error fetching user:', error)
        setUser(null)
      } finally {
        setIsLoading(false)
        setHasInitialized(true)
      }
    }

    initializeUser()
  }, []) // Empty dependency array - only run once

  return (
    <UserContext.Provider value={{ user, setUser, isLoading }}>
      {children}
    </UserContext.Provider>
  )
}

export const useUser = () => {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};

    

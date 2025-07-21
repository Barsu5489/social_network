
'use client';
import React, { createContext, useContext, useState, useEffect, ReactNode, useCallback, useMemo } from 'react';
import { API_BASE_URL } from '@/lib/config';
import { useDebounce } from '@/hooks/use-debounce';

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
  isAuthenticated: boolean;
  isVerifying: boolean;
  verifyUser: () => Promise<void>;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const [isVerifying, setIsVerifying] = useState(false)

  const verifyUser = useCallback(async () => {
    if (isVerifying) {
      console.log('DEBUG: Already verifying user, skipping...')
      return
    }
    
    setIsVerifying(true)
    console.log('DEBUG: Starting user verification...')
    
    try {
      const response = await fetch(`${API_BASE_URL}/api/profile`, {
        credentials: 'include',
      })
      
      if (response.ok) {
        const data = await response.json()
        if (data.success && data.data?.user) {
          console.log('User verified with server:', data.data.user.id)
          setUser(data.data.user)
          setIsAuthenticated(true)
        } else {
          console.log('No valid session, clearing user')
          setUser(null)
          setIsAuthenticated(false)
        }
      } else {
        console.log('No valid session, clearing user')
        setUser(null)
        setIsAuthenticated(false)
      }
    } catch (error) {
      console.error('Error verifying user:', error)
      setUser(null)
      setIsAuthenticated(false)
    } finally {
      setIsVerifying(false)
      setIsLoading(false)
    }
  }, [isVerifying])

  const debouncedVerifyUser = useDebounce(verifyUser, 500)

  useEffect(() => {
    if (!isLoading || isVerifying) return
    debouncedVerifyUser()
  }, [debouncedVerifyUser, isLoading, isVerifying])

  const value = {
    user,
    setUser,
    isAuthenticated,
    setIsAuthenticated,
    isLoading,
    isVerifying,
    verifyUser,
  }

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>
}

export const useUser = () => {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};

    


'use client';
import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';

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

export const UserProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUserState] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // This effect runs only on the client-side
    try {
      const storedUser = localStorage.getItem('connectu-user');
      if (storedUser) {
        setUserState(JSON.parse(storedUser));
      }
    } catch (error) {
      console.error("Failed to parse user from localStorage", error);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const setUser = (user: User | null) => {
    setUserState(user);
    if (user) {
      localStorage.setItem('connectu-user', JSON.stringify(user));
    } else {
      localStorage.removeItem('connectu-user');
    }
  };

  return (
    <UserContext.Provider value={{ user, setUser, isLoading }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};

    
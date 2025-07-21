'use client';

import { useState, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { ScrollArea } from "@/components/ui/scroll-area";
import { API_BASE_URL } from '@/lib/config';

interface Follower {
  id: string;
  first_name: string;
  last_name: string;
  nickname: string;
  avatar_url?: string;
}

interface FollowerSelectorProps {
  selectedUsers: string[];
  onSelectionChange: (users: string[]) => void;
}

export function FollowerSelector({ selectedUsers, onSelectionChange }: FollowerSelectorProps) {
  const [followers, setFollowers] = useState<Follower[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchFollowers = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/api/followers`, {
          credentials: 'include',
        });
        if (response.ok) {
          const data = await response.json();
          setFollowers(data.followers || []);
        }
      } catch (error) {
        console.error('Error fetching followers:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchFollowers();
  }, []);

  const handleUserToggle = (userId: string) => {
    const newSelection = selectedUsers.includes(userId)
      ? selectedUsers.filter(id => id !== userId)
      : [...selectedUsers, userId];
    onSelectionChange(newSelection);
  };

  if (loading) {
    return <div className="p-4">Loading followers...</div>;
  }

  return (
    <div className="border rounded-lg p-4 mt-2">
      <h4 className="font-medium mb-2">Select followers who can see this post:</h4>
      <ScrollArea className="h-32">
        {followers.map((follower) => (
          <div key={follower.id} className="flex items-center space-x-2 py-1">
            <Checkbox
              id={follower.id}
              checked={selectedUsers.includes(follower.id)}
              onCheckedChange={() => handleUserToggle(follower.id)}
            />
            <label htmlFor={follower.id} className="text-sm cursor-pointer">
              {follower.first_name} {follower.last_name} (@{follower.nickname})
            </label>
          </div>
        ))}
      </ScrollArea>
      {followers.length === 0 && (
        <p className="text-sm text-muted-foreground">No followers found.</p>
      )}
    </div>
  );
}
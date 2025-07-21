'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Check, X } from 'lucide-react';
import { API_BASE_URL } from '@/lib/config';
import { useToast } from '@/hooks/use-toast';

interface FollowRequestItemProps {
  notification: {
    id: string;
    reference_id: string; // follower_id
    actor_nickname?: string;
    actor_avatar?: string;
  };
  onResponse: (notificationId: string, accepted: boolean) => void;
}

export function FollowRequestItem({ notification, onResponse }: FollowRequestItemProps) {
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  const handleResponse = async (accept: boolean) => {
    setIsLoading(true);
    try {
      const endpoint = accept ? 'accept' : 'decline';
      const response = await fetch(`${API_BASE_URL}/api/follow-requests/${notification.reference_id}/${endpoint}`, {
        method: 'POST',
        credentials: 'include',
      });

      if (response.ok) {
        onResponse(notification.id, accept);
        toast({
          title: accept ? 'Follow request accepted' : 'Follow request declined',
        });
      } else {
        throw new Error('Failed to respond to follow request');
      }
    } catch (error) {
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'Failed to respond to follow request',
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="bg-card/50">
      <CardContent className="p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Avatar>
              <AvatarImage src={notification.actor_avatar || `https://i.pravatar.cc/40?u=${notification.reference_id}`} />
              <AvatarFallback>{notification.actor_nickname?.[0] || 'U'}</AvatarFallback>
            </Avatar>
            <div>
              <p className="font-medium">@{notification.actor_nickname || 'Unknown'}</p>
              <p className="text-sm text-muted-foreground">wants to follow you</p>
            </div>
          </div>
          <div className="flex gap-2">
            <Button
              size="sm"
              onClick={() => handleResponse(true)}
              disabled={isLoading}
            >
              <Check className="h-4 w-4 mr-1" />
              Accept
            </Button>
            <Button
              size="sm"
              variant="outline"
              onClick={() => handleResponse(false)}
              disabled={isLoading}
            >
              <X className="h-4 w-4 mr-1" />
              Decline
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
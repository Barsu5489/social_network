import React from 'react';
import { useRouter } from 'next/navigation';
import { API_BASE_URL } from '@/lib/config';
import { useToast } from '@/hooks/use-toast';
import { Button } from '@/components/ui/button';
import { MessageCircle } from 'lucide-react';

interface StartChatButtonProps {
  userId: string;
}

const StartChatButton: React.FC<StartChatButtonProps> = ({ userId }) => {
  const router = useRouter();
  const { toast } = useToast();

  const handleStartChat = async () => {
    try {
      // Check follow/following relationship
      const followCheckResponse = await fetch(
        `${API_BASE_URL}/api/follow/check?targetUserId=${userId}`,
        {
          method: 'GET',
          credentials: 'include',
        }
      );

      if (!followCheckResponse.ok) {
        throw new Error('Failed to check follow relationship');
      }

      const followCheckData = await followCheckResponse.json();

      if (!followCheckData.isFollowing && !followCheckData.isFollowedBy) {
        console.log('User does not follow or is followed by this user.');
        toast({
          variant: 'destructive',
          title: 'Error',
          description: 'You can only chat with users you follow or who follow you.',
        });
        return;
      }

      // Create direct chat
      const response = await fetch(`${API_BASE_URL}/api/chats/direct`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ recipient_id: userId }),
        credentials: 'include',
      });

      const data = await response.json();
      console.log('Create chat response:', data);

      if (response.ok && data.chat_id) {
        console.log('Navigating to chat:', data.chat_id);
        router.push(`/chat/${data.chat_id}`);
      } else {
        console.error('Failed to create chat. Status:', response.status, 'Data:', data);
        toast({
          variant: 'destructive',
          title: 'Error',
          description: data.error || 'Failed to create chat.',
        });
      }
    } catch (error: any) {
      console.error('Error in handleStartChat:', error);
      toast({
        variant: 'destructive',
        title: 'Error',
        description: error.message || 'Failed to create chat.',
      });
    }
  };

  return (
      <Button onClick={handleStartChat} variant="outline" size="sm">
          <MessageCircle className="mr-2 h-4 w-4" />
          Message
      </Button>
  );
};

export default StartChatButton;

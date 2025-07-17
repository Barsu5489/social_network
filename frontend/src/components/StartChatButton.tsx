import React from 'react';
import { useRouter } from 'next/navigation';
import { API_BASE_URL } from '@/lib/config';
import { useToast } from '@/hooks/use-toast';

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

      const followCheckData = await followCheckResponse.json();

      if (!followCheckData.isFollowing && !followCheckData.isFollowedBy) {
        toast({
          variant: 'destructive',
          title: 'Error',
          description:
            'You can only chat with users you follow or who follow you.',
        });
        return;
      }

      // Create direct chat
      const response = await fetch(`${API_BASE_URL}/api/chats/direct`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ targetUserId: userId }),
        credentials: 'include',
      });

      const data = await response.json();

      if (response.ok && data.success) {
        router.push(`/chat/${data.chatId}`);
      } else {
        toast({
          variant: 'destructive',
          title: 'Error',
          description: 'Failed to create chat.',
        });
      }
    } catch (error) {
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'Failed to create chat.',
      });
    }
  };

  return (
    <button onClick={handleStartChat}>
      Start Chat
    </button>
  );
};

export default StartChatButton;

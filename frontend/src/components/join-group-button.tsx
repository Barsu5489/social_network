
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from './ui/button';
import { useToast } from '@/hooks/use-toast';
import { API_BASE_URL } from '@/lib/config';
import { Plus, Check } from 'lucide-react';
import { useUser } from '@/contexts/user-context';

interface JoinGroupButtonProps {
    groupId: string;
    isPrivate: boolean;
}

export function JoinGroupButton({ groupId, isPrivate }: JoinGroupButtonProps) {
    const [isLoading, setIsLoading] = useState(false);
    const [isRequested, setIsRequested] = useState(false);
    const { toast } = useToast();
    const router = useRouter();
    const { user } = useUser();

    const handleClick = async () => {
        if (!user) {
            toast({ variant: 'destructive', title: 'Not signed in', description: 'You must be signed in to join a group.' });
            return;
        }

        setIsLoading(true);
        try {
            const response = await fetch(`${API_BASE_URL}/api/groups/join/${groupId}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ user_id: user.id }),
                credentials: 'include',
            });
            
            const data = await response.json();

            if (response.ok) {
                toast({
                    title: 'Request Sent',
                    description: `Your request to join the group has been sent.`,
                });
                setIsRequested(true);
            } else {
                 throw new Error(data.error || 'Failed to send request.');
            }

        } catch (error: any) {
            toast({
                variant: 'destructive',
                title: 'Error',
                description: error.message,
            });
        } finally {
            setIsLoading(false);
        }
    };
    
    if (isRequested) {
        return (
             <Button variant="outline" disabled>
                <Check className="h-4 w-4 mr-1" />
                Requested
            </Button>
        )
    }

    return (
        <Button variant="outline" onClick={handleClick} disabled={isLoading}>
            <Plus className="h-4 w-4 mr-1" />
            {isLoading ? 'Sending...' : isPrivate ? 'Request to Join' : 'Join Group'}
        </Button>
    );
}

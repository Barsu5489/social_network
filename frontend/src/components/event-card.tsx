
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from './ui/card';
import { Button } from './ui/button';
import { Calendar, Check, ThumbsDown, Users } from 'lucide-react';
import type { Event } from '@/types';
import { API_BASE_URL } from '@/lib/config';
import { useToast } from '@/hooks/use-toast';
import { useUser } from '@/contexts/user-context';
import { cn } from '@/lib/utils';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip';

interface EventCardProps {
    event: Event;
}

export function EventCard({ event }: EventCardProps) {
    const { user } = useUser();
    const router = useRouter();
    const { toast } = useToast();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [attendees, setAttendees] = useState(event.attendees);

    const formatTimestamp = (timestamp: number) => {
        return new Date(timestamp * 1000).toLocaleString(undefined, {
            dateStyle: 'medium',
            timeStyle: 'short',
        });
    };

    const handleRsvp = async (status: 'going' | 'not_going') => {
        setIsSubmitting(true);
        try {
            const response = await fetch(`${API_BASE_URL}/api/events/${event.id}/rsvp`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ status }),
                credentials: 'include',
            });
            if (response.ok) {
                toast({ title: 'RSVP submitted!' });
                router.refresh(); 
            } else {
                throw new Error('Failed to RSVP');
            }
        } catch (error: any) {
            toast({ variant: 'destructive', title: 'Error', description: error.message });
        } finally {
            setIsSubmitting(false);
        }
    };
    
    const userRsvpStatus = attendees.find(a => a.user_id === user?.id)?.status;
    const goingAttendees = attendees.filter(a => a.status === 'going');

    return (
        <Card>
            <CardHeader>
                <CardTitle>{event.title}</CardTitle>
                <CardDescription>{event.location}</CardDescription>
            </CardHeader>
            <CardContent>
                <p className="text-sm mb-4">{event.description}</p>
                <div className="flex items-center text-sm text-muted-foreground">
                    <Calendar className="mr-2 h-4 w-4" />
                    <span>{formatTimestamp(event.start_time)} - {formatTimestamp(event.end_time)}</span>
                </div>
                 <div className="flex items-center mt-4">
                    <Users className="mr-2 h-4 w-4 text-muted-foreground" />
                    <div className="flex -space-x-2">
                         <TooltipProvider>
                        {goingAttendees.slice(0, 5).map(attendee => (
                           <Tooltip key={attendee.user_id}>
                               <TooltipTrigger asChild>
                                    <Avatar className="h-6 w-6 border-2 border-background">
                                        <AvatarImage src={`https://i.pravatar.cc/40?u=${attendee.user_id}`} data-ai-hint="person" />
                                        <AvatarFallback>{attendee.user_name[0]}</AvatarFallback>
                                    </Avatar>
                               </TooltipTrigger>
                               <TooltipContent>
                                   <p>{attendee.user_name}</p>
                               </TooltipContent>
                           </Tooltip>
                        ))}
                         </TooltipProvider>
                    </div>
                    {goingAttendees.length > 0 ? (
                        <span className="text-xs text-muted-foreground ml-2">
                           {goingAttendees.length} going
                        </span>
                    ) : (
                         <span className="text-xs text-muted-foreground ml-2">
                           No one is going yet.
                        </span>
                    )}
                </div>
            </CardContent>
            <CardFooter className="flex justify-end gap-2">
                <Button
                    variant={userRsvpStatus === 'going' ? 'default' : 'outline'}
                    onClick={() => handleRsvp('going')}
                    disabled={isSubmitting}
                >
                    <Check className="mr-2 h-4 w-4" />
                    Going
                </Button>
                <Button
                    variant={userRsvpStatus === 'not_going' ? 'destructive' : 'outline'}
                    onClick={() => handleRsvp('not_going')}
                    disabled={isSubmitting}
                >
                    <ThumbsDown className="mr-2 h-4 w-4" />
                    Not Going
                </Button>
            </CardFooter>
        </Card>
    );
}

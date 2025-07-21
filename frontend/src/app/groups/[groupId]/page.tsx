
'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { API_BASE_URL } from "@/lib/config";
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { PostCard } from "@/components/post-card";
import { CreatePost } from "@/components/create-post";
import type { Post, Group, Event } from "@/types";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CreateEventDialog } from "@/components/create-event-dialog";
import { EventCard } from "@/components/event-card";
import Image from "next/image";
import { getGroupCover } from "@/lib/avatars";
import { useUser } from '@/contexts/user-context';
import { Skeleton } from '@/components/ui/skeleton';
import { useToast } from "@/hooks/use-toast";
import { Button } from "@/components/ui/button";
import { InviteUserDialog } from "@/components/invite-user-dialog";
import { JoinGroupButton } from "@/components/join-group-button";

type GroupRequest = {
    id: string;
    user_name: string;
    members?: { id: string }[]; 
};

function GroupRequestsTab({ groupId }: { groupId: string }) {
    const [requests, setRequests] = useState<GroupRequest[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const { toast } = useToast();

    useEffect(() => {
        const fetchRequests = async () => {
            try {
                const res = await fetch(`${API_BASE_URL}/api/groups/${groupId}/requests`, {
                    credentials: 'include'
                });
                if (res.ok) {
                    const data = await res.json();
                    setRequests(data);
                }
            } catch (error) {
                console.error('Error fetching requests:', error);
            } finally {
                setIsLoading(false);
            }
        };
        fetchRequests();
    }, [groupId]);

    const handleRequest = async (requestId: string, action: 'accept' | 'reject') => {
        try {
            const res = await fetch(`${API_BASE_URL}/api/groups/requests/${requestId}/${action}`, {
                method: 'POST',
                credentials: 'include'
            });
            if (res.ok) {
                setRequests(prev => prev.filter(req => req.id !== requestId));
                toast({ title: `Request ${action}ed successfully` });
            }
        } catch (error) {
            toast({ variant: 'destructive', title: `Failed to ${action} request` });
        }
    };

    if (isLoading) return <Skeleton className="h-32 w-full" />;

    return (
        <div className="space-y-4">
            {requests.length > 0 ? (
                requests.map((request) => (
                    <Card key={request.id}>
                        <CardContent className="flex items-center justify-between p-4">
                            <div>
                                <p className="font-medium">{request.user_name}</p>
                                <p className="text-sm text-muted-foreground">Wants to join the group</p>
                            </div>
                            <div className="flex gap-2">
                                <Button size="sm" onClick={() => handleRequest(request.id, 'accept')}>
                                    Accept
                                </Button>
                                <Button size="sm" variant="outline" onClick={() => handleRequest(request.id, 'reject')}>
                                    Reject
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                ))
            ) : (
                <Card>
                    <CardContent className="p-8 text-center text-muted-foreground">
                        <p>No pending requests.</p>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}

export default function SingleGroupPage() {
    const params = useParams();
    const groupId = params.groupId as string;
    const { user } = useUser();

    const [groupData, setGroupData] = useState<{ posts: Post[], events: Event[], groupDetails: Group | null }>({ posts: [], events: [], groupDetails: null });
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isMember, setIsMember] = useState<boolean>(false);

    useEffect(() => {
        if (!groupId || !user) return;

        const getGroupData = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const fetchPosts = fetch(`${API_BASE_URL}/api/groups/${groupId}/posts?user_id=${user.id}`, { credentials: 'include' }).then(res => res.ok ? res.json() : []);
                
                const fetchEvents = fetch(`${API_BASE_URL}/api/groups/${groupId}/events`, { credentials: 'include' }).then(res => res.ok ? res.json() : []);
                
                const fetchGroupDetails = fetch(`${API_BASE_URL}/api/groups`, { credentials: 'include' }).then(async res => {
                    if (!res.ok) return null;
                    const allGroups = await res.json();
                    return Array.isArray(allGroups) ? allGroups.find(g => g.id === groupId) : null;
                });
                
                const [posts, events, groupDetails] = await Promise.all([fetchPosts, fetchEvents, fetchGroupDetails]);

                if (!groupDetails) {
                    throw new Error("Group not found or you do not have permission to view it.");
                }

                setGroupData({ posts, events, groupDetails });
            } catch (err: any) {
                setError(err.message);
            } finally {
                setIsLoading(false);
            }
        };
        
        getGroupData();
    }, [groupId, user]);

    useEffect(() => {
        if (!groupId || !user) return;

        const checkMembership = async () => {
            try {
                const res = await fetch(`${API_BASE_URL}/api/groups/${groupId}/members`, {
                    credentials: 'include'
                });
                if (res.ok) {
                    const members = await res.json();
                    setIsMember(members.some((member: any) => member.user_id === user.id));
                }
            } catch (error) {
                console.error('Error checking membership:', error);
            }
        };

        checkMembership();
    }, [groupId, user]);

    if (isLoading) {
        return (
            <main className="flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
                 <Card className="mb-6 overflow-hidden">
                    <Skeleton className="h-48 w-full" />
                    <CardHeader>
                        <Skeleton className="h-8 w-1/2" />
                        <Skeleton className="h-5 w-3/4 mt-2" />
                    </CardHeader>
                </Card>
                <div className="flex justify-between items-center mb-4">
                    <Skeleton className="h-10 w-48" />
                    <Skeleton className="h-10 w-32" />
                </div>
                <Skeleton className="h-64 w-full" />
            </main>
        )
    }

    if (error || !groupData.groupDetails) {
        return (
            <main className="flex-1 flex items-center justify-center p-6">
                <Card>
                    <CardContent className="p-8 text-center text-muted-foreground">
                        <p>{error || "Group not found."}</p>
                    </CardContent>
                </Card>
            </main>
        )
    }

    const { posts, events, groupDetails } = groupData;
    const cover = getGroupCover(groupId);

    return (
        <main className="flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
            <Card className="mb-6 overflow-hidden">
                <div className="relative h-48 bg-muted">
                    <Image 
                        src={cover.url} 
                        alt={groupDetails.name} 
                        fill
                        className="object-cover"
                        data-ai-hint={cover.hint}
                    />
                </div>
                <CardHeader>
                    <CardTitle className="text-3xl">{groupDetails.name}</CardTitle>
                    <CardDescription>{groupDetails.description}</CardDescription>
                </CardHeader>
                <CardFooter className="!flex justify-between">
                    {isMember ? (
                        <InviteUserDialog groupId={groupDetails.id} />
                    ) : (
                        <JoinGroupButton groupId={groupDetails.id} isPrivate={groupDetails.is_private} />
                    )}
                </CardFooter>
            </Card>

            <Tabs defaultValue="posts">
                <div className="flex justify-between items-center mb-4">
                    <TabsList>
                        <TabsTrigger value="posts">Posts</TabsTrigger>
                        <TabsTrigger value="events">Events</TabsTrigger>
                        {isMember && <TabsTrigger value="requests">Requests</TabsTrigger>}
                    </TabsList>
                    <CreateEventDialog groupId={groupId} />
                </div>
                
                <TabsContent value="posts">
                    <CreatePost groupId={groupId} />
                    <div className="mt-6 space-y-6">
                        {posts && posts.length > 0 ? (
                            posts.map((post: Post) => (
                              <PostCard key={post.id} {...post} />
                            ))
                        ) : (
                            <Card>
                                <CardContent className="p-8 text-center text-muted-foreground">
                                    <p>No posts in this group yet.</p>
                                </CardContent>
                            </Card>
                        )}
                    </div>
                </TabsContent>

                <TabsContent value="events">
                     <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        {events && events.length > 0 ? (
                            events.map((event: Event) => (
                              <EventCard key={event.id} event={event} />
                            ))
                        ) : (
                            <Card className="md:col-span-2">
                                <CardContent className="p-8 text-center text-muted-foreground">
                                    <p>No upcoming events in this group.</p>
                                </CardContent>
                            </Card>
                        )}
                    </div>
                </TabsContent>

                {isMember && (
                    <TabsContent value="requests">
                       <GroupRequestsTab groupId={groupId} />
                    </TabsContent>
                )}
            </Tabs>
        </main>
    );
}

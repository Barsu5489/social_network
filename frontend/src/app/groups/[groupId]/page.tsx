
'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { API_BASE_URL } from "@/lib/config";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
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
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Terminal } from 'lucide-react';

// Placeholder for the requests tab content
function GroupRequestsTab({ groupId }: { groupId: string }) {
    // In a real application, you would fetch pending requests here.
    // const [requests, setRequests] = useState([]);
    // const [isLoading, setIsLoading] = useState(true);
    // useEffect(() => {
    //   fetch(`${API_BASE_URL}/api/groups/${groupId}/requests`) ...
    // }, [groupId]);

    return (
        <Alert>
            <Terminal className="h-4 w-4" />
            <AlertTitle>Under Construction!</AlertTitle>
            <AlertDescription>
                The ability to manage join requests is coming soon. A backend endpoint to fetch pending requests is needed to complete this feature.
            </AlertDescription>
        </Alert>
    );
}

export default function SingleGroupPage() {
    const params = useParams();
    const groupId = params.groupId as string;
    const { user } = useUser();

    const [groupData, setGroupData] = useState<{ posts: Post[], events: Event[], groupDetails: Group | null }>({ posts: [], events: [], groupDetails: null });
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!groupId || !user) return;

        const getGroupData = async () => {
            setIsLoading(true);
            setError(null);
            try {
                // Pass user_id query param to fetch posts
                const fetchPosts = fetch(`${API_BASE_URL}/api/groups/${groupId}/posts?user_id=${user.id}`, { credentials: 'include', cache: 'no-store' }).then(res => res.ok ? res.json() : []);
                
                const fetchEvents = fetch(`${API_BASE_URL}/api/groups/${groupId}/events`, { credentials: 'include', cache: 'no-store' }).then(res => res.ok ? res.json() : []);
                
                const fetchGroupDetails = fetch(`${API_BASE_URL}/api/groups`, { credentials: 'include', cache: 'no-store' }).then(async res => {
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
    const isCreator = user?.id === groupDetails.creator_id;

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
            </Card>

            <Tabs defaultValue="posts">
                <div className="flex justify-between items-center mb-4">
                    <TabsList>
                        <TabsTrigger value="posts">Posts</TabsTrigger>
                        <TabsTrigger value="events">Events</TabsTrigger>
                        {isCreator && <TabsTrigger value="requests">Requests</TabsTrigger>}
                    </TabsList>
                    <CreateEventDialog groupId={groupId} />
                </div>
                
                <TabsContent value="posts">
                    {/* Pass groupId to CreatePost component */}
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

                {isCreator && (
                    <TabsContent value="requests">
                       <GroupRequestsTab groupId={groupId} />
                    </TabsContent>
                )}
            </Tabs>
        </main>
    );
}

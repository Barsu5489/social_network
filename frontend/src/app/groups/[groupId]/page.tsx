
import { cookies } from "next/headers";
import { API_BASE_URL } from "@/lib/config";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { PostCard } from "@/components/post-card";
import { CreatePost } from "@/components/create-post";
import type { Post, Group, Event } from "@/types";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CreateEventDialog } from "@/components/create-event-dialog";
import { EventCard } from "@/components/event-card";
import Image from "next/image";


async function getGroupData(groupId: string, sessionCookie: any) {
    const headers = {
        'Cookie': sessionCookie ? `${sessionCookie.name}=${sessionCookie.value}` : '',
    };

    const fetchPosts = async () => {
        const res = await fetch(`${API_BASE_URL}/api/groups/${groupId}/posts`, { headers, cache: 'no-store' });
        if (!res.ok) return [];
        return res.json();
    };

    const fetchEvents = async () => {
        const res = await fetch(`${API_BASE_URL}/api/groups/${groupId}/events`, { headers, cache: 'no-store' });
        if (!res.ok) return [];
        return res.json();
    };

    const fetchGroupDetails = async () => {
        const res = await fetch(`${API_BASE_URL}/api/groups`, { headers, cache: 'no-store' });
        if (!res.ok) return null;
        const allGroups = await res.json();
        return Array.isArray(allGroups) ? allGroups.find(g => g.id === groupId) : null;
    }

    try {
        const [posts, events, groupDetails] = await Promise.all([
            fetchPosts(),
            fetchEvents(),
            fetchGroupDetails()
        ]);
        return { posts, events, groupDetails };
    } catch (error) {
        console.error('Error fetching group data:', error);
        return { posts: [], events: [], groupDetails: null };
    }
}


export default async function SingleGroupPage({ params }: { params: { groupId: string } }) {
    const sessionCookie = cookies().get('social-network-session');
    
    const { posts, events, groupDetails } = await getGroupData(params.groupId, sessionCookie);
    
    if (!groupDetails) {
        return (
            <main className="flex-1 flex items-center justify-center p-6">
                <Card>
                    <CardContent className="p-8 text-center text-muted-foreground">
                        <p>Group not found or you do not have permission to view it.</p>
                    </CardContent>
                </Card>
            </main>
        )
    }

    return (
        <main className="flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
            <Card className="mb-6 overflow-hidden">
                <div className="relative h-48 bg-muted">
                    <Image 
                        src={`https://placehold.co/1200x400.png`} 
                        alt={groupDetails.name} 
                        fill
                        className="object-cover"
                        data-ai-hint="community abstract"
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
                    </TabsList>
                    <CreateEventDialog groupId={params.groupId} />
                </div>
                
                <TabsContent value="posts">
                    <CreatePost />
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
            </Tabs>
        </main>
    );
}

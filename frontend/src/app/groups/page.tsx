
import { cookies } from "next/headers";
import { API_BASE_URL } from "@/lib/config";
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import Image from "next/image";
import Link from "next/link";
import { CreateGroupDialog } from "@/components/create-group-dialog";
import type { Group } from "@/types";
import { JoinGroupButton } from "@/components/join-group-button";
import { Button } from "@/components/ui/button";


export default async function GroupsPage() {
    let groups: Group[] = [];
    try {
        const cookieStore = cookies();
        const sessionCookie = cookieStore.get('social-network-session');

        const res = await fetch(`${API_BASE_URL}/api/groups`, {
            headers: {
                'Cookie': sessionCookie ? `${sessionCookie.name}=${sessionCookie.value}` : '',
            },
            cache: 'no-store' 
        });
        if (!res.ok) {
            console.error('Failed to fetch groups:', res.statusText);
        } else {
            const data = await res.json();
            groups = Array.isArray(data) ? data : [];
        }
    } catch (error) {
        console.error('Error fetching groups:', error);
    }

    return (
        <main className="flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Discover Groups</h1>
                <CreateGroupDialog />
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {groups.length > 0 ? (
                    groups.map((group) => (
                        <Card key={group.id} className="flex flex-col">
                            <CardHeader>
                                <Image 
                                    src={`https://placehold.co/600x400.png`} 
                                    alt={group.name} 
                                    width={600} 
                                    height={400} 
                                    className="rounded-t-lg aspect-video object-cover"
                                    data-ai-hint="community event"
                                />
                                <CardTitle className="pt-4">{group.name}</CardTitle>
                                <CardDescription className="line-clamp-2 h-[40px]">{group.description}</CardDescription>
                            </CardHeader>
                            <CardContent className="flex-grow">
                                {/* Could add member count here if available */}
                            </CardContent>
                            <CardFooter className="flex justify-between">
                                <Link href={`/groups/${group.id}`} passHref>
                                    <Button>View Group</Button>
                                </Link>
                                <JoinGroupButton groupId={group.id} isPrivate={group.is_private} />
                            </CardFooter>
                        </Card>
                    ))
                ) : (
                    <p className="col-span-full text-center text-muted-foreground py-8">No groups found. Why not create one?</p>
                )}
            </div>
        </main>
    );
}

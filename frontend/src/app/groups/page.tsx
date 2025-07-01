'use client';

import { useEffect, useState } from "react";
import { API_BASE_URL } from "@/lib/config";
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import Image from "next/image";
import Link from "next/link";
import { CreateGroupDialog } from "@/components/create-group-dialog";
import type { Group } from "@/types";
import { JoinGroupButton } from "@/components/join-group-button";
import { Button } from "@/components/ui/button";
import { getGroupCover } from "@/lib/avatars";
import { useUser } from "@/contexts/user-context";
import { useToast } from "@/hooks/use-toast";
import { Skeleton } from "@/components/ui/skeleton";
import { InviteUserDialog } from "@/components/invite-user-dialog";

export default function GroupsPage() {
    const [groups, setGroups] = useState<Group[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const { user } = useUser();
    const { toast } = useToast();

    useEffect(() => {
        const fetchGroups = async () => {
            setIsLoading(true);
            try {
                const res = await fetch(`${API_BASE_URL}/api/groups`, {
                    credentials: 'include'
                });
                if (!res.ok) {
                    throw new Error('Failed to fetch groups');
                }
                const data = await res.json();
                setGroups(Array.isArray(data) ? data : []);
            } catch (error) {
                console.error('Error fetching groups:', error);
                toast({ variant: 'destructive', title: 'Error', description: 'Could not load groups.' });
            } finally {
                setIsLoading(false);
            }
        };

        fetchGroups();
    }, [toast]);

    return (
        <main className="flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Discover Groups</h1>
                <CreateGroupDialog />
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {isLoading ? (
                    Array.from({ length: 3 }).map((_, i) => (
                        <Card key={i}>
                            <CardHeader>
                                <Skeleton className="aspect-video w-full" />
                                <Skeleton className="h-6 w-3/4 mt-4" />
                                <Skeleton className="h-4 w-full mt-2" />
                                <Skeleton className="h-4 w-full" />
                            </CardHeader>
                            <CardContent />
                            <CardFooter className="flex justify-between">
                                <Skeleton className="h-10 w-24" />
                                <Skeleton className="h-10 w-24" />
                            </CardFooter>
                        </Card>
                    ))
                ) : groups.length > 0 ? (
                    groups.map((group) => {
                        const cover = getGroupCover(group.id);
                        const isCreator = user?.id === group.creator_id;

                        return (
                            <Card key={group.id} className="flex flex-col">
                                <CardHeader>
                                    <div className="relative aspect-video w-full overflow-hidden rounded-t-lg">
                                        <Image 
                                            src={cover.url} 
                                            alt={group.name} 
                                            fill
                                            className="object-cover"
                                            data-ai-hint={cover.hint}
                                        />
                                    </div>
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
                                    {isCreator ? (
                                        <InviteUserDialog groupId={group.id} />
                                    ) : (
                                        <JoinGroupButton groupId={group.id} isPrivate={group.is_private} />
                                    )}
                                </CardFooter>
                            </Card>
                        );
                    })
                ) : (
                    <p className="col-span-full text-center text-muted-foreground py-8">No groups found. Why not create one?</p>
                )}
            </div>
        </main>
    );
}

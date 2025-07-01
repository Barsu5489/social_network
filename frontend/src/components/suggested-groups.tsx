

import Image from "next/image";
import { Button } from "./ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./ui/card";
import { Plus } from "lucide-react";

interface Group {
    id: string;
    name: string;
    description: string;
    creator_id: string;
    is_private: boolean;
    created_at: number;
    updated_at: number;
}

export function SuggestedGroups({ groups }: { groups: Group[] }) {

    return (
        <Card>
            <CardHeader>
                <CardTitle className="font-headline">Suggested Groups</CardTitle>
                <CardDescription>Find your community.</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4">
                {groups && groups.length > 0 ? (
                    groups.slice(0, 5).map((group) => ( // Show first 5 groups
                        <div key={group.id} className="flex items-center gap-4">
                            <Image src={"https://placehold.co/40x40.png"} alt={group.name} width={40} height={40} className="rounded-lg" data-ai-hint={"community event"} />
                            <div className="grid gap-1 flex-1">
                                <p className="font-semibold">{group.name}</p>
                                <p className="text-sm text-muted-foreground truncate">{group.description}</p>
                            </div>
                            <Button size="sm" variant="outline" disabled>
                                <Plus className="h-4 w-4 mr-1" />
                                Join
                            </Button>
                        </div>
                    ))
                ) : (
                    <p className="text-sm text-muted-foreground">Could not load suggested groups.</p>
                )}
            </CardContent>
        </Card>
    )
}

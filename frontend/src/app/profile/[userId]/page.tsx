
'use client';
import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { PostCard } from '@/components/post-card';
import { API_BASE_URL } from '@/lib/config';
import { useToast } from '@/hooks/use-toast';
import { Skeleton } from '@/components/ui/skeleton';
import { useUser } from '@/contexts/user-context';
import Link from 'next/link';
import { UserPlus, UserCheck, UserX, Lock, Edit, MessageCircle } from 'lucide-react';
import { EditProfileDialog } from '@/components/edit-profile-dialog';
import StartChatButton from '@/components/StartChatButton';


interface ProfileData {
    user: {
        id: string;
        email: string;
        first_name: string;
        last_name: string;
        nickname: string;
        date_of_birth: string;
        about_me: string;
        avatar_url: string;
        is_private: boolean;
        created_at: number;
    };
    posts: any[];
    followers: any[];
    following: any[];
    follower_count: number;
    following_count: number;
}

const ProfileSkeleton = () => (
    <div className="p-4 lg:p-6 space-y-6">
        <Card className="bg-card/50">
            <CardHeader className="flex flex-col md:flex-row items-start gap-6">
                <Skeleton className="h-32 w-32 rounded-full" />
                <div className="w-full space-y-3">
                    <Skeleton className="h-8 w-1/2" />
                    <Skeleton className="h-5 w-1/3" />
                    <Skeleton className="h-12 w-full" />
                    <div className="flex gap-4">
                        <Skeleton className="h-5 w-24" />
                        <Skeleton className="h-5 w-24" />
                    </div>
                </div>
            </CardHeader>
        </Card>
        <Skeleton className="h-10 w-full" />
    </div>
)

const FollowButton = ({ targetUserId, followers }: { targetUserId: string, followers: any[] }) => {
    const { user: currentUser } = useUser();
    const { toast } = useToast();
    const [isFollowing, setIsFollowing] = useState(() => (followers || []).some(f => f.id === currentUser?.id));
    const [isLoading, setIsLoading] = useState(false);
    const router = useRouter();

    const handleFollow = async () => {
        setIsLoading(true);
        try {
            const response = await fetch(`${API_BASE_URL}/api/users/${targetUserId}/follow`, { 
                method: 'POST', 
                credentials: 'include' 
            });
            if (!response.ok) throw new Error('Failed to follow');
            setIsFollowing(true);
            toast({ title: 'Followed successfully!' });
            router.refresh();
        } catch (error) {
            toast({ variant: 'destructive', title: 'Error', description: 'Could not follow user.' });
        } finally {
            setIsLoading(false);
        }
    };

    const handleUnfollow = async () => {
        setIsLoading(true);
        try {
            const response = await fetch(`${API_BASE_URL}/api/users/${targetUserId}/unfollow`, { 
                method: 'DELETE', 
                credentials: 'include' 
            });
            if (!response.ok) throw new Error('Failed to unfollow');
            setIsFollowing(false);
            toast({ title: 'Unfollowed successfully.' });
            router.refresh();
        } catch (error) {
            toast({ variant: 'destructive', title: 'Error', description: 'Could not unfollow user.' });
        } finally {
            setIsLoading(false);
        }
    };

    return isFollowing ? (
        <Button onClick={handleUnfollow} disabled={isLoading} variant="secondary">
            <UserCheck className="mr-2 h-4 w-4" /> Following
        </Button>
    ) : (
        <Button onClick={handleFollow} disabled={isLoading}>
            <UserPlus className="mr-2 h-4 w-4" /> Follow
        </Button>
    );
};


export default function ProfilePage() {
    const params = useParams();
    const userId = params.userId as string;
    const { toast } = useToast();
    const { user: currentUser } = useUser();

    const [profile, setProfile] = useState<ProfileData | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    const fetchProfile = async () => {
        if (!userId) return;
        setIsLoading(true);
        try {
            const response = await fetch(`${API_BASE_URL}/api/profile?target_id=${userId}`, {
                credentials: 'include'
            });
            const data = await response.json();
            if (response.ok && data.success) {
                setProfile(data.data);
            } else {
                setProfile(null);
                toast({ variant: 'destructive', title: 'Could not load profile.', description: data.error || "This user may not exist or their profile is private." });
            }
        } catch (error: any) {
            toast({ variant: 'destructive', title: 'Error', description: 'Failed to fetch profile data.' });
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchProfile();
    }, [userId]);
    
    if (isLoading) {
        return <ProfileSkeleton />;
    }

    if (!profile || !profile.user) {
        return (
             <main className="flex-1 flex flex-col items-center justify-center gap-4 p-4 lg:gap-6 lg:p-6 text-center">
                 <Card className="p-8 bg-card/50">
                    <Lock size={48} className="mx-auto text-muted-foreground mb-4"/>
                    <h2 className="text-xl font-semibold">Profile Not Available</h2>
                    <p className="text-muted-foreground">This profile is private or does not exist.</p>
                 </Card>
            </main>
        )
    }
    
    const { user, posts, followers, following, follower_count, following_count } = profile;
    const isOwnProfile = currentUser?.id === user.id;

    return (
        <main className="flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
            <Card className="bg-card/50">
                <CardHeader className="flex flex-col md:flex-row items-start gap-6 p-6 ">
                    <Avatar className="h-32 w-32 border-4 border-background ring-4 ring-primary">
                        <AvatarImage src={user.avatar_url || `https://i.pravatar.cc/128?u=${user.id}`} alt={`${user.first_name} ${user.last_name}`} />
                        <AvatarFallback className="text-4xl">{user.first_name?.[0]}{user.last_name?.[0]}</AvatarFallback>
                    </Avatar>
                    <div className="flex-1 space-y-1">
                        <div className="flex items-center justify-between">
                            <h1 className="text-2xl font-bold">{user.first_name} {user.last_name}</h1>
                            {isOwnProfile ? (
                                <EditProfileDialog profile={user} onProfileUpdate={fetchProfile} />
                            ) : (
                                <div className="flex gap-2">
                                    <FollowButton targetUserId={user.id} followers={followers} />
                                    <StartChatButton userId={user.id} />
                                </div>
                            )}
                        </div>
                        <p className="text-muted-foreground">@{user.nickname || user.id}</p>
                        <p className="pt-2">{user.about_me || 'No bio yet.'}</p>
                        <div className="flex items-center gap-4 text-sm pt-2">
                            <span><span className="font-bold">{follower_count}</span> Followers</span>
                            <span><span className="font-bold">{following_count}</span> Following</span>
                        </div>
                    </div>
                </CardHeader>
            </Card>

            <Tabs defaultValue="posts" className="mt-6">
                <TabsList>
                    <TabsTrigger value="posts">Posts</TabsTrigger>
                    <TabsTrigger value="followers">Followers</TabsTrigger>
                    <TabsTrigger value="following">Following</TabsTrigger>
                </TabsList>
                <TabsContent value="posts" className="mt-4">
                    {posts && posts.length > 0 ? (
                        posts.map(post => <PostCard key={post.id} {...post} />)
                    ) : (
                        <Card className="bg-card/50"><CardContent className="p-8 text-center text-muted-foreground">No posts yet.</CardContent></Card>
                    )}
                </TabsContent>
                <TabsContent value="followers" className="mt-4">
                     <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                        {followers && followers.length > 0 ? (
                            followers.map(follower => (
                                <Card key={follower.id} className="bg-card/50">
                                    <CardContent className="p-4 flex items-center gap-4">
                                        <Avatar>
                                            <AvatarImage src={`https://i.pravatar.cc/40?u=${follower.id}`} />
                                            <AvatarFallback>{follower.first_name?.[0]}</AvatarFallback>
                                        </Avatar>
                                        <div className="flex-1">
                                            <Link href={`/profile/${follower.id}`} className="font-semibold hover:underline">
                                                {follower.first_name} {follower.last_name}
                                            </Link>
                                            <p className="text-sm text-muted-foreground">@{follower.nickname}</p>
                                        </div>
                                    </CardContent>
                                </Card>
                            ))
                        ) : (
                           <Card className="sm:col-span-2 lg:col-span-3 bg-card/50"><CardContent className="p-8 text-center text-muted-foreground">No followers yet.</CardContent></Card>
                        )}
                    </div>
                </TabsContent>
                <TabsContent value="following" className="mt-4">
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                        {following && following.length > 0 ? (
                            following.map(followed => (
                                <Card key={followed.id} className="bg-card/50">
                                    <CardContent className="p-4 flex items-center gap-4">
                                        <Avatar>
                                            <AvatarImage src={`https://i.pravatar.cc/40?u=${followed.id}`} />
                                            <AvatarFallback>{followed.first_name?.[0]}</AvatarFallback>
                                        </Avatar>
                                        <div className="flex-1">
                                            <Link href={`/profile/${followed.id}`} className="font-semibold hover:underline">
                                                {followed.first_name} {followed.last_name}
                                            </Link>
                                            <p className="text-sm text-muted-foreground">@{followed.nickname}</p>
                                        </div>
                                    </CardContent>
                                </Card>
                            ))
                        ) : (
                           <Card className="sm:col-span-2 lg:col-span-3 bg-card/50"><CardContent className="p-8 text-center text-muted-foreground">Not following anyone yet.</CardContent></Card> 
                        )}
                    </div>
                </TabsContent>
            </Tabs>
        </main>
    )
}

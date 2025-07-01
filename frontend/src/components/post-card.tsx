
'use client';

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Image from 'next/image';
import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Heart, MessageCircle, MoreHorizontal, Globe, Users, Lock, Trash2 } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { cn } from "@/lib/utils";
import { API_BASE_URL } from "@/lib/config";
import { CommentSection } from "./comment-section";
import type { User } from "@/types";
import { useUser } from "@/contexts/user-context";
import { Skeleton } from "./ui/skeleton";

interface PostCardProps {
  id: string;
  user_id: string;
  content: string;
  privacy: 'public' | 'almost_private' | 'private';
  created_at: number; // unix timestamp
  likes_count: number;
  user_liked: boolean;
  imageUrl?: string;
  imageHint?: string;
}

export function PostCard({ id, user_id, content, privacy, created_at, likes_count, user_liked, imageUrl, imageHint }: PostCardProps) {
    const [isLiked, setIsLiked] = useState(user_liked);
    const [likeCount, setLikeCount] = useState(likes_count);
    const [showComments, setShowComments] = useState(false);
    const [author, setAuthor] = useState<User | null>(null);
    
    const router = useRouter();
    const { toast } = useToast();
    const { user: currentUser } = useUser();
    const isOwnPost = currentUser?.id === user_id;

    useEffect(() => {
        const fetchAuthor = async () => {
            try {
                const response = await fetch(`${API_BASE_URL}/api/profile?target_id=${user_id}`, {
                    credentials: 'include'
                });
                const data = await response.json();
                if (response.ok && data.success) {
                    setAuthor(data.data.user);
                }
            } catch (error) {
                console.error("Failed to fetch post author", error);
            }
        };

        if (user_id) {
            fetchAuthor();
        }
    }, [user_id]);

    const handleLikeToggle = async () => {
        const originalIsLiked = isLiked;
        const originalLikeCount = likeCount;

        const newIsLiked = !isLiked;
        const newLikeCount = newIsLiked ? likeCount + 1 : likeCount - 1;

        // Optimistically update UI
        setIsLiked(newIsLiked);
        setLikeCount(newLikeCount);

        try {
            const response = await fetch(`${API_BASE_URL}/posts/${id}/like`, {
                method: newIsLiked ? 'POST' : 'DELETE',
                credentials: 'include',
            });

            if (!response.ok) {
                // Revert UI on failure
                setIsLiked(originalIsLiked);
                setLikeCount(originalLikeCount);
                toast({
                    variant: "destructive",
                    title: "Something went wrong",
                    description: "Could not update like status. Please try again.",
                });
            }
        } catch (error) {
            // Revert UI on network error
            setIsLiked(originalIsLiked);
            setLikeCount(originalLikeCount);
            toast({
                variant: "destructive",
                title: "Network Error",
                description: "Could not connect to the server.",
            });
        }
    };
    
    const handleDelete = async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/delPost/${id}`, {
                method: 'DELETE',
                credentials: 'include',
            });

            if (response.ok) {
                toast({
                    title: "Post deleted",
                    description: "Your post has been successfully removed.",
                });
                router.refresh(); // Refresh the feed
            } else {
                 const errorData = await response.json().catch(() => ({ message: 'An unknown error occurred.' }));
                 throw new Error(errorData.error || 'Failed to delete post.');
            }
        } catch (error: any) {
             toast({
                variant: "destructive",
                title: "Error",
                description: error.message,
            });
        }
    };

    const PrivacyIcon = privacy === 'public' ? Globe : privacy === 'almost_private' ? Users : Lock;
    const formattedDate = new Date(created_at * 1000).toLocaleDateString(undefined, {
      year: 'numeric', month: 'short', day: 'numeric'
    });

  return (
    <Card>
      <CardHeader className="flex flex-row items-center gap-3 p-4">
        {author ? (
            <Link href={`/profile/${user_id}`}>
                <Avatar>
                    <AvatarImage src={author.avatar_url || `https://i.pravatar.cc/40?u=${user_id}`} alt={author.first_name} data-ai-hint="person portrait" />
                    <AvatarFallback>{author.first_name?.[0]}{author.last_name?.[0]}</AvatarFallback>
                </Avatar>
            </Link>
        ) : (
            <Skeleton className="h-10 w-10 rounded-full" />
        )}

        <div className="grid gap-0.5 flex-1">
          {author ? (
              <Link href={`/profile/${user_id}`} className="font-semibold hover:underline">
                <p className="truncate max-w-xs">{author.first_name} {author.last_name}</p>
              </Link>
          ) : (
              <Skeleton className="h-5 w-24" />
          )}
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <p>{formattedDate}</p>
            <span>&middot;</span>
            <PrivacyIcon className="h-4 w-4" title={privacy} />
          </div>
        </div>
        
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="ml-auto">
                    <MoreHorizontal className="h-5 w-5" />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
                {isOwnPost ? (
                    <DropdownMenuItem onClick={handleDelete} className="text-destructive focus:text-destructive focus:bg-destructive/10">
                        <Trash2 className="mr-2 h-4 w-4" />
                        <span>Delete</span>
                    </DropdownMenuItem>
                ) : (
                    <DropdownMenuItem>Report</DropdownMenuItem>
                )}
            </DropdownMenuContent>
        </DropdownMenu>

      </CardHeader>
      <CardContent className="px-4 pb-2">
        <p className="mb-4 whitespace-pre-wrap">{content}</p>
        {imageUrl && (
            <div className="relative aspect-[16/9] w-full overflow-hidden rounded-lg border">
                <Image src={imageUrl} alt="Post image" fill={true} style={{objectFit: 'cover'}} data-ai-hint={imageHint} />
            </div>
        )}
      </CardContent>
      <CardFooter className="flex justify-between p-4 pt-2 border-b">
        <div className="flex gap-4">
            <Button variant="ghost" className={cn("flex items-center gap-2 text-muted-foreground hover:text-destructive", isLiked && "text-destructive")} onClick={handleLikeToggle}>
                <Heart className={cn("h-5 w-5", isLiked && "fill-current")} />
                <span>{likeCount}</span>
            </Button>
            <Button variant="ghost" className="flex items-center gap-2 text-muted-foreground hover:text-primary" onClick={() => setShowComments(!showComments)}>
                <MessageCircle className="h-5 w-5" />
            </Button>
        </div>
      </CardFooter>
      {showComments && <CommentSection postId={id} />}
    </Card>
  );
}

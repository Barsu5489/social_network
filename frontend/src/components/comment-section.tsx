'use client';
import { useEffect, useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { useToast } from "@/hooks/use-toast";
import { API_BASE_URL } from "@/lib/config";
import { useUser } from "@/contexts/user-context";
import { Skeleton } from "./ui/skeleton";

interface Comment {
    id: string;
    post_id: string;
    user_id: string;
    content: string;
    image_url: string | null;
    created_at: number;
    user_nickname: string;
    user_avatar: string;
    likes_count: number;
    user_liked: boolean;
}

interface CommentSectionProps {
    postId: string;
}

export function CommentSection({ postId }: CommentSectionProps) {
    const [comments, setComments] = useState<Comment[]>([]);
    const [newComment, setNewComment] = useState("");
    const [isLoading, setIsLoading] = useState(true);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { toast } = useToast();
    const { user } = useUser();

    useEffect(() => {
        const fetchComments = async () => {
            setIsLoading(true);
            try {
                const response = await fetch(`${API_BASE_URL}/comments/${postId}`, {
                    credentials: 'include'
                });
                if (response.ok) {
                    const data = await response.json();
                    setComments(data.comments || []);
                } else {
                    toast({ variant: 'destructive', title: 'Failed to load comments.' });
                }
            } catch (error) {
                toast({ variant: 'destructive', title: 'Network error loading comments.' });
            } finally {
                setIsLoading(false);
            }
        };

        fetchComments();
    }, [postId, toast]);

    const handleCommentSubmit = async () => {
        if (!newComment.trim() || !user) return;
        setIsSubmitting(true);
        try {
            const response = await fetch(`${API_BASE_URL}/comment/${postId}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ content: newComment }),
                credentials: 'include'
            });

            if (response.ok) {
                const data = await response.json();
                setComments(prev => [...prev, data.comment]);
                setNewComment("");
            } else {
                toast({ variant: 'destructive', title: 'Failed to post comment.' });
            }
        } catch (error) {
             toast({ variant: 'destructive', title: 'Network error posting comment.' });
        } finally {
            setIsSubmitting(false);
        }
    };
    
    return (
        <div className="p-4 pt-2">
            <div className="flex items-start gap-4 mb-4">
                <Avatar className="h-9 w-9">
                    <AvatarImage src={user?.avatar_url || `https://i.pravatar.cc/40?u=${user?.id}`} />
                    <AvatarFallback>{user?.first_name?.[0]}</AvatarFallback>
                </Avatar>
                <div className="w-full">
                    <Textarea 
                        placeholder="Write a comment..."
                        value={newComment}
                        onChange={(e) => setNewComment(e.target.value)}
                        className="mb-2"
                        disabled={isSubmitting}
                    />
                    <div className="flex justify-end">
                        <Button onClick={handleCommentSubmit} disabled={isSubmitting || !newComment.trim()}>
                            {isSubmitting ? 'Posting...' : 'Post Comment'}
                        </Button>
                    </div>
                </div>
            </div>

            <div className="space-y-4">
                {isLoading ? (
                    Array.from({length: 2}).map((_, i) => (
                        <div key={i} className="flex items-start gap-4">
                            <Skeleton className="h-9 w-9 rounded-full" />
                            <div className="w-full space-y-2">
                                <Skeleton className="h-4 w-1/4" />
                                <Skeleton className="h-4 w-3/4" />
                            </div>
                        </div>
                    ))
                ) : comments.length > 0 ? (
                    comments.map(comment => (
                        <div key={comment.id} className="flex items-start gap-4 text-sm">
                            <Avatar className="h-9 w-9">
                                <AvatarImage src={comment.user_avatar || `https://i.pravatar.cc/40?u=${comment.user_id}`} />
                                <AvatarFallback>{comment.user_nickname?.[0] || 'U'}</AvatarFallback>
                            </Avatar>
                            <div className="flex-1">
                                <div className="bg-secondary rounded-lg px-3 py-2">
                                    <p className="font-semibold">{comment.user_nickname || comment.user_id}</p>
                                    <p>{comment.content}</p>
                                </div>
                                <div className="text-xs text-muted-foreground px-3 pt-1">
                                    {new Date(comment.created_at * 1000).toLocaleString()}
                                </div>
                            </div>
                        </div>
                    ))
                ) : (
                    <p className="text-sm text-center text-muted-foreground py-4">No comments yet. Be the first to comment!</p>
                )}
            </div>
        </div>
    );
}

'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { PostCard } from '@/components/post-card';
import { API_BASE_URL } from '@/lib/config';
import { Heart, MessageCircle, Users, Calendar } from 'lucide-react';
import Link from 'next/link';

interface UserActivityProps {
  userId: string;
  isOwnProfile: boolean;
}

interface ActivityItem {
  id: string;
  type: 'like' | 'comment' | 'follow' | 'post';
  created_at: number;
  post?: any;
  comment?: any;
  user?: any;
}

export function UserActivity({ userId, isOwnProfile }: UserActivityProps) {
  const [likedPosts, setLikedPosts] = useState([]);
  const [recentComments, setRecentComments] = useState([]);
  const [recentActivity, setRecentActivity] = useState<ActivityItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchActivity = async () => {
      try {
        // Fetch liked posts
        const likedResponse = await fetch(`${API_BASE_URL}/api/users/${userId}/liked-posts`, {
          credentials: 'include',
        });
        if (likedResponse.ok) {
          const likedData = await likedResponse.json();
          setLikedPosts(likedData.posts || []);
        }

        // Fetch recent comments (if own profile)
        if (isOwnProfile) {
          const commentsResponse = await fetch(`${API_BASE_URL}/api/users/${userId}/comments`, {
            credentials: 'include',
          });
          if (commentsResponse.ok) {
            const commentsData = await commentsResponse.json();
            setRecentComments(commentsData.comments || []);
          }
        }
      } catch (error) {
        console.error('Error fetching activity:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchActivity();
  }, [userId, isOwnProfile]);

  if (loading) {
    return (
      <Card className="bg-card/50">
        <CardContent className="p-8 text-center">
          <p>Loading activity...</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <Tabs defaultValue="liked" className="w-full">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="liked">Liked Posts</TabsTrigger>
          {isOwnProfile && <TabsTrigger value="comments">Recent Comments</TabsTrigger>}
          <TabsTrigger value="stats">Statistics</TabsTrigger>
        </TabsList>
        
        <TabsContent value="liked" className="mt-4">
          <div className="space-y-4">
            {likedPosts.length > 0 ? (
              likedPosts.map((post: any) => (
                <PostCard key={post.id} {...post} />
              ))
            ) : (
              <Card className="bg-card/50">
                <CardContent className="p-8 text-center text-muted-foreground">
                  <Heart className="mx-auto h-12 w-12 mb-4 opacity-50" />
                  <p>No liked posts yet.</p>
                </CardContent>
              </Card>
            )}
          </div>
        </TabsContent>

        {isOwnProfile && (
          <TabsContent value="comments" className="mt-4">
            <div className="space-y-4">
              {recentComments.length > 0 ? (
                recentComments.map((comment: any) => (
                  <Card key={comment.id} className="bg-card/50">
                    <CardContent className="p-4">
                      <div className="flex items-start gap-3">
                        <MessageCircle className="h-5 w-5 text-muted-foreground mt-1" />
                        <div className="flex-1">
                          <p className="text-sm">{comment.content}</p>
                          <p className="text-xs text-muted-foreground mt-1">
                            Commented on post • {new Date(comment.created_at * 1000).toLocaleDateString()}
                          </p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))
              ) : (
                <Card className="bg-card/50">
                  <CardContent className="p-8 text-center text-muted-foreground">
                    <MessageCircle className="mx-auto h-12 w-12 mb-4 opacity-50" />
                    <p>No recent comments.</p>
                  </CardContent>
                </Card>
              )}
            </div>
          </TabsContent>
        )}

        <TabsContent value="stats" className="mt-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card className="bg-card/50">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Likes Given</CardTitle>
                <Heart className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{likedPosts.length}</div>
              </CardContent>
            </Card>
            
            <Card className="bg-card/50">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Comments Made</CardTitle>
                <MessageCircle className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{recentComments.length}</div>
              </CardContent>
            </Card>

            <Card className="bg-card/50">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Member Since</CardTitle>
                <Calendar className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-sm font-medium">
                  {new Date(Date.now()).toLocaleDateString('en-US', { 
                    year: 'numeric', 
                    month: 'long' 
                  })}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
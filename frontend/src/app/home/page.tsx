
import { CreatePost } from "@/components/create-post";
import { PostCard } from "@/components/post-card";
import { SuggestedGroups } from "@/components/suggested-groups";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Card, CardContent } from "@/components/ui/card";
import { cookies } from 'next/headers';
import { API_BASE_URL } from "@/lib/config";

// This interface matches the backend response for a single post
interface Post {
  id: string;
  user_id: string;
  group_id: string | null;
  content: string;
  privacy: 'public' | 'almost_private' | 'private';
  created_at: number;
  updated_at: number;
  deleted_at: number | null;
  likes_count: number;
  user_liked: boolean;
}

interface Group {
    id: string;
    name: string;
    description: string;
    creator_id: string;
    is_private: boolean;
    created_at: number;
    updated_at: number;
}

async function getPosts(sessionCookie: ReturnType<typeof cookies>['get']) {
  if (!sessionCookie) return [];
  try {
    const res = await fetch(`${API_BASE_URL}/posts`, { 
      headers: {
        'Cookie': `${sessionCookie.name}=${sessionCookie.value}`,
      },
      cache: 'no-store' 
    });

    if (!res.ok) {
      console.error('Failed to fetch posts:', res.status, res.statusText);
      return [];
    }

    const data = await res.json();
    return Array.isArray(data) ? data : [];
  } catch (error) {
    console.error('Error fetching posts:', error);
    return [];
  }
}

async function getGroups(sessionCookie: ReturnType<typeof cookies>['get']) {
    if (!sessionCookie) return [];
    try {
        const res = await fetch(`${API_BASE_URL}/api/groups`, {
            headers: {
                'Cookie': `${sessionCookie.name}=${sessionCookie.value}`,
            },
            cache: 'no-store' 
        });
        if (!res.ok) {
            console.error('Failed to fetch groups:', res.statusText);
            return [];
        }
        const data = await res.json();
        return Array.isArray(data) ? data : [];
    } catch (error) {
        console.error('Error fetching groups:', error);
        return [];
    }
}


export default async function HomePage() {
  const cookieStore = cookies();
  const sessionCookie = cookieStore.get('social-network-session');

  const [posts, groups] = await Promise.all([
    getPosts(sessionCookie),
    getGroups(sessionCookie)
  ]);

  return (
        <main className="flex-1 p-4 lg:p-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
            <div className="lg:col-span-2">
              <ScrollArea className="h-[calc(100vh-100px)]">
                <div className="pr-4 space-y-6">
                  <CreatePost />
                  {posts.length > 0 ? (
                    posts.map((post) => (
                      <PostCard key={post.id} {...post} />
                    ))
                  ) : (
                    <Card>
                      <CardContent className="p-8 text-center text-muted-foreground">
                        <p>No posts found.</p>
                        <p className="text-sm">Be the first to post, or make sure you are logged in.</p>
                      </CardContent>
                    </Card>
                  )}
                </div>
              </ScrollArea>
            </div>
            <aside className="hidden lg:block lg:col-span-1">
              <div className="flex flex-col gap-6 sticky top-24">
                <SuggestedGroups groups={groups} />
              </div>
            </aside>
          </div>
        </main>
  );
}

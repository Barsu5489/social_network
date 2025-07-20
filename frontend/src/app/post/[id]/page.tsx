'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { PostCard } from '@/components/post-card'
import { Skeleton } from '@/components/ui/skeleton'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { ArrowLeft } from 'lucide-react'
import { API_BASE_URL } from '@/lib/config'

interface Post {
  id: string
  user_id: string
  content: string
  privacy: 'public' | 'almost_private' | 'private'
  created_at: number
  likes_count: number
  user_liked: boolean
}

export default function PostPage() {
  const params = useParams()
  const router = useRouter()
  const postId = params.id as string
  
  const [post, setPost] = useState<Post | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const url = `${API_BASE_URL}/api/posts/${postId}`
        console.log('DEBUG: Fetching post from:', url)
        console.log('DEBUG: API_BASE_URL:', API_BASE_URL)
        
        const response = await fetch(url, {
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          }
        })
        
        console.log('DEBUG: Response status:', response.status)
        console.log('DEBUG: Response headers:', response.headers)
        
        if (response.status === 404) {
          setError('Post not found or has been deleted')
          return
        }
        
        if (!response.ok) {
          const errorText = await response.text()
          console.error('DEBUG: Error response:', errorText)
          throw new Error(`Failed to fetch post: ${response.status} ${response.statusText}`)
        }
        
        const data = await response.json()
        console.log('DEBUG: Post data received:', data)
        setPost(data)
      } catch (err) {
        console.error('Error fetching post:', err)
        setError(
          `Failed to load post: ${
            typeof err === 'object' && err !== null && 'message' in err
              ? (err as { message: string }).message
              : String(err)
          }`
        )
      } finally {
        setLoading(false)
      }
    }

    if (postId) {
      fetchPost()
    }
  }, [postId])

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <Skeleton className="h-8 w-32 mb-4" />
        <Skeleton className="h-64 w-full" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <Button 
          variant="ghost" 
          onClick={() => router.back()}
          className="mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Go Back
        </Button>
        <Alert>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      </div>
    )
  }

  if (!post) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <Button 
          variant="ghost" 
          onClick={() => router.back()}
          className="mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Go Back
        </Button>
        <Alert>
          <AlertDescription>Post not found</AlertDescription>
        </Alert>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <Button 
        variant="ghost" 
        onClick={() => router.back()}
        className="mb-4"
      >
        <ArrowLeft className="h-4 w-4 mr-2" />
        Go Back
      </Button>
      <PostCard {...post} />
    </div>
  )
}

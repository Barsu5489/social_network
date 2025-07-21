
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { Image as ImageIcon, Video, Smile, Globe, Users, Lock } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { API_BASE_URL } from '@/lib/config';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { useUser } from '@/contexts/user-context';
import { FollowerSelector } from './follower-selector';

type Privacy = 'public' | 'almost_private' | 'private';

const emojis = ['😀', '😂', '😍', '🤔', '😢', '👍', '❤️', '🔥', '🎉', '🚀', '🤯', '🤬', '😡', '😠', '😤', '🥺', '😩', '🥳', '🤩', '🥸', '😎', '🥰', '🤓', '👏🏽', '🙌🏽', '👐🏾', '🫂', '🙅🏽‍♂️', '🤦🏽‍♂️', '💃', '🕺🏽', '🐶'];

// Add groupId to props
export function CreatePost({ groupId }: { groupId?: string }) {
  const [content, setContent] = useState('');
  const [privacy, setPrivacy] = useState<Privacy>('public');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [allowedUsers, setAllowedUsers] = useState<string[]>([]);
  const [showFollowerSelector, setShowFollowerSelector] = useState(false);
  const router = useRouter();
  const { toast } = useToast();
  const { user } = useUser();

  const handlePost = async () => {
    if (!content.trim() && !selectedImage) {
      return;
    }

    setIsSubmitting(true);
    
    try {
      let imageUrl = null;
      
      // Upload image if selected
      if (selectedImage) {
        const formData = new FormData();
        formData.append('image', selectedImage);
        
        const uploadResponse = await fetch(`${API_BASE_URL}/api/upload`, {
          method: 'POST',
          body: formData,
          credentials: 'include',
        });
        
        if (uploadResponse.ok) {
          const uploadData = await uploadResponse.json();
          imageUrl = uploadData.url;
        }
      }

      const isGroupPost = !!groupId;
      const url = isGroupPost ? `${API_BASE_URL}/api/groups/${groupId}/posts` : `${API_BASE_URL}/api/posts`;
      const body = isGroupPost 
        ? JSON.stringify({ content, image_url: imageUrl }) 
        : JSON.stringify({ content, privacy, image_url: imageUrl, allowed_user_ids: allowedUsers });

      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: body,
        credentials: 'include',
      });

      if (response.ok) {
        toast({
          title: "Post created!",
          description: "Your post is now live.",
        });
        setContent('');
        setPrivacy('public');
        setSelectedImage(null);
        setImagePreview(null);
        setAllowedUsers([]);
        setShowFollowerSelector(false);
        router.refresh();
      }
    } catch (error) {
      toast({
        variant: "destructive",
        title: "Network Error",
        description: "Could not connect to the server. Please try again later.",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const privacyOptions = {
    public: { icon: Globe, label: 'Public' },
    almost_private: { icon: Users, label: 'Followers' },
    private: { icon: Lock, label: 'Specific followers' },
  }

  const CurrentPrivacyIcon = privacyOptions[privacy].icon;

  const handleImageSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedImage(file);
      const reader = new FileReader();
      reader.onload = () => setImagePreview(reader.result as string);
      reader.readAsDataURL(file);
    }
  };

  const removeImage = () => {
    setSelectedImage(null);
    setImagePreview(null);
  };

  const handlePrivacyChange = (newPrivacy: Privacy) => {
    setPrivacy(newPrivacy);
    if (newPrivacy === 'private') {
      setShowFollowerSelector(true);
    } else {
      setShowFollowerSelector(false);
      setAllowedUsers([]);
    }
  };

  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-start gap-4">
          <Avatar>
            <AvatarImage src={user?.avatar_url || `https://i.pravatar.cc/40?u=${user?.id}`} alt={user?.first_name || 'User'} data-ai-hint="woman portrait"/>
            <AvatarFallback>{user?.first_name?.[0]}{user?.last_name?.[0]}</AvatarFallback>
          </Avatar>
          <div className="w-full">
            <Textarea
              placeholder={groupId ? `What's on your mind, member?` : `What's on your mind, ${user?.first_name || 'user'}?`}
              className="min-h-[60px] border-0 focus-visible:ring-0 focus-visible:ring-offset-0 bg-secondary/50"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              disabled={isSubmitting}
            />
            <div className="mt-2 flex items-center justify-between">
              <div className="flex gap-1 text-muted-foreground">
                <input
                  type="file"
                  accept="image/*,video/*,.gif"
                  onChange={handleImageSelect}
                  className="hidden"
                  id="image-upload"
                  disabled={isSubmitting}
                />
                <label htmlFor="image-upload">
                  <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-primary" type="button">
                    <ImageIcon className="h-5 w-5" />
                  </Button>
                </label>
                <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-primary" disabled>
                  <Video className="h-5 w-5" />
                </Button>
                <Popover>
                  <PopoverTrigger asChild>
                    <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-primary">
                      <Smile className="h-5 w-5" />
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-2">
                    <div className="grid grid-cols-5 gap-1">
                      {emojis.map((emoji) => (
                        <Button
                          key={emoji}
                          variant="ghost"
                          size="icon"
                          onClick={() => setContent(content + emoji)}
                          className="text-xl"
                        >
                          {emoji}
                        </Button>
                      ))}
                    </div>
                  </PopoverContent>
                </Popover>
              </div>
              <div className="flex items-center gap-2">
                {!groupId && (
                  <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                          <Button variant="outline" size="sm">
                              <CurrentPrivacyIcon className="h-4 w-4 mr-2"/>
                              {privacyOptions[privacy].label}
                          </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent>
                          <DropdownMenuItem onClick={() => setPrivacy('public')}>
                              <Globe className="h-4 w-4 mr-2"/>
                              Public
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => setPrivacy('almost_private')}>
                              <Users className="h-4 w-4 mr-2"/>
                              Followers
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => handlePrivacyChange('private')}>
                              <Lock className="h-4 w-4 mr-2"/>
                              Specific followers
                          </DropdownMenuItem>
                      </DropdownMenuContent>
                  </DropdownMenu>
                )}

                <Button onClick={handlePost} disabled={isSubmitting || !content.trim()}>
                    {isSubmitting ? 'Posting...' : 'Post'}
                </Button>
              </div>
            </div>
            {imagePreview && (
              <div className="mt-2 relative">
                <img src={imagePreview} alt="Preview" className="max-h-48 rounded-lg" />
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={removeImage}
                  className="absolute top-2 right-2"
                >
                  Remove
                </Button>
              </div>
            )}
            {showFollowerSelector && (
              <FollowerSelector
                selectedUsers={allowedUsers}
                onSelectionChange={setAllowedUsers}
              />
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}


'use client';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useRouter } from 'next/navigation';

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/hooks/use-toast";
import { API_BASE_URL } from '@/lib/config';
import { useUser } from '@/contexts/user-context';
import { Edit } from 'lucide-react';

interface Profile {
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
}

interface EditProfileDialogProps {
  profile: Profile;
  onProfileUpdate: () => void;
}

const profileSchema = z.object({
  nickname: z.string().optional(),
  about_me: z.string().max(200, "Bio must be 200 characters or less.").optional(),
  avatar_url: z.string().url("Must be a valid URL.").optional().or(z.literal('')),
  is_private: z.boolean().default(false),
});

export function EditProfileDialog({ profile, onProfileUpdate }: EditProfileDialogProps) {
  const [open, setOpen] = useState(false);
  const { toast } = useToast();
  const { setUser } = useUser();

  const form = useForm<z.infer<typeof profileSchema>>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      nickname: profile.nickname || "",
      about_me: profile.about_me || "",
      avatar_url: profile.avatar_url || "",
      is_private: profile.is_private,
    },
  });

  const onSubmit = async (values: z.infer<typeof profileSchema>) => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/profile`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
        credentials: 'include',
      });

      const responseData = await response.json();
      if (response.ok && responseData.success) {
        toast({ title: "Profile updated successfully!" });
        // Update user context if nickname or avatar changes, as it's used in the sidebar
        const updatedUser = { ...profile, ...values };
        setUser(updatedUser);
        onProfileUpdate(); // Re-fetch profile data on the page
        setOpen(false);
      } else {
        throw new Error(responseData.error || "Failed to update profile.");
      }
    } catch (error: any) {
      toast({ variant: 'destructive', title: "Error", description: error.message });
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline">
            <Edit className="mr-2 h-4 w-4"/>
            Edit Profile
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit profile</DialogTitle>
          <DialogDescription>
            Make changes to your profile here. Click save when you're done.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 py-4">
            <FormField
              control={form.control}
              name="nickname"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Nickname</FormLabel>
                  <FormControl>
                    <Input {...field} placeholder="Your public display name" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
             <FormField
              control={form.control}
              name="avatar_url"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Avatar URL</FormLabel>
                  <FormControl>
                    <Input {...field} placeholder="https://example.com/image.png" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="about_me"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>About Me</FormLabel>
                  <FormControl>
                    <Textarea {...field} placeholder="A short bio" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="is_private"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
                  <div className="space-y-0.5">
                    <FormLabel>Private Profile</FormLabel>
                    <FormDescription>
                        Only followers will see your full profile.
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
            <DialogFooter>
               <Button type="button" variant="secondary" onClick={() => setOpen(false)}>Cancel</Button>
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting ? 'Saving...' : 'Save changes'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

    

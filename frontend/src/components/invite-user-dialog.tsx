'use client';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { useToast } from "@/hooks/use-toast";
import { API_BASE_URL } from '@/lib/config';
import { Send } from 'lucide-react';

const inviteSchema = z.object({
  invitee_id: z.string().min(1, "User ID is required."),
});

interface InviteUserDialogProps {
  groupId: string;
}

export function InviteUserDialog({ groupId }: InviteUserDialogProps) {
  const [open, setOpen] = useState(false);
  const { toast } = useToast();

  const form = useForm<z.infer<typeof inviteSchema>>({
    resolver: zodResolver(inviteSchema),
    defaultValues: { invitee_id: "" },
  });

  const onSubmit = async (values: z.infer<typeof inviteSchema>) => {
    try {
      const payload = {
        ...values,
        entity_type: 'group',
        entity_id: groupId,
      };
      
      const response = await fetch(`${API_BASE_URL}/api/groups/invite`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
        credentials: 'include',
      });

      if (response.ok) {
        toast({ title: "Invitation sent successfully!" });
        setOpen(false);
        form.reset();
      } else {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to send invitation.");
      }
    } catch (error: any) {
      toast({ variant: 'destructive', title: "Error", description: error.message });
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline">
          <Send className="mr-2 h-4 w-4" />
          Invite
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Invite a user to the group</DialogTitle>
          <DialogDescription>
            Enter the User ID of the person you want to invite.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 py-4">
            <FormField
              control={form.control}
              name="invitee_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>User ID</FormLabel>
                  <FormControl>
                    <Input {...field} placeholder="Enter the user's ID" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
                <Button type="button" variant="secondary" onClick={() => setOpen(false)}>Cancel</Button>
                <Button type="submit" disabled={form.formState.isSubmitting}>
                    {form.formState.isSubmitting ? 'Sending...' : 'Send Invite'}
                </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

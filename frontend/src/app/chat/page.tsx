'use client';

import { ChatLayout } from "@/components/chat/chat-layout";
import { Card, CardContent } from "@/components/ui/card";
import { MessageCircle } from "lucide-react";

function EmptyChatView() {
  return (
    <div className="flex-1 flex items-center justify-center p-8">
      <Card className="bg-card/50 max-w-md w-full">
        <CardContent className="p-8 text-center">
          <MessageCircle className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
          <h3 className="text-lg font-semibold mb-2">No Chat Selected</h3>
          <p className="text-muted-foreground">
            Select a chat from the sidebar to start messaging, or start a new conversation from someone's profile.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}

export default function ChatPage() {
  return (
    <div className="h-[calc(100vh-theme(spacing.14))]">
      <ChatLayout defaultLayout={undefined} navCollapsedSize={8}>
        <EmptyChatView />
      </ChatLayout>
    </div>
  );
}

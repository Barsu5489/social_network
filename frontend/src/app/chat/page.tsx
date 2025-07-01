import { ChatLayout } from "@/components/chat/chat-layout";
import { cookies } from "next/headers";

export default function ChatPage() {
  const layout = cookies().get("react-resizable-panels:layout");
  const defaultLayout = layout ? JSON.parse(layout.value) : undefined;

  return (
      <div className="h-[calc(100vh-theme(spacing.14))]">
        <ChatLayout defaultLayout={defaultLayout} navCollapsedSize={8} />
      </div>
  );
}

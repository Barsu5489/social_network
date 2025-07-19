'use client';

import { useEffect, useState, useRef } from 'react';
import { useParams } from 'next/navigation';
import { useUser } from '@/contexts/user-context';
import { API_BASE_URL } from '@/lib/config';
import { useToast } from '@/hooks/use-toast';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { SendHorizontal } from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';
import { cn } from '@/lib/utils';
import { ChatLayout } from '@/components/chat/chat-layout';

// Interfaces for data structures
interface Message {
  id: string;
  chat_id: string;
  sender_id: string;
  content: string;
  sent_at: number;
  sender: {
    first_name: string;
    last_name: string;
    avatar_url: string;
  };
}

interface ChatDetails {
    id: string;
    type: string;
    name: string;
    avatar_url: string;
    participants: any[]; // Assuming an array of user objects
}

function ChatView({ chatId }: { chatId: string }) {
    const { user } = useUser();
    const { toast } = useToast();
    const [messages, setMessages] = useState<Message[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const [isSending, setIsSending] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [chatDetails, setChatDetails] = useState<ChatDetails | null>(null);
    const ws = useRef<WebSocket | null>(null);
    const messagesEndRef = useRef<HTMLDivElement | null>(null);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(() => {
        scrollToBottom();
    }, [messages]);


    // Fetch initial chat data
    useEffect(() => {
        if (!chatId) return;

        const fetchChatData = async () => {
            setIsLoading(true);
            try {
                // Fetch messages
                const messagesRes = await fetch(`${API_BASE_URL}/api/chats/${chatId}/messages`, { credentials: 'include' });
                if (!messagesRes.ok) {
                    console.error('Failed to load messages. Status:', messagesRes.status);
                    throw new Error('Failed to load messages.');
                }
                const messagesData = await messagesRes.json();
                console.log('Messages data:', messagesData);
                setMessages(messagesData.messages || []);

                // HACK: We don't have a dedicated /api/chats/{id} endpoint to get details.
                // We'll get the details from the main chats list.
                const chatsRes = await fetch(`${API_BASE_URL}/api/chats`, { credentials: 'include' });
                if(!chatsRes.ok) {
                    console.error('Failed to load chat details. Status:', chatsRes.status);
                    throw new Error('Failed to load chat details.');
                }
                const chatsData = await chatsRes.json();
                console.log('Chat details data:', chatsData);
                const currentChat = chatsData.chats.find((c: any) => c.id === chatId);
                console.log('Current chat details:', currentChat);
                setChatDetails(currentChat);

            } catch (error: any) {
                toast({ variant: 'destructive', title: 'Error', description: error.message });
            } finally {
                setIsLoading(false);
            }
        };

        fetchChatData();
    }, [chatId, toast]);

    // Setup WebSocket
    useEffect(() => {
        if (!user?.id || !chatId) return;

        console.log('Attempting WebSocket connection to:', `ws://localhost:3000/ws`);
        ws.current = new WebSocket(`ws://localhost:3000/ws`);
        
        ws.current.onopen = () => {
            console.log('WebSocket connected to backend');
            if (ws.current && user?.id) {
                ws.current.send(JSON.stringify({
                    type: 'join_chat',
                    chat_id: chatId,
                    user_id: user.id
                }));
            }
        };

        ws.current.onmessage = (event) => {
            console.log('WebSocket message received:', event.data);
            try {
                const messageData = JSON.parse(event.data);
                
                if (messageData.type === 'new_message' && messageData.chat_id === chatId) {
                    const newMessage = messageData.data;
                    setMessages((prevMessages) => {
                        if (prevMessages.some(m => m.id === newMessage.id)) {
                            return prevMessages;
                        }
                        return [...prevMessages, newMessage];
                    });
                } else if (messageData.type === 'notification') {
                    console.log('Received notification:', messageData.data);
                }
            } catch (error) {
                console.error('Failed to parse WebSocket message:', error);
            }
        };

        ws.current.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        ws.current.onclose = () => {
            console.log('WebSocket disconnected');
        };

        return () => {
            if (ws.current) {
                ws.current.close();
            }
        };
    }, [user?.id, chatId]);

    const handleSendMessage = async (messageContent: string) => {
        if (!user?.id || !chatId || !messageContent.trim() || isSending) return;
        
        setIsSending(true);
        
        try {
            const url = `${API_BASE_URL}/api/chats/${chatId}/messages`;
            
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    content: messageContent.trim(),
                    type: 'text'
                }),
            });

            if (!response.ok) {
                throw new Error(`Failed to send message: ${response.status}`);
            }

            const responseData = await response.json();
            
            if (responseData.success && responseData.message) {
                // Message will be added via WebSocket, but add locally for immediate feedback
                setMessages(prevMessages => {
                    if (prevMessages.some(m => m.id === responseData.message.id)) {
                        return prevMessages;
                    }
                    return [...prevMessages, responseData.message];
                });
                
                // Clear the input after successful send
                setNewMessage('');
            }
        } catch (error) {
            console.error('Error sending message:', error);
            toast({ 
                variant: 'destructive', 
                title: 'Error', 
                description: 'Failed to send message. Please try again.' 
            });
        } finally {
            setIsSending(false);
        }
    };

    if (isLoading) {
        return <div className="p-4"><Skeleton className="h-full w-full" /></div>;
    }

    if (!chatDetails) {
        return <div className="p-6 text-center">Chat not found.</div>;
    }

    return (
        <div className="flex flex-col h-full">
            {/* Chat Header */}
            <div className="flex items-center p-3 border-b">
                <Avatar className="h-10 w-10">
                    <AvatarImage src={chatDetails.avatar_url || `https://i.pravatar.cc/40?u=${chatDetails.id}`} alt={chatDetails.name} data-ai-hint="person portrait" />
                    <AvatarFallback>{chatDetails.name.substring(0, 2).toUpperCase()}</AvatarFallback>
                </Avatar>
                <div className="ml-4">
                    <p className="font-semibold">{chatDetails.name}</p>
                    <p className="text-sm text-muted-foreground">{chatDetails.type === 'group' ? `${chatDetails.participants?.length || 0} members` : 'Direct Message'}</p>
                </div>
            </div>

            {/* Messages Area */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
                {messages.map((message) => {
                    const isOwnMessage = message.sender_id === user?.id;
                    return (
                        <div key={message.id} className={cn("flex items-end gap-2", isOwnMessage && "justify-end")}>
                            {!isOwnMessage && (
                                <Avatar className="h-8 w-8">
                                    <AvatarImage src={message.sender.avatar_url || `https://i.pravatar.cc/40?u=${message.sender_id}`} />
                                    <AvatarFallback>{message.sender.first_name?.[0]}{message.sender.last_name?.[0]}</AvatarFallback>
                                </Avatar>
                            )}
                            <div className={cn(
                                "max-w-xs md:max-w-md lg:max-w-lg rounded-lg px-4 py-2 shadow-md",
                                isOwnMessage ? "bg-primary text-primary-foreground" : "bg-muted"
                            )}>
                                <p className="text-sm">{message.content}</p>
                            </div>
                             {isOwnMessage && (
                                <Avatar className="h-8 w-8">
                                    <AvatarImage src={user?.avatar_url || `https://i.pravatar.cc/40?u=${user?.id}`} />
                                    <AvatarFallback>{user?.first_name?.[0]}{user?.last_name?.[0]}</AvatarFallback>
                                </Avatar>
                            )}
                        </div>
                    );
                })}
                <div ref={messagesEndRef} />
            </div>

            {/* Message Input */}
            <div className="p-3 border-t">
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        if (newMessage.trim()) {
                            handleSendMessage(newMessage);
                        }
                    }}
                    className="flex items-center gap-2"
                >
                    <Input
                        value={newMessage}
                        onChange={(e) => setNewMessage(e.target.value)}
                        placeholder="Type a message..."
                        disabled={isSending}
                        autoComplete="off"
                        className="bg-secondary border-0 focus-visible:ring-1 focus-visible:ring-ring"
                        onKeyDown={(e) => {
                            if (e.key === 'Enter' && !e.shiftKey) {
                                e.preventDefault();
                                if (newMessage.trim()) {
                                    handleSendMessage(newMessage);
                                }
                            }
                        }}
                    />
                    <Button type="submit" size="icon" disabled={isSending || !newMessage.trim()}>
                        <SendHorizontal className="h-5 w-5" />
                    </Button>
                </form>
            </div>
        </div>
    );
}

// The page component that uses ChatLayout
export default function SingleChatPage() {
    const params = useParams();
    const chatId = params.chatId as string;
    
    // We can't access server-side cookies in a client component.
    // So we'll let the ChatLayout use its default.
    const defaultLayout = undefined; 
    
    return (
        <div className="h-[calc(100vh-theme(spacing.14))]">
            <ChatLayout defaultLayout={defaultLayout} navCollapsedSize={8}>
                <ChatView chatId={chatId} />
            </ChatLayout>
        </div>
    )
}
const [isLoading, setIsLoading] = useState(false);


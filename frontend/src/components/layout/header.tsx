import { Search, MessageSquare, LogOut } from "lucide-react";
import { Input } from "@/components/ui/input";
import { MobileNav } from "./mobile-nav";
import { NotificationBell } from "./notification-bell";
import Link from 'next/link';
import { Button } from "@/components/ui/button";
import { useUser } from "@/contexts/user-context";
import { API_BASE_URL } from "@/lib/config";

export function Header() {
  const { user, setUser } = useUser();

  const handleLogout = async () => {
    try {
      await fetch(`${API_BASE_URL}/api/logout`, {
        method: 'POST',
        credentials: 'include',
      });
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      setUser(null);
      localStorage.clear();
      sessionStorage.clear();
      window.location.href = '/';
    }
  };

  return (
    <header className="flex h-14 items-center gap-4 border-b bg-card px-4 lg:h-[60px] lg:px-6">
      <MobileNav />
      <div className="flex w-full flex-1 justify-center">
        <form className="w-full lg:w-1/2 xl:w-1/3">
          <div className="relative">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search posts, groups, or people..."
              className="w-full appearance-none bg-background pl-8 shadow-none"
            />
          </div>
        </form>
      </div>
      <Link href="/chat">
        <Button variant="ghost" size="icon" className="relative">
          <MessageSquare className="h-5 w-5" />
          <span className="sr-only">Chat</span>
        </Button>
      </Link>
      <NotificationBell />
      {user && (
        <Button 
          variant="outline" 
          size="sm" 
          onClick={handleLogout}
          className="flex items-center gap-2"
        >
          <LogOut className="h-4 w-4" />
          Logout
        </Button>
      )}
    </header>
  );
}

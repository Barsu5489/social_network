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
    <header className="flex h-14 items-center gap-4 border-b bg-card px-4">
      <MobileNav />
      <div className="flex w-full flex-1 justify-center">
        <form className="w-full max-w-md">
          <div className="relative">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search..."
              className="w-full pl-8"
            />
          </div>
        </form>
      </div>
      <Link href="/chat">
        <Button variant="ghost" size="icon">
          <MessageSquare className="h-5 w-5" />
        </Button>
      </Link>
      <NotificationBell />
      <Button variant="outline" size="sm" onClick={handleLogout}>
        <LogOut className="h-4 w-4 mr-2" />
        Logout
      </Button>
    </header>
  );
}

import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import { MobileNav } from "./mobile-nav";
import { NotificationBell } from "./notification-bell";

export function Header() {
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
            <NotificationBell />
        </header>
    );
}

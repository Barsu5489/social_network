import type { Metadata } from 'next';
import './globals.css';
import { Toaster } from "@/components/ui/toaster";
import { cn } from "@/lib/utils";
import { ThemeProvider } from '@/components/theme-provider';
import { UserProvider } from '@/contexts/user-context';
import { WebSocketProvider } from '@/contexts/websocket-context';

export const metadata: Metadata = {
  title: 'ConnectU',
  description: 'A modern social platform to connect with your world.',
};

const favicon = "data:image/svg+xml,%3Csvg%20width='42'%20height='42'%20viewBox='0%200%2042%2042'%20fill='none'%20xmlns='http://www.w3.org/2000/svg'%3E%3Cpath%20d='M34.75%205.25H7.25C6.14543%205.25%205.25%206.14543%205.25%207.25V26.25C5.25%2027.3546%206.14543%2028.25%207.25%2028.25H28L36.75%2036.75V7.25C36.75%206.14543%2035.8546%205.25%2034.75%205.25Z'%20fill='%235D9CEC'/%3E%3Ccircle%20cx='16'%20cy='15'%20r='5'%20fill='%234FC1E9'/%3E%3Cpath%20d='M11%2025C11%2022.2386%2013.2386%2020%2016%2020C18.7614%2020%2021%2022.2386%2021%2025H11Z'%20fill='%234FC1E9'/%3E%3Cg%20opacity='0.8'%3E%3Ccircle%20cx='26'%20cy='15'%20r='5'%20fill='%23AC92EC'/%3E%3Cpath%20d='M21%2025C21%2022.2386%2023.2386%2020%2026%2020C28.7614%2020%2031%2022.2386%2031%2025H21Z'%20fill='%23AC92EC'/%3E%3C/g%3E%3C/svg%3E";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <link rel="icon" href={favicon} type="image/svg+xml" />
      </head>
      <body className={cn("min-h-screen bg-background font-sans antialiased")} suppressHydrationWarning>
        <ThemeProvider
          attribute="class"
          defaultTheme="light"
          enableSystem={false}
          disableTransitionOnChange
        >
          <UserProvider>
            <WebSocketProvider>
              {children}
              <Toaster />
            </WebSocketProvider>
          </UserProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}

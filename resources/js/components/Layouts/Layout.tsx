"use client";

import {LogOut} from "lucide-react";
import * as React from "react";

import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import {Link} from "@inertiajs/react";
import {LanguageSelector} from "../forms/language-selector";
import {ThemeSwitcher} from "@/components/theme-switcher";
import {useTheme} from "@/hooks";
import {useEffect} from "react";
import {router} from "@inertiajs/react";
import {useTranslation} from "react-i18next";

interface NavItem {
  title: string;
  url: string;
  icon: React.ElementType;
  isActive?: boolean;
}

interface User {
  name: string;
  email: string;
  avatar: string;
}

interface SidebarLayoutProps {
  navItems: NavItem[];
  user: User;
  teamName: string;
  logoUrl: string;
  children: React.ReactNode;
}

export const SidebarLayout: React.FC<SidebarLayoutProps> = ({
                                                              navItems,
                                                              user,
                                                              logoUrl,
                                                              teamName,
                                                              children,
                                                            }: SidebarLayoutProps) => {
  const {theme, isDark, applyTheme} = useTheme(); // <- aquí usamos applyTheme directamente

  useEffect(() => {
    applyTheme(theme, isDark);
  }, [theme, isDark]);
  const {t} = useTranslation();

  return (
    <SidebarProvider>
      <div className="flex !h-screen w-full">
        <Sidebar collapsible="icon">
          <SidebarHeader className="flex items-center justify-center">
            <img
              src={logoUrl}
              alt={teamName}
              className="h-full w-auto transition-all duration-300 ease-in-out
                   group-data-[collapsible=icon]:w-10/12 group-data-[collapsible=icon]:h-full
                   group-data-[state=expanded]:w-10/12 group-data-[state=expanded]:mx-auto"
            />
          </SidebarHeader>
          <SidebarContent className="px-2">
            <SidebarMenu>
              {navItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild isActive={item.isActive}>
                    <button onClick={() => router.visit(item.url)} className="flex items-center">
                      <item.icon className="mr-2 h-4 w-4"/>
                      <span>{item.title}</span>
                    </button>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarContent>
          <SidebarFooter>
            <SidebarMenu>
              <SidebarMenuItem>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <SidebarMenuButton
                      size="lg"
                      variant="outline"
                      className="w-full justify-start gap-2 !bg-transparent hover:text-foreground">
                      <Avatar className="h-6 w-6">
                        <AvatarImage src={user.avatar} alt={user.name}/>
                        <AvatarFallback>{user.name.charAt(0)}</AvatarFallback>
                      </Avatar>
                      <div className="flex flex-col items-start text-left">
                        <span className="text-xs font-medium">{user.name}</span>
                        <span className="text-xs text-sidebar-foreground/60">
                          {user.email}
                        </span>
                      </div>
                    </SidebarMenuButton>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent
                    className="w-56"
                    align="start"
                    alignOffset={-8}
                    forceMount>
                    <DropdownMenuLabel>
                      <div className="flex flex-col space-y-1">
                        <p className="text-sm font-medium leading-none">
                          {user.name}
                        </p>
                        <p className="text-xs leading-none text-muted-foreground">
                          {user.email}
                        </p>
                      </div>
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator/>
                    <DropdownMenuItem className="p-0">
                      <LanguageSelector className="w-full p-2"/>
                    </DropdownMenuItem>
                    <DropdownMenuItem asChild>
                      <ThemeSwitcher/>
                    </DropdownMenuItem>
                    <DropdownMenuItem asChild>
                      <Link
                        href="/users/profile"
                        as="button"
                        className="w-full flex items-center">
                        <span>{t("profile")}</span>
                      </Link>
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarFooter>
        </Sidebar>
        <main className="flex-1 !w-full overflow-auto">
          <div className="h-screen">
            <div className="!w-full py-6">
              <div className="flex items-center justify-start gap-x-4 mb-4 pl-3">
                <SidebarTrigger/>
                {teamName && (
                  <span className="text-base lg:text-xl font-semibold">
                    {teamName} API Gateway
                  </span>
                )}
              </div>
              {children}
            </div>
          </div>
        </main>
      </div>
    </SidebarProvider>
  );
};

export default SidebarLayout;

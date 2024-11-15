'use client'
import * as React from 'react'
import { BadgeCheck, Bell, ChevronsUpDown, Command, CreditCard, LogOut, Sparkles } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarInset,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarProvider,
    SidebarRail,
    SidebarTrigger,
} from '@/components/ui/sidebar';

interface NavItem {
    title: string
    url: string
    icon: React.ElementType
    isActive?: boolean
}

interface User {
    name: string
    email: string
    avatar: string
}

interface SidebarLayoutProps {
    navItems: NavItem[]
    user: User
    teamName: string
    children: React.ReactNode
}

export default function SidebarLayout({ navItems, user, teamName, children }: SidebarLayoutProps) {
    return (
        <SidebarProvider className="min-h-screen">
            <div className="flex h-screen w-full">
                <Sidebar collapsible="icon" className="bg-primary text-primary-foreground">
                    <SidebarHeader>
                        <SidebarMenu>
                            <SidebarMenuItem>
                                <SidebarMenuButton
                                    size="lg"
                                    className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                                >
                                    <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                                        <Command className="size-4" />
                                    </div>
                                    <div className="grid flex-1 text-left text-sm leading-tight">
                                        <span className="truncate font-semibold">{teamName}</span>
                                    </div>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        </SidebarMenu>
                    </SidebarHeader>
                    <SidebarContent>
                        <SidebarGroup>
                            <SidebarGroupLabel>Navigation</SidebarGroupLabel>
                            <SidebarMenu>
                                {navItems.map((item) => (
                                    <SidebarMenuItem key={item.title}>
                                        <SidebarMenuButton tooltip={item.title} isActive={item.isActive}>
                                            {item.icon && <item.icon />}
                                            <span>{item.title}</span>
                                        </SidebarMenuButton>
                                    </SidebarMenuItem>
                                ))}
                            </SidebarMenu>
                        </SidebarGroup>
                    </SidebarContent>
                    <SidebarFooter>
                        <SidebarMenu>
                            <SidebarMenuItem>
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <SidebarMenuButton
                                            size="lg"
                                            className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                                        >
                                            <Avatar className="h-8 w-8 rounded-lg">
                                                <AvatarImage src={user.avatar} alt={user.name} />
                                                <AvatarFallback className="rounded-lg">
                                                    {user.name.slice(0, 2).toUpperCase()}
                                                </AvatarFallback>
                                            </Avatar>
                                            <div className="grid flex-1 text-left text-sm leading-tight">
                                                <span className="truncate font-semibold">{user.name}</span>
                                                <span className="truncate text-xs">{user.email}</span>
                                            </div>
                                            <ChevronsUpDown className="ml-auto size-4" />
                                        </SidebarMenuButton>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent
                                        className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                                        side="bottom"
                                        align="end"
                                        sideOffset={4}
                                    >
                                        <DropdownMenuLabel className="p-0 font-normal">
                                            <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                                                <Avatar className="h-8 w-8 rounded-lg">
                                                    <AvatarImage src={user.avatar} alt={user.name} />
                                                    <AvatarFallback className="rounded-lg">
                                                        {user.name.slice(0, 2).toUpperCase()}
                                                    </AvatarFallback>
                                                </Avatar>
                                                <div className="grid flex-1 text-left text-sm leading-tight">
                                                    <span className="truncate font-semibold">{user.name}</span>
                                                    <span className="truncate text-xs">{user.email}</span>
                                                </div>
                                            </div>
                                        </DropdownMenuLabel>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuGroup>
                                            <DropdownMenuItem>
                                                <Sparkles className="mr-2 h-4 w-4" />
                                                <span>Upgrade to Pro</span>
                                            </DropdownMenuItem>
                                        </DropdownMenuGroup>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuGroup>
                                            <DropdownMenuItem>
                                                <BadgeCheck className="mr-2 h-4 w-4" />
                                                <span>Account</span>
                                            </DropdownMenuItem>
                                            <DropdownMenuItem>
                                                <CreditCard className="mr-2 h-4 w-4" />
                                                <span>Billing</span>
                                            </DropdownMenuItem>
                                            <DropdownMenuItem>
                                                <Bell className="mr-2 h-4 w-4" />
                                                <span>Notifications</span>
                                            </DropdownMenuItem>
                                        </DropdownMenuGroup>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuItem>
                                            <LogOut className="mr-2 h-4 w-4" />
                                            <span>Log out</span>
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </SidebarMenuItem>
                        </SidebarMenu>
                    </SidebarFooter>
                    <SidebarRail />
                </Sidebar>
                <SidebarInset className="flex-1 h-screen overflow-auto bg-secondary/10">
                    <header className="flex h-16 shrink-0 items-center gap-2 border-b border-primary/20 px-6 bg-secondary/5">
                        <SidebarTrigger className="-ml-2" />
                        <h1 className="text-lg font-semibold">{teamName}</h1>
                    </header>
                    <main className="container mx-auto py-6">{children}</main>
                </SidebarInset>
            </div>
        </SidebarProvider>
    )
}
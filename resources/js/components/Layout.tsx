'use client'

import * as React from 'react'
import { ChevronLeft, ChevronRight, LogOut } from 'lucide-react'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
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
    useSidebar,
} from '@/components/ui/sidebar'

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
    logoUrl: string
    children: React.ReactNode
}

export default function SidebarLayout({ navItems, user, logoUrl, teamName, children }: SidebarLayoutProps) {
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
                    <SidebarContent>
                        <SidebarMenu>
                            {navItems.map((item) => (
                                <SidebarMenuItem key={item.title}>
                                    <SidebarMenuButton asChild isActive={item.isActive}>
                                        <a href={item.url} className="flex items-center">
                                            <item.icon className="mr-2 h-4 w-4" />
                                            <span>{item.title}</span>
                                        </a>
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
                                        <SidebarMenuButton size="lg" className="w-full justify-start gap-2">
                                            <Avatar className="h-6 w-6">
                                                <AvatarImage src={user.avatar} alt={user.name} />
                                                <AvatarFallback>{user.name.charAt(0)}</AvatarFallback>
                                            </Avatar>
                                            <div className="flex flex-col items-start text-left">
                                                <span className="text-xs font-medium">{user.name}</span>
                                                <span className="text-xs text-sidebar-foreground/60">{user.email}</span>
                                            </div>
                                        </SidebarMenuButton>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent className="w-56" align="start" alignOffset={-8} forceMount>
                                        <DropdownMenuLabel>My Account</DropdownMenuLabel>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuItem>Profile</DropdownMenuItem>
                                        <DropdownMenuItem>Settings</DropdownMenuItem>
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
                </Sidebar>
                <main className="flex-1 !w-full overflow-auto">
                    <div className="h-screen">
                        <div className="!w-full py-6">
                            <div className="flex items-center justify-start gap-x-4 mb-4 pl-3">
                                <SidebarTrigger />
                                {teamName && (
                                <span className="text-base lg:text-xl font-semibold">{teamName} Dashboard</span>
                                )}
                            </div>
                            {children}
                        </div>
                    </div>
                </main>
            </div>
        </SidebarProvider>
    )
}
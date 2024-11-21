import * as React from 'react'
import { AppWindowMac, DoorOpen, Frame, PieChartIcon, Settings2 } from 'lucide-react'
import SidebarLayout from "@/components/Layout"

const defaultNavMain = [
    {
        title: "Dashboard",
        url: "/dashboard",
        icon: PieChartIcon,
        isActive: true,
    },
    {
        title: "Applications",
        url: "#",
        icon: AppWindowMac,
    },
    {
        title: "Gateways",
        url: "#",
        icon: DoorOpen
    },
    {
        title: "Team",
        url: "/dashboard/team",
        icon: Frame,
    },
    {
        title: "Settings",
        url: "#",
        icon: Settings2,
    },
]

interface DashboardLayoutProps {
    children: React.ReactNode
    user?: {
        name: string
        email: string
        avatar: string
    }
    teamName?: string
    logoUrl?: string
}

export default function DashboardLayout({ 
    children,
    user = {
        name: "John Doe",
        email: "john@example.com",
        avatar: "https://unavatar.io/github/shadcn",
    },
    teamName = "Acme",
    logoUrl = "https://unavatar.io/github/shadcn",
}: DashboardLayoutProps) {
    return (
        <SidebarLayout
            navItems={defaultNavMain}
            user={user}
            teamName={teamName}
            logoUrl={logoUrl}
        >
            {children}
        </SidebarLayout>
    )
}

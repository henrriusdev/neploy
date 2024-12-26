import * as React from 'react'
import { AppWindowMac, DoorOpen, Frame, PieChartIcon, Settings2 } from 'lucide-react'
import SidebarLayout from "@/components/Layout"
import { Toaster } from '@/components/ui/toaster'
import { usePage, useRemember } from '@inertiajs/react'

const defaultNavMain = [
    {
        title: "Dashboard",
        url: "/dashboard",
        icon: PieChartIcon,
    },
    {
        title: "Applications",
        url: "/dashboard/applications",
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
  children: React.ReactNode;
  user?: {
    name: string;
    email: string;
    avatar: string;
    provider?: string;
    username?: string;
  };
  teamName?: string;
  logoUrl?: string;
}

export const DashboardLayout: React.FC<DashboardLayoutProps> = ({
    children,
    user,
    teamName = "Acme",
    logoUrl = "https://unavatar.io/github/shadcn",
}) => {
    const { url } = usePage()
    
    // Remember the layout data
    const [layoutData] = useRemember({
        user: {
            name: user?.name,
            email: user?.email,
            avatar: user?.provider === "github" 
                ? `https://unavatar.io/github/${user?.username}` 
                : `https://unavatar.io/${user?.email}`,
        },
        teamName,
        logoUrl,
    }, 'layout')

    // Create navigation with active state based on current URL
    const navigation = defaultNavMain.map(item => ({
        ...item,
        isActive: url.startsWith(item.url) && item.url !== "#"
    }))

    return (
        <SidebarLayout
            user={layoutData.user}
            teamName={layoutData.teamName}
            logoUrl={layoutData.logoUrl}
            navMain={navigation}
        >
            {children}
            <Toaster />
        </SidebarLayout>
    )
}

export default DashboardLayout;
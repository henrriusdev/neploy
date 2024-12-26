import * as React from 'react'
import { AppWindowMac, DoorOpen, Frame, PieChartIcon, Settings2 } from 'lucide-react'
import Layout from "@/components/Layout"
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
        url: "/dashboard/gateways",
        icon: DoorOpen
    },
    {
        title: "Team",
        url: "/dashboard/team",
        icon: Frame,
    },
    {
        title: "Settings",
        url: "/dashboard/settings",
        icon: Settings2,
    },
    {
        title: "Logout",
        url: "/logout",
        icon: DoorOpen,
    },
]

interface DashboardLayoutProps {
  children: React.ReactNode;
  user?: {
    name?: string;
    email?: string;
    username?: string;
    provider?: string;
  };
  teamName?: string;
  logoUrl?: string;
  navItems?: NavigationItem[];
}

interface NavigationItem {
  title: string;
  url: string;
  icon: any;
  isActive?: boolean;
}

export const DashboardLayout: React.FC<DashboardLayoutProps> = ({ children, user, teamName, logoUrl, navItems }) => {
    const { url } = usePage()
    
    const [layoutData] = useRemember({
        user: {
            name: user?.name || '',
            email: user?.email || '',
            avatar: user?.provider === "github" 
                ? `https://unavatar.io/github/${user?.username}` 
                : `https://unavatar.io/${user?.email}`,
        },
        teamName: teamName || '',
        logoUrl: logoUrl || '',
    }, 'dashboard-layout-state');

    // Create navigation with active state based on current URL
    const navigation = (navItems || defaultNavMain).map(item => ({
        ...item,
        isActive: url === item.url || (url.startsWith(item.url) && item.url !== "/dashboard" && item.url !== "/logout")
    }));

    return (
        <div className="min-h-screen bg-background">
            <Layout
                user={layoutData.user}
                teamName={layoutData.teamName}
                logoUrl={layoutData.logoUrl}
                navItems={navigation}
            >
                {children}
                <Toaster />
            </Layout>
        </div>
    );
}

export default DashboardLayout
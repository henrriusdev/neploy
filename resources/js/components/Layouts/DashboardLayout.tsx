import * as React from 'react'
import { AppWindowMac, DoorOpen, Frame, PieChartIcon, Settings2 } from 'lucide-react'
import Layout from "@/components/Layout"
import { Toaster } from '@/components/ui/toaster'
import { usePage, useRemember } from '@inertiajs/react'
import { useTranslation } from 'react-i18next'

const defaultNavMain = [
    {
        title: "sidebar.dashboard",
        url: "/dashboard",
        icon: PieChartIcon,
    },
    {
        title: "sidebar.applications",
        url: "/dashboard/applications",
        icon: AppWindowMac,
    },
    {
        title: "sidebar.gateways",
        url: "/dashboard/gateways",
        icon: DoorOpen
    },
    {
        title: "sidebar.team",
        url: "/dashboard/team",
        icon: Frame,
    },
    {
        title: "sidebar.settings",
        url: "/dashboard/settings",
        icon: Settings2,
    },
    {
        title: "sidebar.logout",
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
    const { t } = useTranslation()
    
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
        title: t(item.title), // Translate the title
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
            </Layout>
            <Toaster />
        </div>
    )
}

export default DashboardLayout
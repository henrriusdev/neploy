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
  children: React.ReactNode;
  user?: {
    name: string;
    email: string;
    avatar: string;
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

export default DashboardLayout;
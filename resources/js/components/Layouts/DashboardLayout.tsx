import Layout from "@/components/Layouts/Layout";
import { Toaster } from "@/components/ui/toaster";
import { usePage, useRemember } from "@inertiajs/react";
import {
  AppWindowMac,
  DoorOpen,
  Frame,
  PieChartIcon,
  Settings2,
} from "lucide-react";
import * as React from "react";
import { useTranslation } from "react-i18next";
import { Applications, Home } from "../views";
import { Gateways } from "../views/gateway";
import { Team } from "../views/team";

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
    icon: DoorOpen,
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
];

interface DashboardLayoutProps {
  user?: {
    name?: string;
    email?: string;
    username?: string;
    provider?: string;
  };
  teamName?: string;
  logoUrl?: string;
  navItems?: NavigationItem[];
  props?: any;
}

interface NavigationItem {
  title: string;
  url: string;
  icon: any;
  isActive?: boolean;
}

export const DashboardLayout: React.FC<DashboardLayoutProps> = ({
  user,
  teamName,
  logoUrl,
  navItems,
  props,
}) => {
  const { url } = usePage();
  const { t } = useTranslation();

  const [layoutData] = useRemember(
    {
      user: {
        name: user?.name || "",
        email: user?.email || "",
        avatar:
          user?.provider === "github"
            ? `https://unavatar.io/github/${user?.username}`
            : `https://unavatar.io/${user?.email}`,
      },
      teamName: teamName || "",
      logoUrl: logoUrl || "",
    },
    "dashboard-layout-state"
  );

  const getComponent = () => {
    const dashboardUrl = url.replace("/dashboard", "");
    switch (dashboardUrl) {
      case "":
        return <Home {...props} />;
      case "/applications":
        return <Applications {...props} />;
      case "/gateways":
        return <Gateways {...props} />;
      case "/team":
        return <Team {...props} />;
      case "/settings":
        return "Settings";
      default:
        return "Dashboard";
    }
  };

  // Create navigation with active state based on current URL
  const navigation = (navItems || defaultNavMain).map((item) => ({
    ...item,
    title: t(item.title), // Translate the title
    isActive:
      url === item.url ||
      (url.startsWith(item.url) &&
        item.url !== "/dashboard" &&
        item.url !== "/logout"),
  }));

  return (
    <div className="min-h-screen bg-background">
      <Layout
        user={layoutData.user}
        teamName={layoutData.teamName}
        logoUrl={layoutData.logoUrl}
        navItems={navigation}>
        {getComponent()}
      </Layout>
      <Toaster />
    </div>
  );
};

export default DashboardLayout;

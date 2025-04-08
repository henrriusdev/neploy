import Layout from "@/components/Layouts/Layout";
import { Toaster } from "@/components/ui/toaster";
import { usePage } from "@inertiajs/react";
import {
  AppWindowMac,
  DoorOpen,
  Frame,
  PieChartIcon,
  Settings2,
} from "lucide-react";
import * as React from "react";
import { useTranslation } from "react-i18next";
import { Applications, Home, Gateways, Team, Config } from "../views";
import {ApplicationView} from "@/components/views/application-view";

const navItems = [
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
  props?: any;
}

export const DashboardLayout: React.FC<DashboardLayoutProps> = ({
  user: backendUser,
  teamName,
  logoUrl,
  props,
}) => {
  const { url } = usePage();
  const { t } = useTranslation();

  const user = {
    name: backendUser?.name || "",
    email: backendUser?.email || "",
    avatar:
      backendUser?.provider === "github"
        ? `https://unavatar.io/github/${backendUser?.username}`
        : `https://unavatar.io/${backendUser?.email}`,
  };

  const getComponent = () => {
    const dashboardUrl = url.replace("/dashboard", "");
    if (dashboardUrl.startsWith("/applications/") && dashboardUrl.split("/").length === 3) {
      return <ApplicationView {...props} />;
    }

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
        return <Config {...props} />;
      default:
        return "Dashboard";
    }
  };

  // Create navigation with active state based on current URL
  const navigation = navItems.map((item) => ({
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
        user={user}
        teamName={teamName}
        logoUrl={logoUrl}
        navItems={navigation}>
        {getComponent()}
      </Layout>
      <Toaster />
    </div>
  );
};

export default DashboardLayout;

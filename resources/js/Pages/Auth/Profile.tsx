import {User} from "@/types";
import {UserProfile} from "@/components/views/user-profile";
import {AppWindowMac, DoorOpen, Frame, PieChartIcon, Settings2} from "lucide-react";
import Layout from "@/components/Layouts/Layout";
import {useTranslation} from "react-i18next";

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

export default function Profile({userData, user: backendUser, teamName, logoUrl}: any) {
  const user = {
    name: backendUser?.name || "",
    email: backendUser?.email || "",
    avatar:
      backendUser?.provider === "github"
        ? `https://unavatar.io/github/${backendUser?.username}`
        : `https://unavatar.io/${backendUser?.email}`,
  };

  const {t} = useTranslation()
  const navigation = navItems.map(item => ({
    ...item,
    title: t(item.title)
  }))

  return (
    <div className="min-h-screen bg-background">
      <Layout
        user={user}
        teamName={teamName}
        logoUrl={logoUrl}
        navItems={navigation}>
        <UserProfile user={userData as User}/>
      </Layout>
    </div>
  )
}
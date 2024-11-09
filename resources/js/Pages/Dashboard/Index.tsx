"use client";

import {
  BadgeCheck,
  Bell,
  ChevronsUpDown,
  Command,
  CreditCard,
  Folder,
  Frame,
  LogOut,
  PieChart as PieChartIcon,
  Settings2,
  Sparkles,
} from "lucide-react";
import * as React from "react";
import {
  Bar,
  BarChart,
  Cell,
  Line,
  LineChart,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
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
} from "@/components/ui/sidebar";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { techStackColors } from "@/lib/colors";

const defaultNavMain = [
  {
    title: "Dashboard",
    url: "#",
    icon: PieChartIcon,
    isActive: true,
  },
  {
    title: "Projects",
    url: "#",
    icon: Folder,
  },
  {
    title: "Team",
    url: "#",
    icon: Frame,
  },
  {
    title: "Settings",
    url: "#",
    icon: Settings2,
  },
];

const defaultRequestsData = [
  { name: "00:00", successful: 165, errors: 5 },
  { name: "01:00", successful: 120, errors: 4 },
  { name: "02:00", successful: 100, errors: 3 },
  { name: "03:00", successful: 140, errors: 3 },
  { name: "04:00", successful: 160, errors: 4 },
  { name: "05:00", successful: 170, errors: 6 },
  { name: "06:00", successful: 180, errors: 8 },
  { name: "07:00", successful: 200, errors: 10 },
  { name: "08:00", successful: 210, errors: 11 },
  { name: "09:00", successful: 220, errors: 12 },
  { name: "10:00", successful: 230, errors: 13 },
  { name: "11:00", successful: 240, errors: 14 },
  { name: "12:00", successful: 250, errors: 15 },
  { name: "13:00", successful: 260, errors: 16 },
  { name: "14:00", successful: 270, errors: 17 },
  { name: "15:00", successful: 280, errors: 18 },
  { name: "16:00", successful: 290, errors: 19 },
  { name: "17:00", successful: 250, errors: 8 },
  { name: "18:00", successful: 240, errors: 10 },
  { name: "19:00", successful: 230, errors: 9 },
  { name: "20:00", successful: 210, errors: 8 },
  { name: "21:00", successful: 200, errors: 7 },
  { name: "22:00", successful: 190, errors: 6 },
  { name: "23:00", successful: 180, errors: 5 },
  { name: "24:00", successful: 165, errors: 5 },
];

const defaultVisitorsData = [
  { name: "Mon", visitors: 2400 },
  { name: "Tue", visitors: 1398 },
  { name: "Wed", visitors: 9800 },
  { name: "Thu", visitors: 3908 },
  { name: "Fri", visitors: 4800 },
  { name: "Sat", visitors: 3800 },
  { name: "Sun", visitors: 4300 },
];

const defaultTechStackData = [
  { name: "React", value: 400 },
  { name: "Vue", value: 300 },
  { name: "Angular", value: 300 },
  { name: "Svelte", value: 200 },
];

const defaultUser = {
  name: "shadcn",
  email: "m@example.com",
  avatar: "/avatars/shadcn.jpg",
};

export default function Dashboard({
  navMain = defaultNavMain,
  requestData = defaultRequestsData,
  techStack = defaultTechStackData,
  user = defaultUser,
  primaryColor = "#8884d8",
  secondaryColor = "#82ca9d",
  visitorData = defaultVisitorsData,
}: {
  navMain?: Array<{
    title: string;
    url: string;
    icon: React.ElementType;
    isActive?: boolean;
  }>;
  requestData?: Array<{ name: string; successful: number; errors: number }>;
  techStack?: Array<{ name: string; value: number }>;
  user?: { name: string; email: string; avatar: string };
  primaryColor?: string;
  secondaryColor?: string;
  visitorData?: Array<{ name: string; visitors: number }>;
}) {
  return (
    <SidebarProvider
      className="!min-h-[90vh] !h-[90vh]"
      style={{
        "--primary-color": primaryColor,
        "--secondary-color": secondaryColor,
      }}>
      <div className="flex h-screen !w-full">
        <Sidebar
          collapsible="icon"
          className="bg-primary text-primary-foreground">
          <SidebarHeader>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton
                  size="lg"
                  className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground">
                  <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                    <Command className="size-4" />
                  </div>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">Dashboard</span>
                  </div>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarHeader>
          <SidebarContent>
            <SidebarGroup>
              <SidebarGroupLabel>Navigation</SidebarGroupLabel>
              <SidebarMenu>
                {navMain.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      tooltip={item.title}
                      isActive={item.isActive}>
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
                      className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground">
                      <Avatar className="h-8 w-8 rounded-lg">
                        <AvatarImage src={user.avatar} alt={user.name} />
                        <AvatarFallback className="rounded-lg">
                          {user.name.slice(0, 2).toUpperCase()}
                        </AvatarFallback>
                      </Avatar>
                      <div className="grid flex-1 text-left text-sm leading-tight">
                        <span className="truncate font-semibold">
                          {user.name}
                        </span>
                        <span className="truncate text-xs">{user.email}</span>
                      </div>
                      <ChevronsUpDown className="ml-auto size-4" />
                    </SidebarMenuButton>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent
                    className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                    side="bottom"
                    align="end"
                    sideOffset={4}>
                    <DropdownMenuLabel className="p-0 font-normal">
                      <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                        <Avatar className="h-8 w-8 rounded-lg">
                          <AvatarImage src={user.avatar} alt={user.name} />
                          <AvatarFallback className="rounded-lg">
                            {user.name.slice(0, 2).toUpperCase()}
                          </AvatarFallback>
                        </Avatar>
                        <div className="grid flex-1 text-left text-sm leading-tight">
                          <span className="truncate font-semibold">
                            {user.name}
                          </span>
                          <span className="truncate text-xs">{user.email}</span>
                        </div>
                      </div>
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuGroup>
                      <DropdownMenuItem>
                        <Sparkles />
                        Upgrade to Pro
                      </DropdownMenuItem>
                    </DropdownMenuGroup>
                    <DropdownMenuSeparator />
                    <DropdownMenuGroup>
                      <DropdownMenuItem>
                        <BadgeCheck />
                        Account
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <CreditCard />
                        Billing
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <Bell />
                        Notifications
                      </DropdownMenuItem>
                    </DropdownMenuGroup>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>
                      <LogOut />
                      Log out
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
            <h1 className="text-lg font-semibold">Dashboard</h1>
          </header>
          <main className="container mx-auto py-6">
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              <Card className="border-primary/10">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Health Apps
                  </CardTitle>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    className="h-4 w-4 text-muted-foreground">
                    <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" />
                  </svg>
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">7/10</div>
                  <p className="text-xs text-muted-foreground">
                    70% of apps are healthy
                  </p>
                </CardContent>
              </Card>
              <Card className="border-primary/10">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Total Requests
                  </CardTitle>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    className="h-4 w-4 text-muted-foreground">
                    <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
                  </svg>
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">1,675,234</div>
                  <p className="text-xs text-muted-foreground">
                    +18% from last month
                  </p>
                </CardContent>
              </Card>
              <Card className="border-primary/10">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Total Visitors
                  </CardTitle>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    className="h-4  w-4 text-muted-foreground">
                    <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                    <circle cx="9" cy="7" r="4" />
                    <path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" />
                  </svg>
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">573,281</div>
                  <p className="text-xs text-muted-foreground">
                    +201 since last hour
                  </p>
                </CardContent>
              </Card>
              <Card className="border-primary/10">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Total Errors
                  </CardTitle>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    className="h-4 w-4 text-muted-foreground">
                    <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
                  </svg>
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">78</div>
                  <p className="text-xs text-muted-foreground">
                    -5% from last hour
                  </p>
                </CardContent>
              </Card>
            </div>
            <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-7">
              <Card className="col-span-full border-primary/10">
                <CardHeader>
                  <CardTitle>Requests by Time</CardTitle>
                </CardHeader>
                <CardContent className="pl-2">
                  <ChartContainer
                    config={{
                      successful: {
                        label: "Successful",
                        color: "var(--primary-color)",
                      },
                      errors: {
                        label: "Errors",
                        color: "var(--secondary-color)",
                      },
                    }}
                    className="h-[350px] w-full">
                    <BarChart data={requestData}>
                      <XAxis dataKey="name" />
                      <YAxis />
                      <Bar
                        dataKey="successful"
                        stackId="a"
                        fill="#4faa4d"
                      />
                      <Bar
                        dataKey="errors"
                        stackId="a"
                        fill="#c00"
                      />
                      <ChartTooltip content={<ChartTooltipContent />} />
                    </BarChart>
                  </ChartContainer>
                </CardContent>
              </Card>
            </div>
            <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-7">
              <Card className="col-span-3 lg:col-span-3 border-primary/10">
                <CardHeader>
                  <CardTitle>Tech Stacks Most Used by Apps</CardTitle>
                </CardHeader>
                <CardContent>
                  <ChartContainer
                    config={techStack.reduce(
                      (acc, { name }, i) => ({
                        ...acc,
                        [name]: {
                          label: name,
                          color: techStackColors[i % techStackColors.length],
                        },
                      }),
                      {}
                    )}
                    className="h-[350px] flex justify-center items-center w-full">
                    <PieChart>
                      <Pie
                        data={techStack}
                        cx="50%"
                        cy="50%"
                        labelLine={true}
                        outerRadius={80}
                        dataKey="value"
                        nameKey="name">
                        {techStack.map((entry, index) => (
                          <Cell
                            key={`cell-${index}`}
                            fill={techStackColors[index % techStackColors.length]}
                          />
                        ))}
                      </Pie>
                      <ChartTooltip content={<ChartTooltipContent />} />
                    </PieChart>
                  </ChartContainer>
                  <div className="mt-4 flex justify-center">
                    {techStack.map((entry, index) => (
                      <div
                        key={`legend-${index}`}
                        className="mx-2 flex items-center">
                        <div
                          className="mr-2 h-3 w-3"
                          style={{
                            backgroundColor:
                              techStackColors[index % techStackColors.length],
                          }}
                        />
                        <span>{entry.name}</span>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
              <Card className="col-span-3 lg:col-span-4 border-primary/10">
                <CardHeader>
                  <CardTitle>Visitor Count by Time</CardTitle>
                </CardHeader>
                <CardContent>
                  <ChartContainer
                    config={{
                      visitors: {
                        label: "Visitors",
                        color: "var(--primary-color)",
                      },
                    }}
                    className="h-[350px]">
                    <LineChart data={visitorData}>
                      <XAxis dataKey="name" />
                      <YAxis />
                      <Line
                        type="monotone"
                        dataKey="visitors"
                        stroke="var(--primary-color)"
                        strokeWidth={2}
                      />
                      <ChartTooltip content={<ChartTooltipContent />} />
                    </LineChart>
                  </ChartContainer>
                </CardContent>
              </Card>
            </div>
          </main>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}

"use client";

import * as React from "react";
import {
  Bar,
  BarChart,
  Cell,
  Line,
  LineChart,
  Pie,
  PieChart,
  XAxis,
  YAxis,
} from "recharts";
import {
  AppWindowMac,
  DoorOpen,
  Frame,
  PieChartIcon,
  Settings2,
} from "lucide-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import DashboardLayout from "@/components/Layouts/DashboardLayout";
import { techStackColors } from "@/lib/colors";

const defaultRequestsData = [
  { name: "00:00", successful: 165, errors: 5 },
  { name: "04:00", successful: 193, errors: 5 },
  { name: "08:00", successful: 165, errors: 5 },
  { name: "20:00", successful: 369, errors: 1 },
  { name: "12:00", successful: 250, errors: 15 },
  { name: "16:00", successful: 402, errors: 15 },
  { name: "20:00", successful: 958, errors: 203 },
  { name: "24:00", successful: 165, errors: 5 },
];

const defaultVisitorsData = [
  { name: "Mon", visitors: 2400 },
  { name: "Wed", visitors: 9800 },
  { name: "Sun", visitors: 4300 },
];

const defaultTechStackData = [
  { name: "React", value: 400 },
  { name: "Vue", value: 250 },
  { name: "Angular", value: 300 },
  { name: "Svelte", value: 200 },
];

const defaultUser = {
  name: "shadcn",
  email: "m@example.com",
  avatar: "https://unavatar.io/github/shadcn",
  provider: "github",
  username: "shadcn",
};

interface DashboardProps {
  teamName?: string;
  navMain?: typeof defaultNavMain;
  requestData?: typeof defaultRequestsData;
  techStack?: typeof defaultTechStackData;
  user?: typeof defaultUser;
  primaryColor?: string;
  secondaryColor?: string;
  visitorData?: typeof defaultVisitorsData;
  health?: string;
  logoUrl?: string;
}

interface PageProps {
  props: DashboardProps
}

function Dashboard({
  requestData = defaultRequestsData,
  techStack = defaultTechStackData,
  user,
  primaryColor = "#8884d8",
  secondaryColor = "#82ca9d",
  visitorData = defaultVisitorsData,
  teamName,
  logoUrl,
  health = "4/10",
}: DashboardProps) {
  const healthPercentage =
    (parseInt(health?.split("/")[0]) / parseInt(health?.split("/")[1])) * 100;
  user.avatar = `https://unavatar.io/${user?.provider ?? "github"}/${user.username}`;

  const dashboardContent = (
    <>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card className="bg-card">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-primary-foreground">Health Apps</CardTitle>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              className="h-4 w-4 text-primary"
            >
              <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" />
            </svg>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary-foreground">{health}</div>
            <p className="text-xs text-primary/80">
              {healthPercentage}% are fully healthy
            </p>
          </CardContent>
        </Card>
        <Card className="bg-card">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-primary-foreground">
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
              className="h-4 w-4 text-primary"
            >
              <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
            </svg>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary-foreground">1,675,234</div>
            <p className="text-xs text-primary/80">
              +18% from last month
            </p>
          </CardContent>
        </Card>
        <Card className="bg-card">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-primary-foreground">
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
              className="h-4  w-4 text-primary"
            >
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary-foreground">573,281</div>
            <p className="text-xs text-primary/80">
              +201 since last hour
            </p>
          </CardContent>
        </Card>
        <Card className="bg-card">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-primary-foreground">Total Errors</CardTitle>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              className="h-4 w-4 text-primary"
            >
              <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
            </svg>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary-foreground">78</div>
            <p className="text-xs text-primary/80">-5% from last hour</p>
          </CardContent>
        </Card>
      </div>
      <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-full bg-card">
          <CardHeader>
            <CardTitle className="text-primary-foreground">Requests by Time</CardTitle>
          </CardHeader>
          <CardContent className="pl-2">
            <ChartContainer
              config={{
                successful: {
                  label: "Successful",
                  color: "hsl(var(--primary))",
                },
                errors: {
                  label: "Errors",
                  color: "hsl(var(--destructive))",
                },
              }}
              className="h-[350px] w-full"
            >
              <BarChart data={requestData}>
                <XAxis dataKey="name" />
                <ChartTooltip content={<ChartTooltipContent />} />
                <YAxis />
                <Bar
                  isAnimationActive={false}
                  dataKey="successful"
                  stackId="a"
                  fill="hsl(var(--primary))"
                />
                <Bar
                  dataKey="errors"
                  stackId="a"
                  fill="hsl(var(--destructive))"
                  isAnimationActive={false}
                />
              </BarChart>
            </ChartContainer>
          </CardContent>
        </Card>
      </div>
      <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-3 lg:col-span-3 bg-card">
          <CardHeader>
            <CardTitle className="text-primary-foreground">Tech Stacks Most Used by Apps</CardTitle>
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
                {},
              )}
              className="h-[350px] flex justify-center items-center w-full"
            >
              <PieChart>
                <Pie
                  data={techStack}
                  cx="50%"
                  cy="50%"
                  labelLine={true}
                  outerRadius={80}
                  dataKey="value"
                  nameKey="name"
                >
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
                <div key={`legend-${index}`} className="mx-2 flex items-center">
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
        <Card className="col-span-3 lg:col-span-4 bg-card border-none">
          <CardHeader>
            <CardTitle className="text-primary-foreground">Visitor Count by Time</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer
              config={{
                visitors: {
                  label: "Visitors",
                  color: "var(--primary-color)",
                },
              }}
              className="h-[350px]"
            >
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
    </>
  );

  return <div className="flex-1 space-y-4 p-8 pt-6">{dashboardContent}</div>;
}

Dashboard.layout = (page: any) => {
  const user = {
    name: page.props.user.name,
    email: page.props.user.email,
    avatar: page.props.user.provider === "github" ? "https://unavatar.io/github/" + page.props.user.username : "https://unavatar.io/" + page.props.user.email
  };
  return (
    <DashboardLayout
      user={user}
      teamName={page.props.teamName}
      logoUrl={page.props.logoUrl}
    >
      {page}
    </DashboardLayout>
  );
};

export default Dashboard;

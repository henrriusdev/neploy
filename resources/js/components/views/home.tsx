import { DashboardProps } from "@/types/props";
import { useTranslation } from "react-i18next";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { Button } from "../ui/button";
import { DashboardCard } from "../dashboard-card";
import { BaseChart } from "../base-chart";
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

export function Home({
  user,
  teamName,
  logoUrl,
  stats,
  requestData = defaultRequestsData,
  techStack = defaultTechStackData,
  visitorData = defaultVisitorsData,
  health = "4/10",
}: DashboardProps) {
  const { t } = useTranslation();

  const healthPercentage =
    (parseInt(health?.split("/")[0]) / parseInt(health?.split("/")[1])) * 100;
  user.avatar = `https://unavatar.io/${user?.provider ?? "github"}/${
    user.username
  }`;

  const dashboardContent = (
    <>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("dashboard.stats.totalApps")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats?.totalApps ?? 0}</div>
            <p className="text-xs text-muted-foreground">
              {t("dashboard.stats.totalAppsDescription")}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("dashboard.stats.runningApps")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats?.runningApps ?? 0}</div>
            <p className="text-xs text-muted-foreground">
              {t("dashboard.stats.runningAppsDescription")}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("dashboard.stats.teamMembers")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats?.teamMembers ?? 0}</div>
            <p className="text-xs text-muted-foreground">
              {t("dashboard.stats.teamMembersDescription")}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("dashboard.stats.deployments")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats?.deployments ?? 0}</div>
            <p className="text-xs text-muted-foreground">
              {t("dashboard.stats.deploymentsDescription")}
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="mt-8 grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>{t("dashboard.recentActivity.title")}</CardTitle>
            <CardDescription>
              {t("dashboard.recentActivity.description")}
            </CardDescription>
          </CardHeader>
          <CardContent>{/* Activity content */}</CardContent>
        </Card>

        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>{t("dashboard.resources.title")}</CardTitle>
            <CardDescription>
              {t("dashboard.resources.description")}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Button variant="outline" className="w-full justify-start">
              {t("dashboard.resources.documentation")}
            </Button>
            <Button variant="outline" className="w-full justify-start">
              {t("dashboard.resources.apiReference")}
            </Button>
            <Button variant="outline" className="w-full justify-start">
              {t("dashboard.resources.guides")}
            </Button>
            <Button variant="outline" className="w-full justify-start">
              {t("dashboard.resources.examples")}
            </Button>
          </CardContent>
        </Card>
      </div>

      <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <DashboardCard
          title={t("dashboard.healthApps")}
          value={health}
          description={`${healthPercentage}% ${t("dashboard.fullyHealthy")}`}
          icon={
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              className="h-4 w-4 text-primary">
              <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" />
            </svg>
          }
        />
        <DashboardCard
          title={t("dashboard.totalRequests")}
          value="1,675,234"
          description={t("dashboard.requestsLastMonth")}
          icon={
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              className="h-4 w-4 text-primary">
              <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
            </svg>
          }
        />
        <DashboardCard
          title={t("dashboard.totalVisitors")}
          value="573,281"
          description={t("dashboard.visitorsLastHour")}
          icon={
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              className="h-4  w-4 text-primary">
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          }
        />
        <DashboardCard
          title={t("dashboard.totalErrors")}
          value="78"
          description={t("dashboard.errorsLastHour")}
          icon={
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              className="h-4 w-4 text-primary">
              <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
            </svg>
          }
        />
      </div>
      <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <BaseChart
          title={t("dashboard.requestsByTime")}
          data={requestData}
          type="bar"
          dataKeys={["successful", "errors"]}
          colors={["hsl(var(--primary))", "hsl(var(--destructive))"]}
          className="col-span-full"
          config={{
            successful: {
              label: t("dashboard.successful"),
              color: "hsl(var(--primary))",
            },
            errors: {
              label: t("dashboard.errors"),
              color: "hsl(var(--destructive))",
            },
          }}
        />
      </div>
      <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <BaseChart
          title={t("dashboard.techStacksMostUsed")}
          data={techStack}
          type="pie"
          dataKeys={["value"]}
          colors={techStackColors}
          className="col-span-3 lg:col-span-3"
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
        />
        <BaseChart
          title={t("dashboard.visitorCountByTime")}
          data={visitorData}
          type="line"
          dataKeys={["visitors"]}
          colors={["var(--primary-color)"]}
          className="col-span-3 lg:col-span-4 border-none"
          config={{
            visitors: {
              label: t("dashboard.visitors"),
              color: "var(--primary-color)",
            },
          }}
        />
      </div>
    </>
  );

  return <div className="flex-1 space-y-4 p-8 pt-6">{dashboardContent}</div>;
}

import {DashboardProps} from "@/types/props";
import {useTranslation} from "react-i18next";
import {Card, CardContent, CardDescription, CardHeader, CardTitle,} from "../ui/card";
import {Button} from "../ui/button";
import {DashboardCard} from "../dashboard-card";
import {BaseChart} from "../base-chart";
import {techStackColors} from "@/lib/colors";
import {useEffect, useState} from "react";
import {Skeleton} from "@/components/ui/skeleton";

const defaultVisitorsData = [
  {name: "Mon", visitors: 2400},
  {name: "Wed", visitors: 9800},
  {name: "Sun", visitors: 4300},
];

export function Home({
                       requests,
                       techStack,
                       visitorData = defaultVisitorsData,
                       health = "4/10",
                     }: DashboardProps) {
  const {t} = useTranslation();
  const [totalRequests, setTotalRequests] = useState(0);
  const [totalErrors, setTotalErrors] = useState(0);

  requests = requests?.map((request) => ({
    ...request,
    total: request.successful + request.errors,
    successful: request.successful - request.errors,
    errors: request.errors ?? 0,
  }));

  useEffect(() => {
    if (requests) {
      const totalSuccessful = requests.reduce(
        (acc, request) => acc + request.total,
        0
      );
      const totalErrors = requests.reduce(
        (acc, request) => acc + request.errors,
        0
      );
      setTotalRequests(totalSuccessful);
      setTotalErrors(totalErrors);
    }
  }, [requests]);

  return <div className="flex-1 space-y-4 p-8 pt-6">
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
            <path
              d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"/>
          </svg>
        }
      />
      <DashboardCard
        title={t("dashboard.totalRequests")}
        value={totalRequests.toString()}
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
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
          </svg>
        }
      />
      <DashboardCard
        title={t("dashboard.totalVisitors")}
        value="573,281"
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
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
        }
      />
      <DashboardCard
        title={t("dashboard.totalErrors")}
        value={totalErrors}
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
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
          </svg>
        }
      />
    </div>
    <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <BaseChart
        title={t("dashboard.requestsByTime")}
        data={requests}
        type="bar"
        dataKeys={["successful", "errors"]}
        colors={["hsl(var(--primary))", "hsl(var(--destructive))"]}
        className="col-span-full"
      />
    </div>
    <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      {techStack?.length === 0 ? (
        <BaseChart
          title={t("dashboard.techStacksMostUsed")}
          data={techStack}
          type="pie"
          dataKeys={["value"]}
          colors={techStackColors}
          className="col-span-3 lg:col-span-3"
        /> ) : (
        <Skeleton className="col-span-3 lg:col-span-3 h-[300px]"/>
      )}
      <BaseChart
        title={t("dashboard.visitorCountByTime")}
        data={defaultVisitorsData}
        type="line"
        dataKeys={["visitors"]}
        colors={["var(--primary)"]}
        className="col-span-3 lg:col-span-4 border-none"
      />
    </div>
  </div>;
}

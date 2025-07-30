import {useMemo, useState} from "react";
import {Bar, BarChart, Cell, Line, LineChart, Pie, PieChart, XAxis, YAxis} from "recharts";
import {Card, CardContent, CardHeader} from "@/components/ui/card";
import {Checkbox} from "@/components/ui/checkbox";
import {Skeleton} from "@/components/ui/skeleton";
import {Label} from "@/components/ui/label";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import {DateRange} from "react-day-picker";
import {DatePicker} from "@/components/forms/date-picker";
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent
} from "@/components/ui/chart";
import {format, parseISO, isValid} from "date-fns";
import {Button} from "../ui/button";
import {Theme, useTheme} from "@/hooks";
import {BaseChart} from "../base-chart";
import {techStackColors} from "@/lib/colors";
import {RequestData, StackData, VisitorData} from "@/types/props";
import {useTranslation} from "react-i18next";

interface ApplicationStat {
  application_id: string;
  date: string;
  requests: number;
  errors: number;
  name?: string;
}

export function Reports({stats, requests, techStack, visitors}: {
  stats: ApplicationStat[];
  requests?: RequestData[];
  techStack?: StackData[];
  visitors?: VisitorData[];
}) {
  const {t} = useTranslation();
  const {applyTheme} = useTheme();
  const [selectedMetrics, setSelectedMetrics] = useState<string[]>(["requests"]);
  const [dateRange, setDateRange] = useState<DateRange | undefined>(undefined);
  const [appFilter, setAppFilter] = useState<string>("all");
  const [chartType, setChartType] = useState<"line" | "bar">("line");
  console.log(techStack)

  const apps = useMemo(() => (stats ? Array.from(new Map(stats.map((s) => [s.application_id, {
    id: s.application_id,
    name: s.name
  }])).values()) : []), [stats]);

  const filteredData = useMemo(() => {
    if (!stats) return [];
    return stats.filter((stat) => {
      const date = new Date(stat.date);
      const inRange = (!dateRange?.from || date >= dateRange.from) && (!dateRange?.to || date <= dateRange.to);
      const appMatch = appFilter === "all" || stat.application_id === appFilter;
      return inRange && appMatch;
    });
  }, [stats, dateRange, appFilter]);

  // Instead of grouping, just map stats to include a 'name' field for recharts, showing only the hour (HH:mm)
  const chartData = useMemo(() => {
    return filteredData
      .map((stat) => {
        let hourLabel = stat.date;
        try {
          const d = parseISO(stat.date);
          if (isValid(d)) hourLabel = format(d, "HH:mm");
        } catch {
        }
        return {
          ...stat,
          name: hourLabel,
        };
      })
      .sort((a, b) => a.name.localeCompare(b.name));
  }, [filteredData]);

  // Filter visitors data by date range
  const filteredVisitors = useMemo(() => {
    if (!visitors) return [];
    return visitors.filter((visitor) => {
      if (!dateRange?.from && !dateRange?.to) return true;
      try {
        const date = parseISO(visitor.name);
        if (!isValid(date)) return true; // If can't parse date, include it
        const inRange = (!dateRange?.from || date >= dateRange.from) && (!dateRange?.to || date <= dateRange.to);
        return inRange;
      } catch {
        return true; // If error parsing, include it
      }
    });
  }, [visitors, dateRange]);

  const toggleMetric = (key: string) => {
    setSelectedMetrics((prev) => (prev.includes(key) ? prev.filter((k) => k !== key) : [...prev, key]));
  };

  const Chart = chartType === "bar" ? BarChart : LineChart;
  const Series = chartType === "bar" ? Bar : Line;

  // Use translation for metric labels
  const metrics = [
    {key: "requests", label: t("dashboard.successful"), color: "#8884d8"},
    {key: "errors", label: t("dashboard.errors"), color: "#ff7300"},
  ];

  const config = Object.fromEntries(metrics.map((m) => [m.key, {label: m.label, color: m.color}]));

  return (
    <Card className="print:shadow-none print:border-none print:bg-white print:p-1">
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-8 gap-4 print:hidden">
          <div className="space-y-1 md:col-span-2">
            <Label>Aplicación</Label>
            <Select onValueChange={setAppFilter} value={appFilter}>
              <SelectTrigger>
                <SelectValue/>
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todas</SelectItem>
                {apps.map((app) => (
                  <SelectItem key={app.id} value={app.id}>
                    {app.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-1 md:col-span-2">
            <Label>Rango de fechas</Label>
            <DatePicker isRangePicker maxYear={new Date().getFullYear()} date={dateRange} onDateChange={setDateRange}/>
          </div>
          <div className="space-y-1 md:col-span-2">
            <Label>Tipo de gráfico</Label>
            <Select onValueChange={(value) => setChartType(value as "line" | "bar")} value={chartType}>
              <SelectTrigger>
                <SelectValue/>
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="line">Línea</SelectItem>
                <SelectItem value="bar">Barras</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="space-x-4 flex items-center justify-center md:col-span-2 2xl:col-span-1">
            {metrics.map((m) => (
              <Label key={m.key} className="flex items-center space-x-2">
                <Checkbox checked={selectedMetrics.includes(m.key)} onCheckedChange={() => toggleMetric(m.key)}/>
                <span>{m.label}</span>
              </Label>
            ))}
          </div>
          {stats && stats.length > 0 && (
            <Button
              onClick={() => {
                // Save current theme
                const currentTheme = localStorage.getItem("theme") || "system";
                const currentDark = localStorage.getItem("darkMode") === "true";

                // Switch to light theme for printing
                applyTheme("neploy", false); // Using 'neploy' as the light theme

                // Trigger print
                setTimeout(() => {
                  window.print();

                  // Restore original theme after printing
                  setTimeout(() => {
                    applyTheme(currentTheme as Theme, currentDark);
                  }, 500);
                }, 300);
              }}
              className="w-fit flex items-center gap-2 place-self-center print:hidden md:col-span-1 mx-auto">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="mr-2 h-4 w-4">
                <polyline points="6 9 6 2 18 2 18 9"></polyline>
                <path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"></path>
                <rect x="6" y="14" width="12" height="8"></rect>
              </svg>
              <span>Imprimir</span>
            </Button>
          )}
        </div>
        {/* Responsive chart wrapper to prevent horizontal scroll */}
        <h2
          className="font-semibold leading-none tracking-tight text-foreground mt-4">{t("dashboard.requestsByTime")}</h2>
        <div style={{width: "100%", overflowX: "auto"}} className="print:mt-0 py-2 print:p-0">
          <div style={{minWidth: 0}} className="print:w-full">
            {stats ? (
              chartData.length > 0 ? (
                <ChartContainer config={config} className="print:w-full print:max-w-full">
                  <Chart data={chartData} margin={{top: 20, right: 30, left: 20, bottom: 5}}
                         className="print:w-full print:max-w-full">
                    <XAxis dataKey="name"/>
                    <YAxis/>
                    <ChartTooltip content={<ChartTooltipContent/>}/>
                    <ChartLegend content={<ChartLegendContent/>}/>
                    {selectedMetrics.map((metric) => {
                      const m = metrics.find((m) => m.key === metric);
                      return (
                        <Series
                          key={metric}
                          type="monotone"
                          dataKey={metric}
                          stroke={chartType === "line" ? m?.color : undefined}
                          fill={chartType === "bar" ? m?.color : undefined}
                          stackId={chartType === "bar" ? "a" : undefined}
                          strokeWidth={2}
                          dot={false}
                        />
                      );
                    })}
                  </Chart>
                </ChartContainer>
              ) : (
                <div className="flex items-center justify-center h-[300px]">
                </div>
              )
            ) : (
              <Skeleton className="h-[300px] w-full"/>
            )}
          </div>
        </div>

        {/* Dashboard Charts - Tech Stack and Visitors */}
        <div className="mt-4 grid gap-4 md:grid-cols-2 print:grid-cols-1">
          {/* Tech Stack Chart */}
          <div>
            {techStack ? (
              techStack.length > 0 ? (
                <BaseChart
                  title={t("dashboard.techStacksMostUsed")}
                  data={techStack}
                  type="pie"
                  dataKeys={["value"]}
                  labelKey="name"
                  colors={techStackColors}
                  className="print:w-full print:mb-8"
                />
              ) : (
                <Card className="flex items-center justify-center h-[300px]">
                </Card>
              )
            ) : (
              <Skeleton className="h-[300px] w-full print:w-full print:mb-8"/>
            )}
          </div>

          {/* Visitors Chart */}
          <div>
            {filteredVisitors ? (
              filteredVisitors.length > 0 ? (
                <BaseChart
                  title={t("dashboard.visitorCountByTime")}
                  data={filteredVisitors}
                  type="line"
                  dataKeys={["value"]}
                  colors={["var(--primary)"]}
                  className="border-none print:w-full print:mb-8"
                />
              ) : (
                <Card className="flex items-center justify-center h-[300px]">
                </Card>
              )
            ) : (
              <Skeleton className="h-[300px] w-full print:w-full print:mb-8"/>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

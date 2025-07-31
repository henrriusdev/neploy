import {useMemo, useState} from "react";
import {Bar, BarChart, Cell, Line, LineChart, Pie, PieChart, ResponsiveContainer, XAxis, YAxis} from "recharts";
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
import {format, parseISO, isValid, addDays} from "date-fns";
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
  const [dateRange, setDateRange] = useState<DateRange>({from: addDays(new Date(), -7), to: new Date()});
  const [appFilter, setAppFilter] = useState<string>("all");
  const [chartType, setChartType] = useState<"line" | "bar">("line");
  const [visitorChartType, setVisitorChartType] = useState<"line" | "bar">("line");
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

  // Map stats to include a proper formatted date for recharts
  const chartData = useMemo(() => {
    return filteredData
      .map((stat) => {
        let dateLabel = stat.date;
        try {
          const d = parseISO(stat.date);
          if (isValid(d)) {
            // Show only date in YYYY-MM-DD format
            dateLabel = format(d, "yyyy-MM-dd");
          }
        } catch (e) {
          console.error("Date parsing error:", e);
        }
        return {
          ...stat,
          name: dateLabel,
        };
      })
      .sort((a, b) => {
        // Sort by the original date for chronological order
        try {
          return new Date(a.date).getTime() - new Date(b.date).getTime();
        } catch {
          return a.name.localeCompare(b.name);
        }
      });
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
        <div className="grid grid-cols-1 md:grid-cols-8 gap-4"> {/* Removed print:hidden to show filters when printing */}
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
            <DatePicker 
              isRangePicker 
              maxYear={new Date().getFullYear()} 
              date={dateRange} 
              onDateChange={(date) => {
                if ('from' in date || date === undefined) {
                  setDateRange(date as DateRange | undefined);
                }
              }}
              className="w-full"
            />
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

                // Trigger print - filters are now visible for printing
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
        <div style={{width: "100%", overflowX: "auto"}} className="mt-4 py-2 print:p-0">
          <div style={{minWidth: 0}} className="print:w-full">
            {stats ? (
              chartData.length > 0 ? (
                <ChartContainer config={config} className="w-[99%] h-[500px] print:w-full print:max-w-full">
                  <ResponsiveContainer width="100%" height={300}>
                    <Chart 
                      data={chartData} 
                      margin={{top: 20, right: 30, left: 20, bottom: 60}} 
                      className="w-full h-[300px] print:w-full print:max-w-full">
                    <XAxis 
                      dataKey="name" 
                      angle={-45} 
                      textAnchor="end" 
                      height={70} 
                      interval={0} 
                      tick={{fontSize: 10}}/>
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
                  </ResponsiveContainer>
                </ChartContainer>
              ) : (
                <div className="flex items-center justify-center h-[300px] bg-white">
                  <p className="text-muted-foreground">{t("dashboard.noData")}</p>
                </div>
              )
            ) : (
              <Skeleton className="h-[300px] w-full"/>
            )}
          </div>
        </div>

        {/* Dashboard Charts - Only Visitors (removed pie chart as requested) */}
        <div className="mt-4 grid gap-4 print:grid-cols-1">

          {/* Visitors Chart - Full width now with chart type selector */}
          <div className="md:col-span-2">
            <div className="flex justify-between items-center mb-2">
              <h2 className="font-semibold leading-none tracking-tight text-foreground">{t("dashboard.visitorCountByTime")}</h2>
            </div>
            {filteredVisitors ? (
              filteredVisitors.length > 0 ? (
                <ChartContainer config={{value: {label: t("dashboard.visitors"), color: "var(--primary)"}}} className="w-[99%] h-[500px] print:w-full print:max-w-full">
                  <ResponsiveContainer width="100%" height={300}>
                  <Chart 
                      data={filteredVisitors} 
                      margin={{top: 20, right: 30, left: 20, bottom: 60}} 
                      className="w-full h-[300px] print:w-full print:max-w-full">
                    <XAxis 
                      dataKey="name" 
                      angle={-45} 
                      textAnchor="end" 
                      height={70} 
                      interval={0} 
                      tick={{fontSize: 10}}/>
                    <YAxis/>
                    <ChartTooltip content={<ChartTooltipContent/>}/>
                    <ChartLegend content={<ChartLegendContent/>}/>
                        <Series
                          key="value"
                          type="monotone"
                          dataKey="value"
                          stroke={chartType === "line" ? "#8884d8" : undefined}
                          fill={chartType === "bar" ? "#8884d8" : undefined}
                          stackId={chartType === "bar" ? "a" : undefined}
                          strokeWidth={2}
                          dot={false}
                        />
                    </Chart>
                  </ResponsiveContainer>
                </ChartContainer>
              ) : (
                <Card className="flex items-center justify-center h-[300px] bg-white">
                  <p className="text-muted-foreground">{t("dashboard.noData")}</p>
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

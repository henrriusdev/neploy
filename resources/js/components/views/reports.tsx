import { useMemo, useState } from "react";
import { Bar, BarChart, Cell, Line, LineChart, Pie, PieChart, XAxis, YAxis } from "recharts";
import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { DateRange } from "react-day-picker";
import { DatePicker } from "@/components/forms/date-picker";
import { ChartContainer, ChartLegend, ChartLegendContent, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart";
import { format, parseISO } from "date-fns";

interface ApplicationStat {
  application_id: string;
  date: string;
  requests: number;
  errors: number;
  name?: string;
}

const metrics = [
  { key: "requests", label: "Requests", color: "#8884d8" },
  { key: "errors", label: "Errors", color: "#ff7300" },
];

export function Reports({ stats }: { stats: ApplicationStat[] }) {
  const [selectedMetrics, setSelectedMetrics] = useState<string[]>(["requests"]);
  const [dateRange, setDateRange] = useState<DateRange | undefined>(undefined);
  const [appFilter, setAppFilter] = useState<string>("all");
  const [chartType, setChartType] = useState<"line" | "bar" | "pie">("line");

  const apps = useMemo(() => Array.from(new Map(stats.map((s) => [s.application_id, { id: s.application_id, name: s.name }])).values()), [stats]);

  const filteredData = useMemo(() => {
    return stats.filter((stat) => {
      const date = new Date(stat.date);
      const inRange = (!dateRange?.from || date >= dateRange.from) && (!dateRange?.to || date <= dateRange.to);
      const appMatch = appFilter === "all" || stat.application_id === appFilter;
      return inRange && appMatch;
    });
  }, [stats, dateRange, appFilter]);

  const groupedData = useMemo(() => {
    // If the date string includes a time, group by hour (YYYY-MM-DD HH:00)
    const map = new Map<string, ApplicationStat & { hour: string }>();
    for (const stat of filteredData) {
      // Try to parse hour from stat.date
      let hour = stat.date;
      try {
        const d = parseISO(stat.date);
        hour = format(d, "yyyy-MM-dd HH:00");
      } catch {}
      const key = `${stat.application_id || "all"}-${hour}`;
      if (!map.has(key)) {
        map.set(key, { ...stat, hour, requests: 0, errors: 0 });
      }
      const agg = map.get(key)!;
      agg.requests += stat.requests;
      agg.errors += stat.errors;
    }
    // Sort by hour ascending
    return Array.from(map.values()).sort((a, b) => a.hour.localeCompare(b.hour));
  }, [filteredData]);

  const toggleMetric = (key: string) => {
    setSelectedMetrics((prev) => (prev.includes(key) ? prev.filter((k) => k !== key) : [...prev, key]));
  };

  const Chart = chartType === "bar" ? BarChart : LineChart;
  const Series = chartType === "bar" ? Bar : Line;

  const config = Object.fromEntries(metrics.map((m) => [m.key, { label: m.label, color: m.color }]));

  return (
    <Card className="pt-3">
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="space-y-1">
            <Label>Aplicación</Label>
            <Select onValueChange={setAppFilter} value={appFilter}>
              <SelectTrigger>
                <SelectValue />
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
          <div className="space-y-1">
            <Label>Rango de fechas</Label>
            <DatePicker isRangePicker maxYear={new Date().getFullYear()} date={dateRange} onDateChange={setDateRange} />
          </div>
          <div className="space-y-1">
            <Label>Tipo de gráfico</Label>
            <Select onValueChange={(value) => setChartType(value as "line" | "bar" | "pie")} value={chartType}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="line">Línea</SelectItem>
                <SelectItem value="bar">Barras</SelectItem>
                <SelectItem value="pie">Circular</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="space-x-4 pb-3 flex items-end justify-center">
            {metrics.map((m) => (
              <Label key={m.key} className="flex items-center space-x-2">
                <Checkbox checked={selectedMetrics.includes(m.key)} onCheckedChange={() => toggleMetric(m.key)} />
                <span>{m.label}</span>
              </Label>
            ))}
          </div>
        </div>
        {/* Responsive chart wrapper to prevent horizontal scroll */}
        <div style={{ width: "100%", overflowX: "auto" }}>
          <div style={{ minWidth: 600 }}>
            <ChartContainer config={config}>
              {chartType === "pie" ? (
                <PieChart>
                  <Pie
                    data={selectedMetrics.map((metric) => {
                      const m = metrics.find((m) => m.key === metric);
                      const total = groupedData.reduce((acc, item) => acc + (item[metric as keyof ApplicationStat] as number), 0);
                      return { name: m?.label, value: total, color: m?.color };
                    })}
                    dataKey="value"
                    nameKey="name"
                    cx="50%"
                    cy="50%"
                    outerRadius={120}
                    label>
                    {selectedMetrics.map((metric, i) => {
                      const m = metrics.find((m) => m.key === metric);
                      return <Cell key={`cell-${metric}`} fill={m?.color} />;
                    })}
                  </Pie>
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <ChartLegend content={<ChartLegendContent />} />
                </PieChart>
              ) : (
                <Chart data={groupedData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
                  <XAxis dataKey="hour" />
                  <YAxis />
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <ChartLegend content={<ChartLegendContent />} />
                  {selectedMetrics.map((metric) => {
                    const m = metrics.find((m) => m.key === metric);
                    return (
                      <Series
                        key={metric}
                        type="monotone"
                        dataKey={metric}
                        stroke={chartType === "line" ? m?.color : undefined}
                        fill={chartType === "bar" ? m?.color : undefined}
                        strokeWidth={2}
                        dot={false}
                      />
                    );
                  })}
                </Chart>
              )}
            </ChartContainer>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

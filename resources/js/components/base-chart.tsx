import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ResponsiveContainer } from "recharts";
import { ChartContainer, ChartLegend, ChartLegendContent, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart";
import { Bar, BarChart, CartesianGrid, Cell, Line, LineChart, Pie, PieChart, XAxis, YAxis } from "recharts";

interface BaseChartProps {
  title: string;
  data: any[];
  type: "bar" | "line" | "pie";
  dataKeys: string[];
  colors: string[];
  className?: string;
  config?: Record<string, { label: string; color: string }>;
}

export function BaseChart({ title, data, type, dataKeys, colors, className, config }: BaseChartProps) {
  const chartConfig = React.useMemo(() => {
    const baseConfig = config || {};
    
    // Ensure each dataKey has a config entry
    return dataKeys.reduce((acc, key, index) => {
      if (!acc[key]) {
        acc[key] = {
          label: key,
          color: colors[index % colors.length]
        };
      }
      return acc;
    }, {...baseConfig});
  }, [config, dataKeys, colors]);

  const renderChart = () => {
    switch (type) {
      case "bar":
        return (
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={data}>
              <XAxis dataKey="name" />
              <YAxis />
              <ChartTooltip content={<ChartTooltipContent />} />
              <ChartLegend content={<ChartLegendContent />} />
              {dataKeys.map((key, index) => (
                <Bar key={key} dataKey={key} fill={colors[index]} stackId="a" />
              ))}
            </BarChart>
          </ResponsiveContainer>
        );
      case "line":
        return (
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={data}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <ChartTooltip content={<ChartTooltipContent />} />
              <ChartLegend content={<ChartLegendContent />} />
              {dataKeys.map((key, index) => (
                <Line 
                  key={key} 
                  type="monotone" 
                  dataKey={key} 
                  stroke={colors[index].startsWith("var") ? `hsl(${colors[index]})` : colors[index]} 
                  strokeWidth={2}
                  dot={false}
                />
              ))}
            </LineChart>
          </ResponsiveContainer>
        );
      case "pie":
        return (
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <ChartTooltip content={<ChartTooltipContent />} />
              <ChartLegend content={<ChartLegendContent />} />
              <Pie 
                data={data} 
                cx="50%" 
                cy="50%" 
                labelLine={true} 
                outerRadius={80} 
                dataKey={dataKeys[0]} 
                nameKey="name"
                label
              >
                {data.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={colors[index % colors.length]} />
                ))}
              </Pie>
            </PieChart>
          </ResponsiveContainer>
        );
      default:
        return null;
    }
  };

  return (
    <Card className={`bg-card ${className}`}>
      {title && (
        <CardHeader>
          <CardTitle className="text-foreground">{title}</CardTitle>
        </CardHeader>
      )}
      <CardContent>
        <ChartContainer config={chartConfig}>
          {renderChart()}
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

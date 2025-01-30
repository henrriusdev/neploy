import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import {
  BarChart,
  LineChart,
  PieChart,
  Bar,
  Line,
  Pie,
  Cell,
  XAxis,
  YAxis,
} from "recharts";

interface BaseChartProps {
  title: string;
  data: any[];
  type: "bar" | "line" | "pie";
  dataKeys: string[];
  colors: string[];
  className?: string;
  config?: Record<string, { label: string; color: string }>;
}

export function BaseChart({
  title,
  data,
  type,
  dataKeys,
  colors,
  className,
  config,
}: BaseChartProps) {
  const renderChart = () => {
    switch (type) {
      case "bar":
        return (
          <BarChart data={data}>
            <XAxis dataKey="name" />
            <YAxis />
            {dataKeys.map((key, index) => (
              <Bar
                key={key}
                dataKey={key}
                fill={colors[index]}
                stackId="a"
                isAnimationActive={false}
              />
            ))}
          </BarChart>
        );
      case "line":
        return (
          <LineChart data={data}>
            <XAxis dataKey="name" />
            <YAxis />
            {dataKeys.map((key, index) => (
              <Line
                key={key}
                type="monotone"
                dataKey={key}
                stroke={colors[index]}
                strokeWidth={2}
              />
            ))}
          </LineChart>
        );
      case "pie":
        return (
          <PieChart>
            <Pie
              data={data}
              cx="50%"
              cy="50%"
              labelLine={true}
              outerRadius={80}
              dataKey={dataKeys[0]}
              nameKey="name">
              {data.map((entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={colors[index % colors.length]}
                />
              ))}
            </Pie>
          </PieChart>
        );
      default:
        return null;
    }
  };

  return (
    <Card className={`bg-card ${className}`}>
      <CardHeader>
        <CardTitle className="text-foreground">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <ChartContainer config={config} className="h-[350px] w-full">
          <>
            {renderChart()}
            <ChartTooltip content={<ChartTooltipContent />} />
          </>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

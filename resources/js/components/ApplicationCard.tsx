import { Application } from "@/types/common";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { TechIcon } from "@/components/TechIcon";
import { Play, Square, Trash2 } from "lucide-react";

interface ApplicationCardProps {
  app: Application;
  onStart: (id: string) => void;
  onStop: (id: string) => void;
  onDelete: (id: string) => void;
}

export function ApplicationCard({
  app,
  onStart,
  onStop,
  onDelete,
}: ApplicationCardProps) {
  const getStatusBadgeColor = (status: Application["status"]) => {
    switch (status) {
      case "Running":
        return "bg-green-500";
      case "Building":
        return "bg-yellow-500";
      case "Stopped":
        return "bg-gray-500";
      case "Error":
        return "bg-red-500";
      default:
        return "bg-gray-500";
    }
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-start justify-between space-y-0">
        <div>
          <CardTitle className="text-xl">{app.appName}</CardTitle>
          <CardDescription>
            {app?.techStack === null ? (
              "Auto detected"
            ) : (
              <TechIcon name={app.techStack.name} />
            )}
          </CardDescription>
        </div>
        <Badge className={`${getStatusBadgeColor(app.status)} text-white`}>
          {app.status}
        </Badge>
      </CardHeader>
      <CardContent>
        <div className="flex gap-2">
          {app.status !== "Running" && (
            <Button size="sm" variant="outline" onClick={() => onStart(app.id)}>
              <Play className="h-4 w-4 mr-1" />
              Start
            </Button>
          )}
          {app.status === "Running" && (
            <Button size="sm" variant="outline" onClick={() => onStop(app.id)}>
              <Square className="h-4 w-4 mr-1" />
              Stop
            </Button>
          )}
          <Button
            size="sm"
            variant="destructive"
            onClick={() => onDelete(app.id)}>
            <Trash2 className="h-4 w-4 mr-1" />
            Delete
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}

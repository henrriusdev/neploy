import {Application} from "@/types/common";
import {Badge} from "@/components/ui/badge";
import {Button} from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {TechIcon} from "@/components/icons/tech-icon";
import {Play, Square, Trash2} from "lucide-react";
import {useTranslation} from "react-i18next";
import "@/i18n";
import {router} from "@inertiajs/react";

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
  const {t} = useTranslation();

  const getStatusBadgeColor = (status: Application["status"]) => {
    switch (status) {
      case "Running":
        return "!bg-green-500";
      case "Building":
        return "!bg-yellow-500";
      case "Stopped":
        return "!bg-gray-500";
      case "Error":
        return "!bg-red-500";
      case "Created":
        return "!bg-blue-500 !text-white";
      default:
        return "!bg-gray-500";
    }
  };

  const translateStatus = (status: Application["status"]) => {
    switch (status) {
      case "Running":
        return t("dashboard.applications.status.running");
      case "Building":
        return t("dashboard.applications.status.building");
      case "Stopped":
        return t("dashboard.applications.status.stopped");
      case "Error":
        return t("dashboard.applications.status.error");
      case "Created":
        return t("dashboard.applications.status.created");
      default:
        return t("dashboard.applications.status.unknown");
    }
  };

  return (
    <Card className={"hover:!bg-primary hover:cursor-pointer transition-colors"}
          onClick={() => router.visit(`/dashboard/applications/${app.id}`)}>
      <CardHeader className="flex flex-row items-start justify-between space-y-0">
        <div>
          <CardTitle className="text-xl">{app.appName}</CardTitle>
          <CardDescription>
            {app?.techStack === null ? (
              "Auto detected"
            ) : (
              <TechIcon name={app.techStack.name}/>
            )}
          </CardDescription>
        </div>
        <Badge className={`${getStatusBadgeColor(app.status)}`} >
          {translateStatus(app.status)}
        </Badge>
      </CardHeader>
      <CardContent>
        <div className="flex gap-2">
          {app.status !== "Running" && (
            <Button size="sm" variant="outline" onClick={() => onStart(app.id)}>
              <Play className="h-4 w-4 mr-1"/>
              {t("dashboard.applications.start")}
            </Button>
          )}
          {app.status === "Running" && (
            <Button size="sm" variant="outline" onClick={() => onStop(app.id)}>
              <Square className="h-4 w-4 mr-1"/>
              {t("dashboard.applications.stop")}
            </Button>
          )}
          <Button
            size="sm"
            variant="destructive"
            onClick={() => onDelete(app.id)}>
            <Trash2 className="h-4 w-4 mr-1"/>
            {t("dashboard.applications.delete")}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}

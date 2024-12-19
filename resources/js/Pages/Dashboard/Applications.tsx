import * as React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { PlusCircle, Grid, List, Upload, Settings } from "lucide-react";
import DashboardLayout from "@/components/Layouts/DashboardLayout";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

interface ApplicationStat {
  id: string;
  applicationId: string;
  environmentId: string;
  date: string;
  requests: number;
  errors: number;
  averageResponseTime: number;
  dataTransfered: number;
  uniqueVisitors: number;
  healthy: boolean;
  createdAt: string;
  updatedAt: string;
}

interface Application {
  id: string;
  appName: string;
  storageLocation: string;
  deployLocation: string;
  techStackId: string;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
  stats: ApplicationStat[];
}

interface ApplicationsProps {
  user?: {
    name: string;
    email: string;
  };
  teamName: string;
  logoUrl: string;
  applications?: Application[] | null;
}

function Applications({
  user,
  teamName,
  logoUrl,
  applications = null,
}: ApplicationsProps) {
  const [viewMode, setViewMode] = React.useState<"grid" | "list">("grid");

  // Calculate stats from all applications
  const aggregateStats = React.useMemo(() => {
    if (!applications) return { total: 0, healthy: 0, deployed: 0 };
    
    return {
      total: applications.length,
      healthy: applications.filter(app => 
        app.stats.some(stat => stat.healthy)
      ).length,
      deployed: applications.filter(app => 
        app.deployLocation && app.stats.length > 0
      ).length,
    };
  }, [applications]);

  return (
    <div className="space-y-6 p-3">
      {/* Stats Section */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Apps</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{aggregateStats.total}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Healthy Apps</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{aggregateStats.healthy}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Deployed Apps</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{aggregateStats.deployed}</div>
          </CardContent>
        </Card>
      </div>

      {/* Actions Bar */}
      <div className="flex justify-between items-center">
        <div className="flex gap-2">
          <Button
            variant={viewMode === "grid" ? "default" : "outline"}
            size="icon"
            onClick={() => setViewMode("grid")}
          >
            <Grid className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === "list" ? "default" : "outline"}
            size="icon"
            onClick={() => setViewMode("list")}
          >
            <List className="h-4 w-4" />
          </Button>
        </div>
        <div className="flex gap-2">
          <Dialog>
            <DialogTrigger asChild>
              <Button>
                <Upload className="mr-2 h-4 w-4" />
                Upload App
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Upload New Application</DialogTitle>
                <DialogDescription>
                  Upload your application package or connect your repository.
                </DialogDescription>
              </DialogHeader>
              {/* Add upload form here */}
            </DialogContent>
          </Dialog>
          <Button variant="outline">
            <Settings className="mr-2 h-4 w-4" />
            Configure
          </Button>
        </div>
      </div>

      {/* Applications List/Grid */}
      {applications === null || applications.length === 0 ? (
        <Card className="p-12">
          <div className="flex flex-col items-center justify-center text-center space-y-4">
            <div className="p-3 bg-primary/10 rounded-full">
              <PlusCircle className="w-8 h-8 text-primary" />
            </div>
            <div>
              <h3 className="text-lg font-semibold">No applications found</h3>
              <p className="text-sm text-muted-foreground">
                Get started by uploading your first application or connecting your repository.
              </p>
            </div>
            <Dialog>
              <DialogTrigger asChild>
                <Button>
                  <Upload className="mr-2 h-4 w-4" />
                  Upload App
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Upload New Application</DialogTitle>
                  <DialogDescription>
                    Upload your application package or connect your repository.
                  </DialogDescription>
                </DialogHeader>
                {/* Add upload form here */}
              </DialogContent>
            </Dialog>
          </div>
        </Card>
      ) : (
        <div className={viewMode === "grid" ? "grid grid-cols-3 gap-4" : "space-y-4"}>
          {applications.map((app) => {
            const latestStat = app.stats[app.stats.length - 1];
            return (
              <Card key={app.id}>
                <CardHeader>
                  <CardTitle>{app.appName}</CardTitle>
                  <CardDescription>{app.deployLocation}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {latestStat && (
                      <div className="grid grid-cols-2 gap-2">
                        <div className="text-sm">
                          <span className="text-muted-foreground">Status: </span>
                          <Badge variant={latestStat.healthy ? "secondary" : "destructive"}>
                            {latestStat.healthy ? "Healthy" : "Unhealthy"}
                          </Badge>
                        </div>
                        <div className="text-sm">
                          <span className="text-muted-foreground">Requests: </span>
                          {latestStat.requests}
                        </div>
                        <div className="text-sm">
                          <span className="text-muted-foreground">Errors: </span>
                          {latestStat.errors}
                        </div>
                        <div className="text-sm">
                          <span className="text-muted-foreground">Response Time: </span>
                          {latestStat.averageResponseTime}ms
                        </div>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}

Applications.layout = (page: any) => {
  const { user: pageUser, teamName, logoUrl } = page.props;
  const user = {
    name: pageUser.name,
    email: pageUser.email,
    avatar:
      pageUser.provider === "github"
        ? "https://unavatar.io/github/" + pageUser.username
        : "https://unavatar.io/" + pageUser.email,
  };
  return (
    <DashboardLayout teamName={teamName} logoUrl={logoUrl} user={user}>
      {page}
    </DashboardLayout>
  );
};

export default Applications;

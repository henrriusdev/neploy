import { DynamicForm } from "@/components/DynamicForm";
import DashboardLayout from "@/components/Layouts/DashboardLayout";
import { ApplicationCard } from "@/components/ApplicationCard";
import { ApplicationForm } from "@/components/ApplicationForm";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useToast } from "@/hooks/use-toast";
import { useWebSocket } from "@/hooks/useWebSocket";
import { Application } from "@/types/common";
import type { ActionMessage, ActionResponse, Input as InputType, ProgressMessage } from "@/types/websocket";
import axios from "axios";
import { debounce } from "lodash";
import { Grid, List, PlusCircle } from "lucide-react";
import * as React from "react";
import { useEffect, useMemo, useState } from "react";
import * as z from "zod";
import { ApplicationsProps } from "@/types/props";


const uploadFormSchema = z.object({
  appName: z.string().min(1, "Application name is required"),
  description: z.string().optional(),
  repoUrl: z
    .string()
    .refine(
      (value) => {
        if (!value) return true; // Allow empty string
        try {
          const url = new URL(value);
          if (!["github.com", "gitlab.com"].includes(url.hostname)) {
            return false;
          }
          const parts = url.pathname.split("/").filter(Boolean);
          return parts.length === 2; // Should have exactly user and repo
        } catch {
          return false;
        }
      },
      { message: "Must be a valid GitHub or GitLab repository URL" }
    )
    .optional(),
  branch: z.string().optional(),
});

function Applications({
  user,
  teamName,
  logoUrl,
  applications: initialApplications = null,
}: ApplicationsProps) {
  const [applications, setApplications] = useState<Application[] | null>(
    initialApplications
  );
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [isUploading, setIsUploading] = useState(false);
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);
  const [isLoadingBranches, setIsLoadingBranches] = useState(false);
  const [branches, setBranches] = useState<string[]>([]);

  const [actionDialog, setActionDialog] = useState<{
    show: boolean;
    title: string;
    description: string;
    fields: InputType[];
    onSubmit: (data: any) => void;
  }>({
    show: false,
    title: "",
    description: "",
    fields: [],
    onSubmit: () => {},
  });
  const { toast } = useToast();
  const { onNotification, onInteractive, sendMessage } = useWebSocket();

  useEffect(() => {
    const unsubProgress = onNotification((message: ProgressMessage) => {
      if (message.type === "progress") {
        toast({
          title: "Deployment Progress",
          description: message.message,
        });
      }
    });

    const unsubInteractive = onInteractive((message: ActionMessage) => {
      console.log("Received interactive message:", message);
      if (!message?.inputs || !Array.isArray(message.inputs)) {
        console.error("Invalid message inputs:", message.inputs);
        return;
      }

      setActionDialog({
        show: true,
        title: message.title || "Action Required",
        description: message.message || "",
        fields: message.inputs.map((input) => ({
          ...input,
          // Add validation for port number
          validate:
            input.name === "port"
              ? (value: string) => {
                  const port = parseInt(value);
                  if (isNaN(port) || port < 1 || port > 65535) {
                    return "Please enter a valid port number (1-65535)";
                  }
                  return true;
                }
              : undefined,
        })),
        onSubmit: (data) => {
          console.log("Submitting form data:", data);
          const response: ActionResponse = {
            type: message.type,
            action: message.action,
            data: {
              ...data,
              action: message.action,
            },
          };
          console.log("Sending response:", response);
          sendMessage(response.type, response.action, response.data);
          setActionDialog((prev) => ({ ...prev, show: false }));

          // Show confirmation toast
          toast({
            title: "Port Configuration",
            description: `Port ${data.port} will be exposed for this application.`,
          });
        },
      });
    });

    // Store unsubscribe functions
    const unsubFunctions = [unsubProgress, unsubInteractive];

    return () => {
      // Call all unsubscribe functions
      unsubFunctions.forEach((unsub) => unsub && unsub());
    };
  }, [onNotification, onInteractive, sendMessage, toast]);

  const refreshApplications = async () => {
    try {
      const { data } = await axios.get("/applications");
      setApplications(data);
    } catch (error) {
      toast({
        title: "Error",
        description: axios.isAxiosError(error)
          ? error.response?.data?.message || "Failed to fetch applications"
          : "Failed to fetch applications",
        variant: "destructive",
      });
    }
  };

  const fetchBranches = async (repoUrl: string) => {
    if (!repoUrl) {
      setBranches([]);
      return;
    }

    try {
      setIsLoadingBranches(true);
      const { data } = await axios.post("/applications/branches", {
        repoUrl: repoUrl,
      });
      setBranches(data.branches);
    } catch (error) {
      console.error("Error fetching branches:", error);
      toast({
        title: "Error",
        description: axios.isAxiosError(error)
          ? error.response?.data?.message ||
            "Failed to fetch repository branches"
          : "Failed to fetch repository branches",
        variant: "destructive",
      });
    } finally {
      setIsLoadingBranches(false);
    }
  };

  const debouncedFetchBranches = useMemo(
    () => debounce(fetchBranches, 1000),
    []
  );

  const onSubmit = async (
    values: z.infer<typeof uploadFormSchema>,
    file: File | null
  ) => {
    setIsUploading(true);

    try {
      if (values.repoUrl && !values.branch) {
        toast({
          title: "Error",
          description: "Please select a branch",
          variant: "destructive",
        });
        return;
      }

      const {
        data: { id: applicationId },
      } = await axios.post("/applications", {
        appName: values.appName,
        description:
          values.description ||
          (values.repoUrl
            ? `Deployed from GitHub: ${values.repoUrl}`
            : "Uploaded from ZIP file"),
      });

      if (values.repoUrl) {
        if (file) {
          toast({
            title: "Error",
            description:
              "Please provide either a GitHub URL or a ZIP file, not both",
            variant: "destructive",
          });
          return;
        }

        await axios.post(`/applications/${applicationId}/deploy`, {
          repoUrl: values.repoUrl,
          branch: values.branch,
        });

        toast({
          title: "Success",
          description: "GitHub repository deployment started",
        });
      } else if (file) {
        const formData = new FormData();
        formData.append("file", file);

        await axios.post(`/applications/${applicationId}/upload`, formData, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        });

        toast({
          title: "Success",
          description: "Application file uploaded successfully",
        });
      } else {
        toast({
          title: "Error",
          description: "Please provide either a GitHub URL or a ZIP file",
          variant: "destructive",
        });
        return;
      }

      await refreshApplications();
      setUploadDialogOpen(false);
    } catch (error) {
      toast({
        title: "Error",
        description: axios.isAxiosError(error)
          ? error.response?.data?.message || error.message
          : "An error occurred",
        variant: "destructive",
      });
    } finally {
      setIsUploading(false);
    }
  };

  const handleApplicationAction = async (
    appId: string,
    action: "start" | "stop" | "delete"
  ) => {
    try {
      if (action === "delete") {
        await axios.delete(`/applications/${appId}`);
      } else {
        await axios.post(`/applications/${appId}/${action}`);
      }

      toast({
        title: "Success",
        description: `Application ${action} request sent`,
      });

      await refreshApplications();
    } catch (error) {
      toast({
        title: "Error",
        description: axios.isAxiosError(error)
          ? error.response?.data?.message || `Failed to ${action} application`
          : `Failed to ${action} application`,
        variant: "destructive",
      });
    }
  };

  React.useEffect(() => {
    refreshApplications();
  }, []);

  React.useEffect(() => {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host}/ws/interactive`;
    const ws = new WebSocket(wsUrl);

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === "APPLICATION_UPDATE") {
        setApplications((prev) => {
          if (!prev) return prev;
          return prev.map((app) =>
            app.id === data.applicationId ? { ...app, ...data.updates } : app
          );
        });
      }
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <div className="space-y-6 p-3">
      {/* Stats Section */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Apps</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {applications?.length || 0}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Running Apps</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {applications?.filter((app) => app.status === "Running").length ||
                0}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Failed Apps</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {applications?.filter((app) => app.status === "Error").length ||
                0}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Actions Bar */}
      <div className="flex justify-between items-center">
        <Dialog open={uploadDialogOpen} onOpenChange={setUploadDialogOpen}>
          <DialogTrigger asChild>
            <Button className="flex items-center gap-2">
              <PlusCircle className="h-4 w-4" />
              New Application
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Deploy New Application</DialogTitle>
              <DialogDescription>
                Upload a zip file or provide a GitHub repository URL to deploy
                your application.
              </DialogDescription>
            </DialogHeader>
            <ApplicationForm
              onSubmit={onSubmit}
              isUploading={isUploading}
              branches={branches}
              isLoadingBranches={isLoadingBranches}
              onRepoUrlChange={debouncedFetchBranches}
            />
          </DialogContent>
        </Dialog>
        <div className="flex items-center gap-2">
          <Button
            variant={viewMode === "grid" ? "default" : "outline"}
            size="icon"
            onClick={() => setViewMode("grid")}>
            <Grid className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === "list" ? "default" : "outline"}
            size="icon"
            onClick={() => setViewMode("list")}>
            <List className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Applications List/Grid */}
      {!applications || applications.length === 0 ? (
        <Card className="p-12">
          <div className="flex flex-col items-center justify-center text-center space-y-4">
            <div className="p-3 bg-primary/10 rounded-full">
              <PlusCircle className="w-8 h-8 text-primary" />
            </div>
            <div>
              <h3 className="text-lg font-semibold">No applications found</h3>
              <p className="text-sm text-muted-foreground">
                Get started by clicking the "New Application" button above to
                deploy your first application.
              </p>
            </div>
          </div>
        </Card>
      ) : (
        <div
          className={
            viewMode === "grid"
              ? "grid gap-4 md:grid-cols-2 lg:grid-cols-3"
              : "space-y-4"
          }>
          {applications.map((app: Application) => (
            <ApplicationCard
              key={app.id}
              app={app}
              onStart={(id) => handleApplicationAction(id, "start")}
              onStop={(id) => handleApplicationAction(id, "stop")}
              onDelete={(id) => handleApplicationAction(id, "delete")}
            />
          ))}
        </div>
      )}
      <Dialog
        open={actionDialog.show}
        onOpenChange={(open) =>
          setActionDialog((prev) => ({ ...prev, show: open }))
        }>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{actionDialog.title}</DialogTitle>
            <DialogDescription>{actionDialog.description}</DialogDescription>
          </DialogHeader>
          {actionDialog.show &&
            actionDialog.fields &&
            actionDialog.fields.length > 0 && (
              <DynamicForm
                fields={actionDialog.fields}
                onSubmit={actionDialog.onSubmit}
                className="mt-4"
              />
            )}
        </DialogContent>
      </Dialog>
    </div>
  );
}

Applications.layout = (page: any) => {
  return <DashboardLayout>{page}</DashboardLayout>;
};

export default Applications;
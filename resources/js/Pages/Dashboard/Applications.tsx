import { DynamicForm } from "@/components/DynamicForm";
import DashboardLayout from "@/components/Layouts/DashboardLayout";
import { TechIcon } from "@/components/TechIcon";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/hooks/use-toast";
import { useWebSocket } from "@/hooks/useWebSocket";
import type { ActionMessage, Input as InputType, ProgressMessage } from "@/types/websocket";
import { zodResolver } from "@hookform/resolvers/zod";
import axios from "axios";
import {
  Grid,
  List,
  Play,
  PlusCircle,
  Square,
  Trash2
} from "lucide-react";
import * as React from "react";
import { useEffect, useState } from "react";
import { useDropzone } from "react-dropzone";
import { useForm } from "react-hook-form";
import * as z from "zod";

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

interface TechStack {
  id: string;
  name: string;
  description: string;
}

interface Application {
  id: string;
  appName: string;
  storageLocation: string;
  deployLocation: string;
  techStackId: string;
  status: "Building" | "Running" | "Stopped" | "Error";
  language?: string;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
  stats: ApplicationStat[];
  techStack: TechStack;
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

const uploadFormSchema = z.object({
  appName: z.string().min(1, "Application name is required"),
  description: z.string().optional(),
  repoUrl: z
    .string()
    .refine((value) => value === "" || (value && typeof value === "string" && Boolean(new URL(value))),
      { message: "Invalid URL" })
    .optional(),
});

function Applications({
  user,
  teamName,
  logoUrl,
  applications: initialApplications = null,
}: ApplicationsProps) {
  const [applications, setApplications] = useState<Application[] | null>(initialApplications);
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [isUploading, setIsUploading] = useState(false);
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);
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
    onSubmit: () => { }
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
      setActionDialog({
        show: true,
        title: message.title,
        description: message.message,
        fields: message.inputs,
        onSubmit: (data) => {
          // Format response to match backend expectations
          const response = {
            action: data.action,
            data: {
              port: data.port
            }
          };
          // Send response back through websocket
          sendMessage(message.type, response.action, response.data);
          setActionDialog((prev) => ({ ...prev, show: false }));
        },
      });
    });

    return () => {
      unsubProgress();
      unsubInteractive();
    };
  }, [onNotification, onInteractive, sendMessage, toast]);

  const form = useForm<z.infer<typeof uploadFormSchema>>({
    resolver: zodResolver(uploadFormSchema),
    defaultValues: {
      appName: "",
      description: "",
      repoUrl: "",
    },
  });

  const onDrop = React.useCallback((acceptedFiles: File[]) => {
    const file = acceptedFiles[0];
    if (!file) return;

    if (!file.name.endsWith(".zip")) {
      toast({
        title: "Invalid file type",
        description: "Please upload a .zip file",
        variant: "destructive",
      });
      return;
    }

    setUploadedFile(file);
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      "application/zip": [".zip"],
    },
    maxFiles: 1,
    multiple: undefined,
    onDragEnter: undefined,
    onDragOver: undefined,
    onDragLeave: undefined,
  });

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

  const onSubmit = async (values: z.infer<typeof uploadFormSchema>) => {
    setIsUploading(true);

    try {
      // Create application
      const { data: { id: applicationId } } = await axios.post("/applications", {
        appName: values.appName,
        description:
          values.description ||
          (values.repoUrl
            ? `Deployed from GitHub: ${values.repoUrl}`
            : "Uploaded from ZIP file"),
      });

      // Deploy either from GitHub URL or file upload, not both
      if (values.repoUrl) {
        if (uploadedFile) {
          toast({
            title: "Error",
            description: "Please provide either a GitHub URL or a ZIP file, not both",
            variant: "destructive",
          });
          return;
        }

        await axios.post(`/applications/${applicationId}/deploy`, {
          repoUrl: values.repoUrl,
        });

        toast({
          title: "Success",
          description: "GitHub repository deployment started",
        });
      } else if (uploadedFile) {
        const formData = new FormData();
        formData.append("file", uploadedFile);

        await axios.post(`/applications/${applicationId}/upload`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
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

      // Refresh the applications list
      await refreshApplications();

      // Reset form and close dialog
      form.reset();
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

      // Refresh the applications list
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

  // Initial load of applications
  React.useEffect(() => {
    refreshApplications();
  }, []);

  // WebSocket connection for real-time updates
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
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-4"
              >
                <FormField
                  control={form.control}
                  name="appName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Application Name</FormLabel>
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Description</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Enter application description"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="repoUrl"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>GitHub Repository URL (Optional)</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          placeholder="https://github.com/username/repo"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div
                  {...getRootProps()}
                  className="border-2 border-dashed rounded-lg p-6 text-center cursor-pointer hover:border-primary"
                >
                  <input {...getInputProps()} />
                  {isDragActive ? (
                    <p>Drop the ZIP file here...</p>
                  ) : (
                    <p>Drag & drop a ZIP file here, or click to select</p>
                  )}
                </div>
                <Button type="submit" className="w-full" disabled={isUploading}>
                  {isUploading ? "Deploying..." : "Deploy Application"}
                </Button>
              </form>
            </Form>
          </DialogContent>
        </Dialog>
        <div className="flex items-center gap-2">
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
          }
        >
          {applications.map((app: Application) => (
            <Card key={app.id}>
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
                <Badge
                  className={`${getStatusBadgeColor(app.status)} text-white`}
                >
                  {app.status}
                </Badge>
              </CardHeader>
              <CardContent>
                <div className="flex gap-2">
                  {app.status !== "Running" && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleApplicationAction(app.id, "start")}
                    >
                      <Play className="h-4 w-4 mr-1" />
                      Start
                    </Button>
                  )}
                  {app.status === "Running" && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleApplicationAction(app.id, "stop")}
                    >
                      <Square className="h-4 w-4 mr-1" />
                      Stop
                    </Button>
                  )}
                  <Button
                    size="sm"
                    variant="destructive"
                    onClick={() => handleApplicationAction(app.id, "delete")}
                  >
                    <Trash2 className="h-4 w-4 mr-1" />
                    Delete
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
      <Dialog open={actionDialog.show} onOpenChange={(open) => setActionDialog(prev => ({ ...prev, show: open }))}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{actionDialog.title}</DialogTitle>
            <DialogDescription>{actionDialog.description}</DialogDescription>
          </DialogHeader>
          <DynamicForm
            fields={actionDialog.fields}
            onSubmit={actionDialog.onSubmit}
            className="mt-4"
          />
        </DialogContent>
      </Dialog>
    </div>
  );
}

Applications.layout = (page: any) => {
  return (
    <DashboardLayout>
      {page}
    </DashboardLayout>
  );
};

export default Applications;

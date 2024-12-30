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
  FormDescription
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
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
import { useEffect, useState, useMemo } from "react";
import { useDropzone } from "react-dropzone";
import { useForm } from "react-hook-form";
import * as z from "zod";
import { debounce, DebouncedFunc } from 'lodash';

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
    .refine(
      (value) => {
        if (!value) return true; // Allow empty string
        try {
          const url = new URL(value);
          // Check if it's GitHub or GitLab
          if (!['github.com', 'gitlab.com'].includes(url.hostname)) {
            return false;
          }
          // Check if it has the pattern: hostname/user/repo
          const parts = url.pathname.split('/').filter(Boolean);
          return parts.length === 2; // Should have exactly user and repo
        } catch {
          return false;
        }
      },
      { message: "Must be a valid GitHub or GitLab repository URL (e.g., https://github.com/user/repo)" }
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
  const [applications, setApplications] = useState<Application[] | null>(initialApplications);
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);
  const [isLoadingBranches, setIsLoadingBranches] = useState(false);
  const [branches, setBranches] = useState<string[]>([]);
  const [selectedBranch, setSelectedBranch] = useState<string>("");

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
      branch: "",
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

  const fetchBranches = async (repoUrl: string) => {
    if (!repoUrl) {
      setBranches([]);
      setSelectedBranch('');
      form.setValue('branch', '');
      return;
    }

    try {
      // Validate URL format before making API call
      try {
        const url = new URL(repoUrl);
        if (!['github.com', 'gitlab.com'].includes(url.hostname)) {
          return;
        }
        const parts = url.pathname.split('/').filter(Boolean);
        if (parts.length !== 2) {
          return;
        }
      } catch {
        return;
      }

      setIsLoadingBranches(true);
      const { data } = await axios.post('/applications/branches', {
        repoUrl: repoUrl
      });
      setBranches(data.branches);
      
      // Set default branch if available (usually 'main' or 'master')
      const defaultBranch = data.branches.find((b: string) => 
        ['main', 'master'].includes(b)
      ) || data.branches[0];
      
      if (defaultBranch) {
        setSelectedBranch(defaultBranch);
        form.setValue('branch', defaultBranch);
      }
    } catch (error) {
      console.error('Error fetching branches:', error);
      toast({
        title: "Error",
        description: axios.isAxiosError(error)
          ? error.response?.data?.message || "Failed to fetch repository branches"
          : "Failed to fetch repository branches",
        variant: "destructive",
      });
    } finally {
      setIsLoadingBranches(false);
    }
  };

  const debouncedFetchBranches = useMemo(
    () => debounce(fetchBranches, 1000),
    [] // Empty deps since we want to create this only once
  );

  useEffect(() => {
    const subscription = form.watch((value, { name }) => {
      if (name === 'repoUrl') {
        debouncedFetchBranches(value.repoUrl);
      }
    });
    
    return () => {
      subscription.unsubscribe();
      debouncedFetchBranches.cancel();
    };
  }, [debouncedFetchBranches, form]);

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
      if (values.repoUrl && !values.branch) {
        toast({
          title: "Error",
          description: "Please select a branch",
          variant: "destructive",
        });
        return;
      }

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
          branch: values.branch,
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
                      <FormLabel>GitHub/GitLab Repository URL (Optional)</FormLabel>
                      <FormControl>
                        <Input 
                          placeholder="https://github.com/username/repository" 
                          {...field} 
                        />
                      </FormControl>
                      <FormDescription>
                        Enter a valid GitHub or GitLab repository URL (e.g., https://github.com/user/repo)
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                {form.watch('repoUrl') && (
                  <FormField
                    control={form.control}
                    name="branch"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Branch</FormLabel>
                        <Select
                          disabled={isLoadingBranches}
                          value={field.value}
                          onValueChange={(value) => {
                            field.onChange(value);
                            setSelectedBranch(value);
                          }}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue 
                                placeholder={
                                  isLoadingBranches 
                                    ? "Loading branches..." 
                                    : "Select a branch"
                                } 
                              />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {branches.map((branch) => (
                              <SelectItem key={branch} value={branch}>
                                {branch}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                )}
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

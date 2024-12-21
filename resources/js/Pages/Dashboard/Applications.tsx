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
import {
  PlusCircle,
  Grid,
  List,
  Upload,
  Settings,
  Play,
  Square,
  Trash2,
  AlertCircle,
} from "lucide-react";
import DashboardLayout from "@/components/Layouts/DashboardLayout";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useDropzone } from "react-dropzone";
import { useToast } from "@/hooks/use-toast";
import axios from "axios";
import { Textarea } from "@/components/ui/textarea";
import { TechIcon } from "@/components/TechIcon";

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
  language: z.string().optional(),
});

const SUPPORTED_LANGUAGES = ["Node.js", "Go", "Python", "Java"];

function Applications({
  user,
  teamName,
  logoUrl,
  applications: initialApplications = null,
}: ApplicationsProps) {
  const [viewMode, setViewMode] = React.useState<"grid" | "list">("grid");
  const [applications, setApplications] = React.useState<Application[] | null>(
    initialApplications
  );
  const [isUploading, setIsUploading] = React.useState(false);
  const [uploadDialogOpen, setUploadDialogOpen] = React.useState(false);
  const [uploadedFile, setUploadedFile] = React.useState<File | null>(null);
  const { toast } = useToast();

  const form = useForm<z.infer<typeof uploadFormSchema>>({
    resolver: zodResolver(uploadFormSchema),
    defaultValues: {
      appName: "",
      description: "",
      repoUrl: "",
      language: undefined,
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

  const onSubmit = async (values: z.infer<typeof uploadFormSchema>) => {
    setIsUploading(true);
    try {
      // First, create the application record
      const createResponse = await axios.post("/applications", {
        appName: values.appName,
        description:
          values.description ||
          (values.repoUrl
            ? `Deployed from GitHub: ${values.repoUrl}`
            : "Uploaded from ZIP file"),
        techStack: values.language || "auto-detect",
      });

      // If application creation fails (remember that is an axios call), throw an error
      if (createResponse.status >= 400) {
        throw new Error(
          "Failed to create application" + createResponse.statusText
        );
      }

      // Extract the application ID from the response
      const applicationId = createResponse.data.id;

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

        const { data: deployData } = await axios.post(
          `/applications/${applicationId}/deploy`,
          {
            repoUrl: values.repoUrl,
          }
        );

        toast({
          title: "Success",
          description: "GitHub repository deployment started",
        });
      } else if (uploadedFile) {
        const formData = new FormData();
        formData.append("file", uploadedFile);

        const { data: uploadData } = await axios.post(
          `/applications/${applicationId}/upload`,
          formData,
          {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          }
        );

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
      const { data: updatedApplications } = await axios.get("/applications");
      setApplications(updatedApplications);

      // Reset form and close dialog
      form.reset();
      setUploadDialogOpen(false);
    } catch (error) {
      toast({
        title: "Error",
        description:
          error instanceof Error
            ? error.message
            : "Failed to deploy application",
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
      const response = await fetch(`/api/applications/${appId}/${action}`, {
        method: "POST",
      });

      if (!response.ok) {
        throw new Error(`Failed to ${action} application`);
      }

      toast({
        title: "Success",
        description: `Application ${action} request sent`,
      });

      // For delete action, remove from local state
      if (action === "delete") {
        setApplications((prev) =>
          prev ? prev.filter((app) => app.id !== appId) : null
        );
      }
    } catch (error) {
      toast({
        title: "Error",
        description:
          error instanceof Error
            ? error.message
            : `Failed to ${action} application`,
        variant: "destructive",
      });
    }
  };

  // WebSocket connection for real-time updates
  React.useEffect(() => {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host}/ws`;
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
                className="space-y-4">
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
                  className="border-2 border-dashed rounded-lg p-6 text-center cursor-pointer hover:border-primary">
                  <input {...getInputProps()} />
                  {isDragActive ? (
                    <p>Drop the ZIP file here...</p>
                  ) : (
                    <p>Drag & drop a ZIP file here, or click to select</p>
                  )}
                </div>
                <FormField
                  control={form.control}
                  name="language"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Programming Language</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Select language" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {SUPPORTED_LANGUAGES.map((lang) => (
                            <SelectItem key={lang} value={lang}>
                              {lang}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
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
                  className={`${getStatusBadgeColor(app.status)} text-white`}>
                  {app.status}
                </Badge>
              </CardHeader>
              <CardContent>
                <div className="flex gap-2">
                  {app.status !== "Running" && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleApplicationAction(app.id, "start")}>
                      <Play className="h-4 w-4 mr-1" />
                      Start
                    </Button>
                  )}
                  {app.status === "Running" && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleApplicationAction(app.id, "stop")}>
                      <Square className="h-4 w-4 mr-1" />
                      Stop
                    </Button>
                  )}
                  <Button
                    size="sm"
                    variant="destructive"
                    onClick={() => handleApplicationAction(app.id, "delete")}>
                    <Trash2 className="h-4 w-4 mr-1" />
                    Delete
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
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

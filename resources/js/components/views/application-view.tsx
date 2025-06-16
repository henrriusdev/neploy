"use client";

import { ApplicationForm, DynamicForm } from "@/components/forms";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Progress } from "@/components/ui/progress";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useToast, useWebSocket } from "@/hooks";
import { sanitizeAppName } from "@/lib/utils";
import {
  useDeleteVersionMutation,
  useDeployApplicationMutation,
  useLoadBranchesQuery,
  useStartApplicationMutation,
  useStopApplicationMutation,
  useUploadApplicationMutation,
} from "@/services/api/applications";
import { ActionMessage, ActionResponse, ApplicationProps, ProgressMessage } from "@/types";
import type { Input as InputInterface } from "@/types/websocket";
import { router } from "@inertiajs/react";
import { DialogTrigger } from "@radix-ui/react-dialog";
import { CirclePlay, Pause, Plus, Trash2 } from "lucide-react";
import { FC, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";

const uploadFormSchema = z.object({
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
      { message: "Must be a valid GitHub or GitLab repository URL" },
    )
    .optional(),
  branch: z.string().optional(),
});

export const ApplicationView: FC<ApplicationProps> = ({ application }) => {
  const [versions, setVersions] = useState(application.versions);
  const [currentRepoUrl, setCurrentRepoUrl] = useState("");
  const [branches, setBranches] = useState<string[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const { toast } = useToast();
  const { t } = useTranslation();
  const { onNotification, onInteractive, sendMessage } = useWebSocket();

  const [actionDialog, setActionDialog] = useState<{
    show: boolean;
    title: string;
    description: string;
    fields: InputInterface[];
    onSubmit: (data: any) => void;
  }>({
    show: false,
    title: "",
    description: "",
    fields: [],
    onSubmit: () => {},
  });

  const [deleteVersion] = useDeleteVersionMutation();
  const [deployApplication] = useDeployApplicationMutation();
  const [uploadApplication] = useUploadApplicationMutation();
  const [startApplication] = useStartApplicationMutation();
  const [stopApplication] = useStopApplicationMutation();
  const { data: branchesData, isFetching: isLoadingBranches, error: branchesError } = useLoadBranchesQuery({ repoUrl: currentRepoUrl }, { skip: !currentRepoUrl });

  useEffect(() => {
    application.versions && setVersions(application.versions);
  }, [application.versions]);

  useEffect(() => {
    if (branchesData?.branches) {
      setBranches(branchesData.branches);
    }
  }, [branchesData]);

  useEffect(() => {
    if (branchesError) {
      toast({
        title: t("common.error"),
        description: t("dashboard.applications.errors.branchesFetchFailed"),
        variant: "destructive",
      });
    }
  }, [branchesError, t]);

  useEffect(() => {
    const unsubProgress = onNotification((message: ProgressMessage) => {
      if (message.type === "progress") {
        toast({
          title: t("dashboard.applications.actions.deploymentProgress"),
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
        title: message.title || t("dashboard.applications.actions.required"),
        description: message.message || "",
        fields: message.inputs.map((input) => ({
          ...input,
          // Add validation for port number
          validate:
            input.name === "port"
              ? (value: string) => {
                  const port = parseInt(value);
                  if (isNaN(port) || port < 1 || port > 65535) {
                    return t("dashboard.applications.errors.portInvalid");
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
            title: t("dashboard.applications.actions.portConfiguration"),
            description: t("dashboard.applications.actions.portExposed", {
              port: data.port,
            }),
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
  }, [onNotification, onInteractive, sendMessage, toast, t]);

  const handleVersionAction = async (appId: string, versionId: string, action: "start" | "stop") => {
    try {
      if (action === "start") {
        await startApplication({ appId, versionId }).unwrap();
      } else if (action === "stop") {
        await stopApplication({ appId, versionId }).unwrap();
      }

      toast({
        title: t("common.success"),
        description: t(`applications.actions.${action}Success`),
      });

      router.reload({ only: ["application"] });
    } catch (error: any) {
      toast({
        title: t("common.error"),
        description: error,
        variant: "destructive",
      });
    }
  };

  const handleDeleteVersion = async (appId: string, versionId: string) => {
    try {
      await deleteVersion({ appId, versionId }).unwrap();
      toast({
        title: "Success",
        description: "Version deleted successfully",
      });
      router.reload({ only: ["application"] });
    } catch (error) {
      console.error(error);
    }
  };

  let debounceTimer: ReturnType<typeof setTimeout> | null = null;

  const handleRepoUrlChange = (url: string) => {
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }

    debounceTimer = setTimeout(() => {
      if (!url) {
        setBranches([]);
        setCurrentRepoUrl("");
      } else {
        setCurrentRepoUrl(url);
      }
    }, 1000);
  };
  const handleVersionSubmit = async (values: z.infer<typeof uploadFormSchema>, file: File | null) => {
    setIsUploading(true);

    try {
      if (!values.repoUrl && !file) {
        toast({
          title: "Error",
          description: "You must provide a repo URL or upload a .zip file",
          variant: "destructive",
        });
        return;
      }

      if (values.repoUrl && values.branch) {
        await deployApplication({
          appId: application.id,
          repoUrl: values.repoUrl,
          branch: values.branch,
        }).unwrap();
      }

      if (file) {
        await uploadApplication({
          appId: application.id,
          file: file,
        }).unwrap();
      }

      toast({
        title: "Success",
        description: "Version created successfully",
      });

      // refrescar app o versiones si quieres
    } catch (err: any) {
      toast({
        title: "Error",
        description: err?.message || "Something went wrong",
        variant: "destructive",
      });
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="p-6 space-y-6 max-w-7xl mx-auto">
      {/* Header Section */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-bold ">{application.appName}</h1>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30">
              Running
            </Badge>
            <span className="text-sm text-muted-foreground">ID: {application.id}</span>
          </div>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Overview Section */}
        <Card className="border-border/50">
          <CardHeader>
            <CardTitle>Overview</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Created At</p>
                <p className="text-sm font-medium">{application.createdAt}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Updated At</p>
                <p className="text-sm font-medium">{application.updatedAt}</p>
              </div>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Description</p>
              <p className="text-sm">{application.description}</p>
            </div>
          </CardContent>
        </Card>

        {/* Metrics Section */}
        <Card className=" border-border/50">
          <CardHeader>
            <CardTitle>Metrics</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm">CPU Usage</span>
                <span className="text-sm font-medium">{application.cpuUsage.toFixed(2)}%</span>
              </div>
              <Progress value={application.cpuUsage} className="h-2" />
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm">Memory Usage</span>
                <span className="text-sm font-medium">{application.memoryUsage.toFixed(2)}%</span>
              </div>
              <Progress value={application.memoryUsage} className="h-2" />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Uptime</p>
                <p className="text-sm font-medium">{application.uptime}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Requests/min</p>
                <p className="text-sm font-medium">{application.requestsPerMin}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* API Versions Section */}
        <Card className="md:col-span-2  border-border/50">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>API Versions</CardTitle>
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="outline" size="sm">
                    <Plus className="w-4 h-4 mr-2" />
                    New Version
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Create New Version</DialogTitle>
                  </DialogHeader>
                  <ApplicationForm
                    mode="create-version"
                    onSubmit={handleVersionSubmit}
                    isUploading={isUploading}
                    branches={branches}
                    isLoadingBranches={isLoadingBranches}
                    onRepoUrlChange={handleRepoUrlChange}
                  />
                </DialogContent>
              </Dialog>
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Version</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Path</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created At</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {versions?.length ? (
                  versions.map((version, i) => (
                    <TableRow key={i}>
                      <TableCell className="font-mono">{version.versionTag}</TableCell>
                      <TableCell>{version.description}</TableCell>
                      <TableCell>
                        <a target="_blank" href={`/${version.versionTag}/${sanitizeAppName(application.appName)}/`}>{`/${version.versionTag}/${sanitizeAppName(application.appName)}/`}</a>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline" className="text-xs capitalize">
                          {version.status}
                        </Badge>
                      </TableCell>
                      <TableCell>{new Date(version.createdAt).toLocaleDateString()}</TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 text-blue-400 hover:bg-blue-400/10"
                          onClick={() => handleVersionAction(application.id, version.id, "start")}
                          disabled={version.status === "active"}>
                          <CirclePlay className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 text-yellow-400 hover:bg-yellow-400/10"
                          onClick={() => handleVersionAction(application.id, version.id, "stop")}
                          disabled={version.status !== "active"}>
                          <Pause className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 text-red-400 hover:bg-red-400/10"
                          onClick={() => handleDeleteVersion(application.id, version.id)}
                          disabled={version.status === "Running"}>
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center text-muted-foreground">
                      No versions found.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>

      <Dialog open={actionDialog.show} onOpenChange={(open) => !open && setActionDialog({ ...actionDialog, show: false })}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{actionDialog.title}</DialogTitle>
            <DialogDescription>{actionDialog.description}</DialogDescription>
          </DialogHeader>
          <DynamicForm fields={actionDialog.fields} onSubmit={actionDialog.onSubmit} />
        </DialogContent>
      </Dialog>
    </div>
  );
};

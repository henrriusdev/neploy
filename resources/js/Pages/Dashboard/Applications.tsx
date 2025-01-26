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
import { debounce } from "lodash";
import { Grid, List, PlusCircle } from "lucide-react";
import * as React from "react";
import { useEffect, useMemo, useState } from "react";
import * as z from "zod";
import { ApplicationsProps } from "@/types/props";
import { useGetAllApplicationsQuery, useLoadBranchesQuery, useCreateApplicationMutation, useDeployApplicationMutation, useUploadApplicationMutation, useStartApplicationMutation, useStopApplicationMutation, useDeleteApplicationMutation } from "@/services/api/applications";
import { useTranslation } from 'react-i18next';
import '@/i18n';

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
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [isUploading, setIsUploading] = useState(false);
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);
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
  const { t } = useTranslation();
  const { onNotification, onInteractive, sendMessage } = useWebSocket();

  const { data: applications, refetch: refreshApplications, error: applicationsError } = useGetAllApplicationsQuery(undefined, {
    // Refetch cada 30 segundos
    pollingInterval: 30000,
    refetchOnFocus: true,
    refetchOnReconnect: true
  });

  const { data: branchesData, isFetching: isLoadingBranches, error: branchesError } = useLoadBranchesQuery(
    { repoUrl: "" },
    {
      skip: true,
    }
  );

  useEffect(() => {
    if (applicationsError) {
      toast({
        title: t('common.error'),
        description: t('applications.errors.fetchFailed'),
        variant: "destructive",
      });
    }
  }, [applicationsError, t]);

  useEffect(() => {
    if (branchesError) {
      toast({
        title: t('common.error'),
        description: t('applications.errors.branchesFetchFailed'),
        variant: "destructive",
      });
    }
  }, [branchesError, t]);

  useEffect(() => {
    if (branchesData?.branches) {
      setBranches(branchesData.branches);
    }
  }, [branchesData]);

  const debouncedFetchBranches = useMemo(
    () => debounce((repoUrl: string) => {
      if (!repoUrl) {
        setBranches([]);
        return;
      }
      // Refetch con el nuevo repoUrl
      useLoadBranchesQuery({ repoUrl }, { skip: false });
    }, 1000),
    []
  );

  const [createApplication] = useCreateApplicationMutation();
  const [deployApplication] = useDeployApplicationMutation();
  const [uploadApplication] = useUploadApplicationMutation();
  const [startApplication] = useStartApplicationMutation();
  const [stopApplication] = useStopApplicationMutation();
  const [deleteApplication] = useDeleteApplicationMutation();

  const handleApplicationAction = async (
    appId: string,
    action: "start" | "stop" | "delete"
  ) => {
    try {
      if (action === "delete") {
        await deleteApplication({ appId });
      } else if (action === "start") {
        await startApplication({ appId });
      } else if (action === "stop") {
        await stopApplication({ appId });
      }

      toast({
        title: t('common.success'),
        description: t(`applications.actions.${action}Success`),
      });
      
      refreshApplications();
    } catch (error: any) {
      toast({
        title: t('common.error'),
        description: error.message || t(`applications.errors.${action}Failed`),
        variant: "destructive",
      });
    }
  };

  const onSubmit = async (
    values: z.infer<typeof uploadFormSchema>,
    file: File | null
  ) => {
    if (!values.appName) {
      toast({
        title: t('common.error'),
        description: t('applications.errors.nameRequired'),
        variant: "destructive",
      });
      return;
    }

    if (!file && !values.repoUrl) {
      toast({
        title: t('common.error'),
        description: t('applications.errors.fileOrRepoRequired'),
        variant: "destructive",
      });
      return;
    }

    setIsUploading(true);

    try {
      const response = await createApplication({
        appName: values.appName,
        description:
          values.description ||
          `Application created from ${
            file ? "file upload" : "repository " + values.repoUrl
          }`,
      });

      if ('error' in response) {
        throw new Error('Failed to create application');
      }

      const applicationId = response.data.id;

      if (values.repoUrl) {
        if (!values.branch) {
          toast({
            title: t('common.error'),
            description: t('applications.errors.branchRequired'),
            variant: "destructive",
          });
          return;
        }

        await deployApplication({
          appId: applicationId,
          repoUrl: values.repoUrl,
          branch: values.branch,
        });

        toast({
          title: t('common.success'),
          description: t('applications.actions.deploySuccess'),
        });
      }

      if (file) {
        await uploadApplication({
          appId: applicationId,
          file: file
        });

        toast({
          title: t('common.success'),
          description: t('applications.actions.uploadSuccess'),
        });
      }

      refreshApplications();
      setUploadDialogOpen(false);
    } catch (error: any) {
      toast({
        title: t('common.error'),
        description: error.message || t('applications.errors.unknown'),
        variant: "destructive",
      });
    } finally {
      setIsUploading(false);
    }
  };

  useEffect(() => {
    const unsubProgress = onNotification((message: ProgressMessage) => {
      if (message.type === "progress") {
        toast({
          title: t('applications.actions.deploymentProgress'),
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
        title: message.title || t('applications.actions.required'),
        description: message.message || "",
        fields: message.inputs.map((input) => ({
          ...input,
          // Add validation for port number
          validate:
            input.name === "port"
              ? (value: string) => {
                  const port = parseInt(value);
                  if (isNaN(port) || port < 1 || port > 65535) {
                    return t('applications.errors.portInvalid');
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
            title: t('applications.actions.portConfiguration'),
            description: t('applications.actions.portExposed', { port: data.port }),
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

  React.useEffect(() => {
    refreshApplications();
  }, []);

  React.useEffect(() => {
    const ws = new WebSocket(
      `${window.location.protocol === "https:" ? "wss:" : "ws:"}//${
        window.location.host
      }/ws`
    );

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === "APPLICATION_UPDATE") {
        // Refetch para obtener los datos actualizados
        refreshApplications();
      }
    };

    return () => {
      ws.close();
    };
  }, [refreshApplications]);

  return (
    <div className="space-y-6 p-3">
      {/* Stats Section */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">{t('dashboard.applications.stats.totalApplications')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {applications?.length || 0}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">{t('dashboard.applications.stats.runningApplications')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {applications?.filter((app) => app.status === "Running").length || 0}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">{t('dashboard.applications.stats.failedApplications')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {applications?.filter((app) => app.status === "Error").length || 0}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Actions Bar */}
      <div className="flex justify-between items-center">
          <Dialog open={uploadDialogOpen} onOpenChange={setUploadDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <PlusCircle className="mr-2 h-4 w-4" />
                {t('dashboard.applications.create')}
              </Button>
            </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle>{t('dashboard.applications.createNew.title')}</DialogTitle>
                <DialogDescription>
                  {t('dashboard.applications.createNew.description')}
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
        </div>

      <div className={viewMode === "grid" ? "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4" : "space-y-4"}>
        {applications?.map((app) => (
          <ApplicationCard
            key={app.id}
            app={app}
            onStart={() => handleApplicationAction(app.id, "start")}
            onStop={() => handleApplicationAction(app.id, "stop")}
            onDelete={() => handleApplicationAction(app.id, "delete")}
          />
        ))}
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
}

Applications.layout = (page: any) => {
  return <DashboardLayout>{page}</DashboardLayout>;
};

export default Applications;
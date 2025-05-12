import {useToast, useWebSocket} from "@/hooks";
import {
  useCreateApplicationMutation,
  useDeleteApplicationMutation,
  useDeployApplicationMutation,
  useGetAllApplicationsQuery,
  useLoadBranchesQuery,
  useStartApplicationMutation,
  useStopApplicationMutation,
  useUploadApplicationMutation,
} from "@/services/api/applications";
import {ActionMessage, ActionResponse, ApplicationsProps, Input, ProgressMessage,} from "@/types";
import {debounce} from "lodash";
import {useEffect, useMemo, useState} from "react";
import {useTranslation} from "react-i18next";
import {z} from "zod";
import {Card, CardContent, CardHeader, CardTitle} from "../ui/card";
import {Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger,} from "../ui/dialog";
import {Button} from "../ui/button";
import {PlusCircle} from "lucide-react";
import {ApplicationForm, DynamicForm} from "../forms";
import {ApplicationCard} from "../application-card";

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
      {message: "Must be a valid GitHub or GitLab repository URL"}
    )
    .optional(),
  branch: z.string().optional(),
});

export function Applications({
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
    fields: Input[];
    onSubmit: (data: any) => void;
  }>({
    show: false,
    title: "",
    description: "",
    fields: [],
    onSubmit: () => {
    },
  });
  const [currentRepoUrl, setCurrentRepoUrl] = useState("");
  const {toast} = useToast();
  const {t} = useTranslation();
  const {onNotification, onInteractive, sendMessage} = useWebSocket();

  const {
    data: applications,
    refetch: refreshApplications,
    error: applicationsError,
  } = useGetAllApplicationsQuery(undefined, {
    // Refetch cada 30 segundos
    pollingInterval: 30000,
    refetchOnFocus: true,
    refetchOnReconnect: true,
  });

  const {
    data: branchesData,
    isFetching: isLoadingBranches,
    error: branchesError,
  } = useLoadBranchesQuery(
    {repoUrl: currentRepoUrl},
    {skip: !currentRepoUrl}
  );

  useEffect(() => {
    if (applicationsError) {
      toast({
        title: t("common.error"),
        description: t("applications.errors.fetchFailed"),
        variant: "destructive",
      });
    }
  }, [applicationsError, t]);

  useEffect(() => {
    if (branchesError) {
      toast({
        title: t("common.error"),
        description: t("dashboard.applications.errors.branchesFetchFailed"),
        variant: "destructive",
      });
    }
  }, [branchesError]);

  useEffect(() => {
    if (branchesData?.branches) {
      setBranches(branchesData.branches);
    }
  }, [branchesData]);

  const debouncedFetchBranches = useMemo(
    () =>
      debounce((repoUrl: string) => {
        if (!repoUrl) {
          setBranches([]);
          setCurrentRepoUrl("");
          return;
        }
        setCurrentRepoUrl(repoUrl);
      }, 1000),
    []
  );

  const handleRepoUrlChange = (repoUrl: string) => {
    debouncedFetchBranches(repoUrl);
  };

  const [createApplication] = useCreateApplicationMutation();
  const [deployApplication] = useDeployApplicationMutation();
  const [uploadApplication] = useUploadApplicationMutation();
  const [startApplication] = useStartApplicationMutation();
  const [stopApplication] = useStopApplicationMutation();
  const [deleteApplication] = useDeleteApplicationMutation();

  const handleApplicationAction = async (appId: string) => {
    try {
      await deleteApplication({appId});

      toast({
        title: t("common.success"),
        description: t(`applications.actions.deleteSuccess`),
      });

      refreshApplications();
    } catch (error: any) {
      toast({
        title: t("common.error"),
        description: error.message || t(`applications.errors.deleteFailed`),
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
        title: t("common.error"),
        description: t("applications.errors.nameRequired"),
        variant: "destructive",
      });
      return;
    }

    if (!file && !values.repoUrl) {
      toast({
        title: t("common.error"),
        description: t("applications.errors.fileOrRepoRequired"),
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

      if ("error" in response) {
        throw new Error("Failed to create application");
      }

      const applicationId = response.data.id;

      if (values.repoUrl) {
        if (!values.branch) {
          toast({
            title: t("common.error"),
            description: t("applications.errors.branchRequired"),
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
          title: t("common.success"),
          description: t("applications.actions.deploySuccess"),
        });
      }

      if (file) {
        await uploadApplication({
          appId: applicationId,
          file: file,
        });

        toast({
          title: t("common.success"),
          description: t("applications.actions.uploadSuccess"),
        });
      }

      refreshApplications();
      setUploadDialogOpen(false);
    } catch (error: any) {
      toast({
        title: t("common.error"),
        description: error.message || t("applications.errors.unknown"),
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
          title: t("applications.actions.deploymentProgress"),
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
        title: message.title || t("applications.actions.required"),
        description: message.message || "",
        fields: message.inputs.map((input: Input) => ({
          ...input,
          // Add validation for port number
          validate:
            input.name === "port"
              ? (value: string) => {
                const port = parseInt(value);
                if (isNaN(port) || port < 1 || port > 65535) {
                  return t("applications.errors.portInvalid");
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
          setActionDialog((prev) => ({...prev, show: false}));

          // Show confirmation toast
          toast({
            title: t("applications.actions.portConfiguration"),
            description: t("applications.actions.portExposed", {
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

  useEffect(() => {
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
            <CardTitle className="text-sm font-medium">
              {t("dashboard.applications.stats.totalApplications")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {applications?.length || 0}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("dashboard.applications.stats.runningApplications")}
            </CardTitle>
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
            <CardTitle className="text-sm font-medium">
              {t("dashboard.applications.stats.failedApplications")}
            </CardTitle>
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
      <div className="flex justify-between items-center px-3 py-1">
        <h1 className="font-bold text-3xl">{t('dashboard.applications.title')}</h1>
        <Dialog open={uploadDialogOpen} onOpenChange={setUploadDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <PlusCircle className="mr-2 h-4 w-4"/>
              {t("dashboard.applications.create")}
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>
                {t("dashboard.applications.createNew.title")}
              </DialogTitle>
              <DialogDescription>
                {t("dashboard.applications.createNew.description")}
              </DialogDescription>
            </DialogHeader>
            <ApplicationForm
              onSubmit={onSubmit}
              isUploading={isUploading}
              branches={branches}
              isLoadingBranches={isLoadingBranches}
              onRepoUrlChange={handleRepoUrlChange}
            />
          </DialogContent>
        </Dialog>
      </div>

      <div
        className={
          viewMode === "grid"
            ? "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
            : "space-y-4"
        }>
        {applications?.map((app) => (
          <ApplicationCard
            key={app.id}
            app={app}
            onDelete={() => handleApplicationAction(app.id)}
          />
        ))}
      </div>

      <Dialog
        open={actionDialog.show}
        onOpenChange={(open) =>
          !open && setActionDialog({...actionDialog, show: false})
        }>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{actionDialog.title}</DialogTitle>
            <DialogDescription>{actionDialog.description}</DialogDescription>
          </DialogHeader>
          <DynamicForm
            fields={actionDialog.fields}
            onSubmit={actionDialog.onSubmit}
          />
        </DialogContent>
      </Dialog>
    </div>
  );
}

import React, { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { TechStacksSettingsProps } from "@/types/props";
import { TechIcon } from "@/components/icons/tech-icon";
import { useTranslation } from "react-i18next";
import { DialogButton } from "../forms/dialog-button";
import { TechStackForm } from "../forms/tech-stack-form";
import { useToast } from "@/hooks";
import { z } from "zod";
import {
  useCreateTechStackMutation,
  useDeleteTechStackMutation,
  useGetTechStacksQuery,
  useUpdateTechStackMutation,
} from "@/services/api/tech-stack";
import { Pencil, PlusCircle, Trash2 } from "lucide-react";
import { TooltipButton } from "../ui/tooltip-button";

const techStackSchema = z.object({
  name: z.string().min(2).max(64),
  description: z.string().min(2).max(128),
});

const TechStackTab: React.FC<TechStacksSettingsProps> = ({
  techStacks: initialTechStacks,
}) => {
  const { t } = useTranslation();
  const { toast } = useToast();
  const getTechStacks = useGetTechStacksQuery();
  const [createTechStack] = useCreateTechStackMutation();
  const [updateTechStack] = useUpdateTechStackMutation();
  const [deleteTechStack] = useDeleteTechStackMutation();
  const [openTechStackId, setOpenTechStackId] = useState<string | null>(null);
  const [techStacks, setTechStacks] = useState(initialTechStacks);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (getTechStacks.currentData) {
      setTechStacks(getTechStacks.currentData);
    }
  }, [getTechStacks.data]);

  async function update(
    techStackId: string,
    data: z.infer<typeof techStackSchema>
  ) {
    try {
      await updateTechStack({
        id: techStackId,
        name: data.name,
        description: data.description,
      }).unwrap();

      toast({
        title: t("settings.techStack.updateSuccess"),
        description: t("settings.techStack.updateSuccessDescription"),
      });
      setOpenTechStackId(null);
      getTechStacks.refetch();
    } catch (error) {
      toast({
        title: t("common.error"),
        description: t("settings.techStack.updateError"),
        variant: "destructive",
      });
    }
  }

  async function create(data: z.infer<typeof techStackSchema>) {
    try {
      await createTechStack({
        name: data.name,
        description: data.description,
      }).unwrap();
      toast({
        title: t("settings.techStack.createSuccess"),
        description: t("settings.techStack.createSuccessDescription"),
      });
      getTechStacks.refetch();
      setOpen(false);
    } catch (error) {
      toast({
        title: t("common.error"),
        description: t("settings.techStack.createError"),
        variant: "destructive",
      });
    }
  }

  async function del(techStackId: string) {
    try {
      await deleteTechStack(techStackId).unwrap();
      toast({
        title: t("settings.techStack.deleteSuccess"),
        description: t("settings.techStack.deleteSuccessDescription"),
      });
      getTechStacks.refetch();
    } catch (error) {
      toast({
        title: t("common.error"),
        description: t("settings.techStack.deleteError"),
        variant: "destructive",
      });
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <DialogButton
          title={t("settings.techStack.add")}
          open={open}
          onOpen={setOpen}
          description={t("settings.techStack.editDescriptionDialog")}
          icon={PlusCircle}
          buttonText={t("settings.techStack.add")}>
          <TechStackForm onSubmit={create} />
        </DialogButton>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("settings.techStack.title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("settings.techStack.tableLogo")}</TableHead>
                <TableHead>{t("settings.techStack.tableTechnology")}</TableHead>
                <TableHead>
                  {t("settings.techStack.tableDescription")}
                </TableHead>
                <TableHead>{t("settings.techStack.tableTotalApps")}</TableHead>
                <TableHead>{t("settings.techStack.tableActions")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {techStacks.map((techStack) => (
                <TableRow key={techStack.id}>
                  <TableCell>
                    <TechIcon name={techStack.name} size={40} />
                  </TableCell>
                  <TableCell>{techStack.name}</TableCell>
                  <TableCell>{techStack.description}</TableCell>
                  <TableCell>
                    <Badge
                      variant={
                        techStack.applications?.length > 0
                          ? "destructive"
                          : "default"
                      }>
                      {techStack.applications?.length ?? 0}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <div className="flex space-x-2">
                      <DialogButton
                        title={t("settings.techStack.edit")}
                        description={t(
                          "settings.techStack.editDescriptionDialog"
                        )}
                        icon={Pencil}
                        variant="tooltip"
                        open={openTechStackId === techStack.id}
                        onOpen={() => setOpenTechStackId(techStack.id)}
                        buttonText={t("settings.techStack.editAction")}>
                        <TechStackForm
                          defaultValues={techStack}
                          onSubmit={(data) => update(techStack.id, data)}
                        />
                      </DialogButton>
                      <TooltipButton
                        icon={Trash2}
                        tooltip={t("settings.techStack.deleteTooltip")}
                        variant="destructive"
                        size="icon"
                        disabled={techStack.applications?.length > 0}
                        onClick={() => del(techStack.id)}>
                        {t("settings.techStack.deleteAction")}
                      </TooltipButton>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default TechStackTab;

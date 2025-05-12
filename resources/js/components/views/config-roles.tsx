import React, { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { RoleIcon } from "@/components/icons/role-icon";
import {Eye, Pencil, PlusCircle, Trash2, Users} from "lucide-react";
import { TooltipButton } from "@/components/ui/tooltip-button";
import { useTranslation } from "react-i18next";
import { DialogButton } from "../forms/dialog-button";
import { RoleForm } from "../forms/role-form";
import { RolesSettingsProps, RoleWithUsers } from "@/types";
import { useToast } from "@/hooks";
import { roleSchema } from "@/lib/validations/role";
import { z } from "zod";
import {
  useCreateRoleMutation,
  useDeleteRoleMutation,
  useGetRolesQuery,
  useUpdateRoleMutation,
} from "@/services/api/role";
import {RoleUserManagerDialog} from "@/components/role-user-manager-dialog";

const RolesTab: React.FC<RolesSettingsProps> = ({ roles: initialRoles }) => {
  const { t } = useTranslation();
  const { toast } = useToast();
  const getRoles = useGetRolesQuery();
  const [createRole] = useCreateRoleMutation();
  const [updateRole] = useUpdateRoleMutation();
  const [deleteRole] = useDeleteRoleMutation();
  const [openRoleId, setOpenRoleId] = useState<string | null>(null);
  const [roles, setRoles] = useState<RoleWithUsers[]>(initialRoles);
  const [openManageUsersRole, setOpenManageUsersRole] = useState("")

  const [open, setOpen] = useState(false);

  useEffect(() => {
    console.log(getRoles.currentData, getRoles.data);
    if (getRoles.currentData) {
      setRoles(getRoles.currentData);
    }
  }, [getRoles.data]);

  async function update(roleId: string, data: z.infer<typeof roleSchema>) {
    try {
      await updateRole({
        id: roleId,
        name: data.name,
        description: data.description,
        icon: data.icon,
        color: data.color,
      }).unwrap();

      toast({
        title: t("settings.roles.updateSuccess"),
        description: t("settings.roles.updateSuccessDescription", {
          role: data.name,
        }),
      });
      setOpenRoleId(null);
      getRoles.refetch();
    } catch (error) {
      toast({
        title: t("common.error"),
        description: t("settings.roles.updateError"),
        variant: "destructive",
      });
    }
  }

  async function create(data: z.infer<typeof roleSchema>) {
    try {
      await createRole({
        name: data.name,
        description: data.description,
        icon: data.icon,
        color: data.color,
      }).unwrap();
      toast({
        title: t("settings.roles.createSuccess"),
        description: t("settings.roles.createSuccessDescription", {
          role: data.name,
        }),
      });
      getRoles.refetch();
      setOpen(false);
    } catch (error) {
      toast({
        title: t("common.error"),
        description: t("settings.roles.createError"),
        variant: "destructive",
      });
    }
  }

  async function del(roleId: string) {
    try {
      await deleteRole({ id: roleId }).unwrap();
      toast({
        title: t("settings.roles.deleteSuccess"),
        description: t("settings.roles.deleteSuccessDescription"),
      });
    } catch (error) {
      toast({
        title: t("common.error"),
        description: t("settings.roles.deleteError"),
        variant: "destructive",
      });
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex justify-between items-center">
          {t("settings.roles.title")}
          <DialogButton
            buttonText={t("settings.roles.add")}
            title={t("settings.roles.editDialog")}
            description={t("settings.roles.editDescriptionDialog")}
            icon={PlusCircle}
            open={open}
            variant="text"
            onOpen={setOpen}>
            <RoleForm
              onSubmit={create}
              onCancel={() => setOpen(false)}
              defaultValues={null}
            />
          </DialogButton>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t("settings.roles.tableIcon")}</TableHead>
              <TableHead>{t("settings.roles.tableName")}</TableHead>
              <TableHead>{t("settings.roles.tableDescription")}</TableHead>
              <TableHead>{t("settings.roles.tableUsers")}</TableHead>
              <TableHead>{t("settings.roles.tableActions")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {roles.map((role) => (
              <TableRow key={role.id}>
                <TableCell>
                  <RoleIcon icon={role.icon} color={role.color} size={60} />
                </TableCell>
                <TableCell>{role.name}</TableCell>
                <TableCell>{role.description}</TableCell>
                <TableCell>{role.users.length}</TableCell>
                <TableCell>
                  <div className="flex items-center gap-2">
                    <TooltipButton
                      tooltip={t("settings.roles.showAction")}
                      icon={Users}
                      variant="ghost"
                      size="icon"
                      onClick={() => setOpenManageUsersRole(role.id)}
                    />
                    <RoleUserManagerDialog
                      open={openManageUsersRole === role.id}
                      onOpenChange={() => setOpenManageUsersRole(null)}
                      roleId={role.id}
                      roleName={role.name}
                      assignedUsers={role.users}
                      />
                    <DialogButton
                      buttonText={t("settings.roles.editAction")}
                      title={t("settings.roles.editDialog")}
                      description={t("settings.roles.editDescriptionDialog")}
                      icon={Pencil}
                      variant="tooltip"
                      open={openRoleId === role.id}
                      onOpen={() => setOpenRoleId(role.id)}>
                      <RoleForm
                        defaultValues={role}
                        onSubmit={(data) => {
                          update(role.id, data);
                        }}
                        onCancel={() => {
                          setOpenRoleId(null);
                        }}
                      />
                    </DialogButton>
                    <TooltipButton
                      tooltip={t("settings.roles.deleteAction")}
                      icon={Trash2}
                      variant="destructive"
                      size="icon"
                      disabled={role.users.length > 0}
                      onClick={() => del(role.id)}
                    />
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};

export default RolesTab;

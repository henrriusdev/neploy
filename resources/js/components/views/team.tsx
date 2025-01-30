import { useToast } from "@/hooks";
import { TeamProps } from "@/types";
import { router } from "@inertiajs/react";
import { PlusCircle, Trash } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { Badge } from "../ui/badge";
import { Button } from "../ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";

interface InviteMemberData {
  email: string;
  role: string;
}

export function Team({ user, teamName, logoUrl, team, roles }: TeamProps) {
  const [open, setOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState<InviteMemberData>({
    email: "",
    role: "",
  });
  const [teamState, setTeamState] = useState(team);
  const { toast } = useToast();
  const { t } = useTranslation();

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    const inviteData = {
      email: formData.email,
      role: formData.role,
    };

    router.post("/users/invite", inviteData, {
      onSuccess: () => {
        toast({
          title: t("dashboard.team.inviteSuccess"),
          description: t("dashboard.team.inviteSuccess"),
        });
        setOpen(false);
        setFormData({ email: "", role: "" });
      },
      onError: () => {
        toast({
          title: t("dashboard.team.inviteError"),
          description: t("dashboard.team.inviteError"),
          variant: "destructive",
        });
      },
      onFinish: () => setIsLoading(false),
    });
  };

  const handleRemoveMember = async (memberId: string) => {
    if (!confirm(t("dashboard.team.confirmRemove"))) return;

    router.delete(`/users/${memberId}`, {
      onSuccess: () => {
        toast({
          title: t("dashboard.team.removeSuccess"),
          description: t("dashboard.team.removeSuccess"),
        });
        setTeamState(teamState.filter((member) => member.id !== memberId));
      },
      onError: () => {
        toast({
          title: t("dashboard.team.removeError"),
          description: t("dashboard.team.removeError"),
          variant: "destructive",
        });
      },
    });
  };

  return (
    <div className="container mx-auto py-6">
      <Card>
        <CardHeader>
          <div className="flex justify-between items-center">
            <div>
              <CardTitle>{t("dashboard.team.title")}</CardTitle>
              <CardDescription>
                {t("dashboard.team.description")}
              </CardDescription>
            </div>
            <Dialog open={open} onOpenChange={setOpen}>
              <DialogTrigger asChild>
                <Button>
                  <PlusCircle className="mr-2 h-4 w-4" />
                  {t("dashboard.team.inviteMember")}
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>{t("dashboard.team.inviteMember")}</DialogTitle>
                  <DialogDescription>
                    {t("dashboard.team.inviteDescription")}
                  </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleInvite} className="space-y-4">
                  <div>
                    <Label htmlFor="email">{t("dashboard.team.email")}</Label>
                    <Input
                      id="email"
                      type="email"
                      value={formData.email}
                      onChange={(e) =>
                        setFormData({ ...formData, email: e.target.value })
                      }
                      required
                    />
                  </div>
                  <div>
                    <Label htmlFor="role">{t("dashboard.team.role")}</Label>
                    <Select
                      value={formData.role}
                      onValueChange={(value) =>
                        setFormData({ ...formData, role: value })
                      }>
                      <SelectTrigger>
                        <SelectValue
                          placeholder={t("dashboard.team.selectRole")}
                        />
                      </SelectTrigger>
                      <SelectContent>
                        {roles.map((role) => (
                          <SelectItem key={role.name} value={role.name}>
                            {role.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <Button type="submit" className="w-full" disabled={isLoading}>
                    {isLoading
                      ? t("dashboard.team.inviting")
                      : t("dashboard.team.invite")}
                  </Button>
                </form>
              </DialogContent>
            </Dialog>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("dashboard.team.member")}</TableHead>
                <TableHead>{t("dashboard.team.role")}</TableHead>
                <TableHead>{t("dashboard.team.status")}</TableHead>
                <TableHead className="text-right">
                  {t("dashboard.team.actions")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {teamState.map((member) => (
                <TableRow key={member.id}>
                  <TableCell className="font-medium">
                    <div className="flex items-center space-x-4">
                      <Avatar>
                        <AvatarImage
                          src={`https://unavatar.io/${member.provider}/${member.username}`}
                          alt={member.firstName + " " + member.lastName}
                        />
                        <AvatarFallback>
                          {member.firstName
                            .split(" ")
                            .map((n) => n[0])
                            .join("")}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <div className="font-medium">
                          {member.firstName + " " + member.lastName}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          {member.email}
                        </div>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    {member.roles.map((role) => (
                      <Badge
                        key={role.name}
                        variant="default"
                        style={{ backgroundColor: role.color }}>
                        {role.name}
                      </Badge>
                    ))}
                  </TableCell>
                  <TableCell>
                    <span className="text-xs">Active</span>
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="destructive"
                      size="icon"
                      onClick={() => handleRemoveMember(member.id)}>
                      <Trash className="h-4 w-4 text-destructive-foreground" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}

import * as React from "react";
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
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { PlusCircle } from "lucide-react";
import DashboardLayout from "@/components/Layouts/DashboardLayout";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { router } from "@inertiajs/react";
import { useToast } from "@/hooks/use-toast";
import { RoleIcon } from "@/components/RoleIcon";
import { TrashIcon } from "@radix-ui/react-icons";

interface TeamMember {
  id: string;
  username: string;
  email: string;
  firstName: string;
  lastName: string;
  provider: string;
  roles: Array<{
    name: string;
    description: string;
    icon: string;
    color: string;
  }>;
  techStacks: Array<{
    name: string;
    description: string;
  }>;
}

const defaultTeam: TeamMember[] = [
  {
    id: "1",
    username: "johndoe",
    email: "john@example.com",
    firstName: "John",
    lastName: "Doe",
    provider: "github",
    roles: [
      {
        name: "Admin",
        description: "Administrator",
        icon: "crown",
        color: "#FFB020",
      },
    ],
    techStacks: [
      {
        name: "React",
        description: "Frontend Framework",
      },
    ],
  },
];

interface TeamProps {
  user?: {
    name: string;
    email: string;
    avatar: string;
  };
  teamName?: string;
  logoUrl?: string;
  team?: TeamMember[];
  roles?: Array<{
    name: string;
    description: string;
    icon: string;
    color: string;
  }>;
}

interface InviteMemberData {
  email: string;
  role: string;
}

function Team({
  user,
  teamName,
  logoUrl,
  team = defaultTeam,
  roles,
}: TeamProps) {
  const [open, setOpen] = React.useState(false);
  const [isLoading, setIsLoading] = React.useState(false);
  const [formData, setFormData] = React.useState<InviteMemberData>({
    email: "",
    role: "member", // default role
  });
  const [teamState, setTeam] = React.useState(team);
  const { toast } = useToast();
  console.log(roles);

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    
    router.post("/users/invite", formData, {
      onSuccess: () => {
        toast({
          title: "Success",
          description: "Invitation sent successfully",
        });
        setFormData({ email: "", role: "member" });
        setOpen(false);
        
        // Refresh the team list
        router.visit("/team", {
          method: 'get',
          onSuccess: (response) => {
            setTeam(response.props.team);
          },
        });
      },
      onError: (error) => {
        toast({
          title: "Error",
          description: error || "Failed to send invitation",
          variant: "destructive",
        });
      },
      onFinish: () => {
        setIsLoading(false);
      }
    });
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Team Members</h2>
          <p className="text-muted-foreground">
            Manage your team members and their access levels
          </p>
        </div>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button variant="outline">
              <PlusCircle className="mr-2 h-4 w-4" />
              Add Member
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Invite Team Member</DialogTitle>
              <DialogDescription>
                Enter the email address and select a role for the new team
                member.
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleInvite} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="email">Email address</Label>
                <Input
                  id="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) =>
                    setFormData((prev) => ({ ...prev, email: e.target.value }))
                  }
                  placeholder="colleague@company.com"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="role">Role</Label>
                <Select
                  value={formData.role}
                  onValueChange={(value) =>
                    setFormData((prev) => ({ ...prev, role: value }))
                  }>
                  <SelectTrigger className="h-20">
                    <SelectValue placeholder="Select a role" />
                  </SelectTrigger>
                  <SelectContent>
                    {roles?.map((role) => (
                      <SelectItem key={role.name} value={role.name}>
                        <div className="!flex items-center gap-2 justify-start flex-row">
                          <RoleIcon icon={role.icon} color={role.color} />
                          <div>
                            <span className="capitalize font-bold text-gray-100">{role.name}</span>
                            <span className="block text-xs text-gray-200">
                              {role.description}
                            </span>
                          </div>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={isLoading}>
                  {isLoading ? "Sending..." : "Send Invitation"}
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Team Overview</CardTitle>
          <CardDescription>
            A list of all team members including their roles and tech stacks.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Member</TableHead>
                <TableHead>Roles</TableHead>
                <TableHead>Tech Stacks</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {teamState.length > 0 ? (
                teamState.map((member) => (
                  <TableRow key={member.id}>
                    <TableCell>
                      <div className="flex items-center space-x-4">
                        <Avatar>
                          <AvatarImage
                            src={`https://unavatar.io/${member.provider === 'github' ? `github/${member.username}` : member.email}`}
                            alt={`${member.firstName} ${member.lastName}`}
                          />
                          <AvatarFallback>
                            {member.firstName[0]}
                            {member.lastName[0]}
                          </AvatarFallback>
                        </Avatar>
                        <div>
                          <div className="font-medium">
                            {member.firstName} {member.lastName}
                          </div>
                          <div className="text-sm text-muted-foreground">
                            {member.email}
                          </div>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex flex-wrap gap-2">
                        {member.roles && member.roles.map((role) => (
                          <Badge
                            key={role.name}
                            variant="secondary"
                            style={{
                              backgroundColor: role.color + "20",
                              color: role.color,
                            }}>
                            {role.name}
                          </Badge>
                        ))}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex flex-wrap gap-2">
                        {member.techStacks && member.techStacks.map((tech) => (
                          <Badge
                            key={tech.name}
                            variant="outline"
                            className="text-foreground">
                            {tech.name}
                          </Badge>
                        ))}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex space-x-2">
                        <Button variant="destructive" size="icon" className="flex items-center justify-center">
                          <TrashIcon className="!h-7 !w-7" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={4} className="h-24 text-center">
                    No team members found.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}

Team.layout = (page: any) => {
  return (
    <DashboardLayout>
      {page}
    </DashboardLayout>
  );
};

export default Team;

import * as React from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
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
import { Controller } from "react-hook-form";
import { InputAutoComplete } from "@/components/forms";
import { ColorPicker } from "@/components/forms/color-picker";
import { RoleIcon } from "@/components/icons/role-icon";
import { Trash2 } from "lucide-react";
import { icons } from "@/lib/icons";
import { RoleForm } from "@/components/forms";

const roleSchema = z.object({
  name: z.string().min(1, "Role name is required"),
  description: z.string().min(1, "Description is required"),
  icon: z.string().min(1, "Icon is required"),
  color: z.string().min(1, "Color is required"),
});

interface Props {
  onNext: () => void;
  onBack: () => void;
  roles: any[];
  setRoles: (roles: any[]) => void;
}

export function RolesStep({ onNext, onBack, roles, setRoles }: Props) {
  const onSubmit = (data: z.infer<typeof roleSchema>) => {
    setRoles([...roles, data]);
  };

  return (
    <Card className="w-full max-w-screen-md mx-auto">
      <CardHeader>
        <CardTitle>Create Roles</CardTitle>
        <CardDescription>Define roles for your organization</CardDescription>
      </CardHeader>
      <CardContent>
        <RoleForm 
          onSubmit={onSubmit}
          renderFooter={(form) => (
            <div className="flex justify-between mt-6">
              <Button type="button" variant="outline" onClick={onBack}>
                Back
              </Button>
              <div className="space-x-2">
                <Button type="submit">Add Role</Button>
                <Button 
                  type="button" 
                  onClick={onNext}
                  disabled={roles.length === 0}>
                  Next
                </Button>
              </div>
            </div>
          )}
        />
        {roles.length > 0 && (
          <div className="mt-6">
            <h3 className="font-semibold text-xl">Roles Selected</h3>
            <ul>
              {roles.map((role, index) => (
                <li
                  key={index}
                  className="flex justify-between items-center space-x-2 my-1">
                  <div className="flex items-center space-x-2 space-y-2">
                    <RoleIcon icon={role.icon} color={role.color} />
                    <div>
                      <p className="font-semibold">{role.name}</p>
                      <p>{role.description}</p>
                    </div>
                  </div>
                  <Button
                    type="button"
                    variant="destructive"
                    size="icon"
                    className="w-12 h-12"
                    onClick={() => {
                      setRoles(roles.filter((r, i) => i !== index));
                    }}>
                    <Trash2 className="!w-7 !h-7" />
                  </Button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

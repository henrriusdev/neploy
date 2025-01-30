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
import { InputAutoComplete } from "@/components/input-autocomplete";
import { ColorPicker } from "@/components/forms/color-picker";
import { RoleIcon } from "@/components/role-icon";
import { Trash2 } from "lucide-react";
import { icons } from "@/lib/icons";

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

export default function RolesStep({ onNext, onBack, roles, setRoles }: Props) {
  const form = useForm<z.infer<typeof roleSchema>>({
    resolver: zodResolver(roleSchema),
    defaultValues: {
      name: "",
      description: "",
      icon: "",
      color: "#000000",
    },
  });

  const onSubmit = (data: z.infer<typeof roleSchema>) => {
    setRoles([...roles, data]);
    form.reset();
  };

  return (
    <Card className="w-full max-w-screen-md mx-auto">
      <CardHeader>
        <CardTitle>Create Roles</CardTitle>
        <CardDescription>Define roles for your organization</CardDescription>
      </CardHeader>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <CardContent className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Role Name</FormLabel>
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
                    <Textarea {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="icon"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Icon</FormLabel>
                  <FormControl>
                    <div className="flex gap-2 items-center">
                      <div className="flex-grow">
                        <Controller
                          control={form.control}
                          name="icon"
                          render={({ field }) => (
                            <InputAutoComplete
                              field={field}
                              OPTIONS={icons.map((icon) => ({
                                value: icon,
                                label: icon,
                                icon: (
                                  <RoleIcon
                                    icon={icon}
                                    color={form.getValues("color") || "#000000"}
                                    size={24}
                                  />
                                ),
                              }))}
                              placeholder="Search for an icon..."
                            />
                          )}
                        />
                      </div>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="color"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Color</FormLabel>
                  <FormControl>
                    <Controller
                      control={form.control}
                      name="color"
                      render={({ field }) => <ColorPicker field={field} />}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="flex items-start flex-col space-y-2">
              <h2 className="font-semibold text-lg">Icon Preview</h2>
              <RoleIcon icon={form.watch("icon")} color={form.watch("color")} />
            </div>

            {roles.length > 0 && (
              <div>
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
          <CardFooter className="flex justify-between">
            <Button type="button" variant="outline" onClick={onBack}>
              Back
            </Button>
            <div className="space-x-2">
              <Button type="submit">Add Role</Button>
              <Button type="button" onClick={onNext}>
                Next
              </Button>
            </div>
          </CardFooter>
        </form>
      </Form>
    </Card>
  );
}

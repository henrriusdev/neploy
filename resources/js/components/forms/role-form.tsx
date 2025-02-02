import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, Controller, UseFormReturn } from "react-hook-form";
import * as z from "zod";
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
import { InputAutoComplete } from "@/components/forms/input-autocomplete";
import { ColorPicker } from "@/components/forms";
import { RoleIcon } from "@/components/icons";
import { icons } from "@/lib/icons";
import { roleSchema } from "@/lib/validations/role";
import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

interface RoleFormProps {
  defaultValues?: z.infer<typeof roleSchema>;
  onSubmit: (data: z.infer<typeof roleSchema>) => void;
  onCancel?: () => void;
  renderFooter?: (form: UseFormReturn<z.infer<typeof roleSchema>>) => React.ReactNode;
}

export function RoleForm({ defaultValues, onSubmit, onCancel, renderFooter }: RoleFormProps) {
  const form = useForm<z.infer<typeof roleSchema>>({
    resolver: zodResolver(roleSchema),
    defaultValues: defaultValues || {
      name: "",
      description: "",
      icon: "",
      color: "#000000",
    },
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
        {renderFooter ? (
          renderFooter(form)
        ) : (
          <DialogFooter>
            <Button variant="outline" type="button" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="submit">Save Changes</Button>
          </DialogFooter>
        )}
      </form>
    </Form>
  );
}

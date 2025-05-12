import {zodResolver} from "@hookform/resolvers/zod";
import {useForm, UseFormReturn} from "react-hook-form";
import * as z from "zod";
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage,} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {InputAutoComplete} from "@/components/forms/input-autocomplete";
import {techIcons} from "@/lib/icons";
import {Button} from "@/components/ui/button";
import {DialogFooter} from "@/components/ui/dialog";
import {useTranslation} from "react-i18next";

const techStackSchema = z.object({
  name: z.string().min(2).max(64),
  description: z.string().min(2).max(128),
});

interface TechStackFormProps {
  defaultValues?: z.infer<typeof techStackSchema>;
  onSubmit: (data: z.infer<typeof techStackSchema>) => void;
  onCancel?: () => void;
  renderFooter?: (form: UseFormReturn<z.infer<typeof techStackSchema>>) => React.ReactNode;
}

export function TechStackForm({defaultValues, onSubmit, onCancel, renderFooter}: TechStackFormProps) {
  const {t} = useTranslation();

  const form = useForm<z.infer<typeof techStackSchema>>({
    resolver: zodResolver(techStackSchema),
    defaultValues: defaultValues || {
      name: "",
      description: "",
    },
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="name"
          render={({field}) => (
            <FormItem>
              <FormLabel>{t("settings.techStack.name")}</FormLabel>
              <FormControl>
                <InputAutoComplete
                  field={field}
                  OPTIONS={techIcons.map((icon) => ({
                    label: icon.name,
                    value: icon.name,
                  }))}
                  placeholder={t("settings.techStack.name")}
                />
              </FormControl>
              <FormMessage/>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="description"
          render={({field}) => (
            <FormItem>
              <FormLabel>{t("settings.techStack.description")}</FormLabel>
              <FormControl>
                <Input {...field} placeholder={t("settings.techStack.description")}/>
              </FormControl>
              <FormMessage/>
            </FormItem>
          )}
        />

        {renderFooter ? (
          renderFooter(form)
        ) : (
          <DialogFooter>
            {onCancel && (
              <Button type="button" variant="outline" onClick={onCancel}>
                {t("common.cancel")}
              </Button>
            )}
            <Button type="submit">{t("common.save")}</Button>
          </DialogFooter>
        )}
      </form>
    </Form>
  );
}

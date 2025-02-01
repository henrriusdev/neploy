import * as React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Switch } from "@/components/ui/switch";
import { Input } from "@/components/ui/input";
import { GeneralSettingsProps } from "@/types/props";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useTranslation } from "react-i18next";
import { useUpdateMetadataMutation } from "@/services/api/metadata";
import { useToast } from "@/hooks";

const formSchema = z.object({
  teamName: z.string().min(1, "Team name is required"),
  logoUrl: z.string().url("Must be a valid URL"),
  language: z.enum(["en", "es", "fr", "pt", "zh"]),
  darkMode: z.boolean(),
  emailNotifications: z.boolean(),
});

type FormValues = z.infer<typeof formSchema>;

const GeneralTab: React.FC<GeneralSettingsProps> = ({
  teamName: originalTeamName,
  logoUrl: originalLogoUrl,
  language: originalLanguage,
}) => {
  const { t } = useTranslation();
  const { toast } = useToast();
  const [updateMetadata] = useUpdateMetadataMutation();
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      teamName: originalTeamName,
      logoUrl: originalLogoUrl,
      language: originalLanguage as "en" | "es" | "fr" | "pt" | "zh",
      darkMode: false,
      emailNotifications: false,
    },
  });

  async function onSubmit(values: FormValues) {
    try {
      await updateMetadata({
        data: {
          teamName: values.teamName,
          logoUrl: values.logoUrl,
          language: values.language,
        },
      }).unwrap();

      // Show success toast or notification
      toast({
        title: t("common.success"),
        description: t("settings.general.updateSuccess"),
      });
    } catch (error) {
      console.error("Failed to save settings:", error);
      // Show error toast or notification
      toast({
        title: t("common.error"),
        description: t("settings.general.updateError"),
        variant: "destructive",
      });
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("settings.general.title")}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid gap-4 grid-cols-6">
              <FormField
                control={form.control}
                name="teamName"
                render={({ field }) => (
                  <FormItem className="space-y-2 col-span-3">
                    <FormLabel>{t("settings.general.teamName")}</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("settings.general.teamNamePlaceholder")}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="logoUrl"
                render={({ field }) => (
                  <FormItem className="space-y-2 col-span-3">
                    <FormLabel>{t("settings.general.logoUrl")}</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("settings.general.logoUrlPlaceholder")}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="language"
                render={({ field }) => (
                  <FormItem className="space-y-2 col-span-2">
                    <FormLabel>{t("settings.general.language")}</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={t(
                              "settings.general.languagePlaceholder"
                            )}
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="en">{t("languages.en")}</SelectItem>
                        <SelectItem value="es">{t("languages.es")}</SelectItem>
                        <SelectItem value="fr">{t("languages.fr")}</SelectItem>
                        <SelectItem value="pt">{t("languages.pt")}</SelectItem>
                        <SelectItem value="zh">{t("languages.zh")}</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="darkMode"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-start gap-x-14 col-span-2">
                    <div>
                      <FormLabel>{t("settings.general.darkMode")}</FormLabel>
                      <FormDescription>
                        {t("settings.general.darkModeDescription")}
                      </FormDescription>
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="emailNotifications"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-start gap-x-14 col-span-2">
                    <div>
                      <FormLabel>
                        {t("settings.general.emailNotifications")}
                      </FormLabel>
                      <FormDescription>
                        {t("settings.general.emailNotificationsDescription")}
                      </FormDescription>
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <div className="flex items-center justify-end gap-x-4">
              <Button
                type="button"
                variant="ghost"
                onClick={() => form.reset()}
                disabled={form.formState.isSubmitting}>
                {t("common.cancel")}
              </Button>
              <Button
                type="submit"
                disabled={
                  !form.formState.isDirty || form.formState.isSubmitting
                }>
                {form.formState.isSubmitting
                  ? t("common.saving")
                  : t("common.save")}
              </Button>
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
};

export default GeneralTab;

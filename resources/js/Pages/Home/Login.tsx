import { useState } from "react";
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
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { router } from "@inertiajs/react";
import { useTranslation } from "react-i18next";
import { LanguageSelector } from "@/components/forms/language-selector";
import "@/i18n";
import { useLoginMutation } from "@/services/api/auth";

const formSchema = z.object({
  email: z.string().min(1, "Email is required").email("Invalid email address"),
  password: z
    .string()
    .min(1, "Password is required")
    .min(6, { message: "Password must be at least 6 characters" }),
});

export default function Component() {
  const [isLoading, setIsLoading] = useState(false);
  const { t } = useTranslation();
  const [login] = useLoginMutation();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    setIsLoading(true);
    try {
      // @ts-expect-error
      await login(values).unwrap();
      // Redirect after successful login - you might want to use router.visit here
      router.visit("/dashboard");
    } catch (error: any) {
      console.log(error);
      if (error.data?.message) {
        form.setError("root", { message: error.data.message });
      } else if (error.status === 401) {
        form.setError("root", { message: t("errors.invalidCredentials") });
      } else {
        form.setError("root", { message: t("errors.serverError") });
      }
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="min-h-screen bg-background flex flex-col md:flex-row">
      {/* Side Content */}
      <div className="md:w-2/5 bg-gradient-to-r from-[#2b354c] to-background from-30% to-100% p-8 flex flex-col justify-center">
        <div className="mb-8">
          <img
            src="/placeholder.svg?height=80&width=80"
            alt="Company Logo"
            width={80}
            height={80}
            className="rounded-full bg-white p-2"
          />
        </div>
        <h2 className="text-3xl font-bold text-white mb-4">
          {t("auth.welcomeTitle")}
        </h2>
        <p className="text-white mb-4">{t("auth.welcomeDescription")}</p>
        <ul className="text-white list-disc list-inside">
          <li>{t("auth.feature1")}</li>
          <li>{t("auth.feature2")}</li>
          <li>{t("auth.feature3")}</li>
        </ul>
      </div>

      {/* Login Form */}
      <div className="md:w-1/2 flex items-center justify-center p-8">
        <Card className="w-full max-w-[400px]">
          <CardHeader>
            <div className="flex justify-between items-center">
              <CardTitle>{t("auth.login")}</CardTitle>
              <LanguageSelector />
            </div>
            <CardDescription>{t("auth.enterEmail")}</CardDescription>
          </CardHeader>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <CardContent className="space-y-4">
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("auth.email")}</FormLabel>
                      <FormControl>
                        <Input placeholder={t("auth.enterEmail")} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("auth.password")}</FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder={t("auth.enterPassword")}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                {form.formState.errors.root && (
                  <FormMessage>
                    {form.formState.errors.root.message}
                  </FormMessage>
                )}
              </CardContent>
              <CardFooter>
                <Button className="w-full" type="submit" disabled={isLoading}>
                  {isLoading ? t("auth.loggingIn") : t("auth.logIn")}
                </Button>
              </CardFooter>
            </form>
          </Form>
        </Card>
      </div>
    </div>
  );
}

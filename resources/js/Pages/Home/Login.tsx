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
import { router } from '@inertiajs/react';

const formSchema = z.object({
  email: z.string().email({ message: "Invalid email address" }),
  password: z
    .string()
    .min(6, { message: "Password must be at least 6 characters" }),
});

export default function Component() {
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    setIsLoading(true);
    router.post("/login", values, {
      onSuccess: () => {
        setIsLoading(false);
      },
      onError: (errors) => {
        console.error(errors);
        form.setError("root", { message: "An error occurred during login" });
      },
      onFinish: () => {
        setIsLoading(false);
      }
    });
  }

  return (
    <div className="min-h-screen bg-background flex flex-col md:flex-row">
      {/* Side Content */}
      <div className="md:w-2/5 bg-accent p-8 flex flex-col justify-center">
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
          Welcome to Our Platform
        </h2>
        <p className="text-white mb-4">
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do
          eiusmod tempor incididunt ut labore et dolore magna aliqua.
        </p>
        <ul className="text-white list-disc list-inside">
          <li>Feature 1: Lorem ipsum dolor sit amet</li>
          <li>Feature 2: Consectetur adipiscing elit</li>
          <li>Feature 3: Sed do eiusmod tempor incididunt</li>
        </ul>
      </div>

      {/* Login Form */}
      <div className="md:w-1/2 flex items-center justify-center p-8">
        <Card className="w-full max-w-[400px]">
          <CardHeader>
            <CardTitle>Login</CardTitle>
            <CardDescription>
              Enter your email and password to log in
            </CardDescription>
          </CardHeader>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <CardContent className="space-y-4">
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input placeholder="Enter your email" {...field} />
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
                      <FormLabel>Password</FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder="Enter your password"
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
                  {isLoading ? "Logging in..." : "Log in"}
                </Button>
              </CardFooter>
            </form>
          </Form>
        </Card>
      </div>
    </div>
  );
}

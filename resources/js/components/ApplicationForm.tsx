import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";
import { useDropzone } from "react-dropzone";
import React from "react";

const uploadFormSchema = z.object({
  appName: z.string().min(1, "Application name is required"),
  description: z.string().optional(),
  repoUrl: z
    .string()
    .refine(
      (value) => {
        if (!value) return true; // Allow empty string
        try {
          const url = new URL(value);
          // Check if it's GitHub or GitLab
          if (!["github.com", "gitlab.com"].includes(url.hostname)) {
            return false;
          }
          // Check if it has the pattern: hostname/user/repo
          const parts = url.pathname.split("/").filter(Boolean);
          return parts.length === 2; // Should have exactly user and repo
        } catch {
          return false;
        }
      },
      {
        message:
          "Must be a valid GitHub or GitLab repository URL (e.g., https://github.com/user/repo)",
      }
    )
    .optional(),
  branch: z.string().optional(),
});

interface ApplicationFormProps {
  onSubmit: (
    values: z.infer<typeof uploadFormSchema>,
    file: File | null
  ) => void;
  isUploading: boolean;
  branches: string[];
  isLoadingBranches: boolean;
  onRepoUrlChange: (url: string) => void;
}

export function ApplicationForm({
  onSubmit,
  isUploading,
  branches,
  isLoadingBranches,
  onRepoUrlChange,
}: ApplicationFormProps) {
  const [uploadedFile, setUploadedFile] = React.useState<File | null>(null);

  const form = useForm<z.infer<typeof uploadFormSchema>>({
    resolver: zodResolver(uploadFormSchema),
    defaultValues: {
      appName: "",
      description: "",
      repoUrl: "",
      branch: "",
    },
  });

  const onDrop = React.useCallback((acceptedFiles: File[]) => {
    const file = acceptedFiles[0];
    if (file) {
      setUploadedFile(file);
    }
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      "application/zip": [".zip"],
    },
    maxFiles: 1,
    multiple: undefined,
    onDragEnter: undefined,
    onDragOver: undefined,
    onDragLeave: undefined,
  });

  React.useEffect(() => {
    const subscription = form.watch((value, { name }) => {
      if (name === "repoUrl") {
        onRepoUrlChange(value.repoUrl || "");
      }
    });

    return () => subscription.unsubscribe();
  }, [form, onRepoUrlChange]);

  const handleSubmit = (values: z.infer<typeof uploadFormSchema>) => {
    onSubmit(values, uploadedFile);
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="appName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Application Name</FormLabel>
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
                <Textarea
                  placeholder="Enter application description"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="repoUrl"
          render={({ field }) => (
            <FormItem>
              <FormLabel>GitHub/GitLab Repository URL (Optional)</FormLabel>
              <FormControl>
                <Input
                  placeholder="https://github.com/username/repository"
                  {...field}
                />
              </FormControl>
              <FormDescription>
                Enter a valid GitHub or GitLab repository URL (e.g.,
                https://github.com/user/repo)
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        {form.watch("repoUrl") && (
          <FormField
            control={form.control}
            name="branch"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Branch</FormLabel>
                <Select
                  disabled={isLoadingBranches}
                  value={field.value}
                  onValueChange={field.onChange}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue
                        placeholder={
                          isLoadingBranches
                            ? "Loading branches..."
                            : "Select a branch"
                        }
                      />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {branches.map((branch) => (
                      <SelectItem key={branch} value={branch}>
                        {branch}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
        )}
        <div
          {...getRootProps()}
          className="border-2 border-dashed rounded-lg p-6 text-center cursor-pointer hover:border-primary">
          <input {...getInputProps()} />
          {isDragActive ? (
            <p>Drop the ZIP file here...</p>
          ) : (
            <p>Drag & drop a ZIP file here, or click to select</p>
          )}
        </div>
        <Button type="submit" className="w-full" disabled={isUploading}>
          {isUploading ? "Deploying..." : "Deploy Application"}
        </Button>
      </form>
    </Form>
  );
}

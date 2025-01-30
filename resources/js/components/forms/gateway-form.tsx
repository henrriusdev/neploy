import { useEffect } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Switch } from "@/components/ui/switch"

const formSchema = z.object({
  name: z.string().min(1, "Name is required"),
  path: z.string().min(1, "Path is required"),
  httpMethod: z.enum(["GET", "POST", "PUT", "DELETE"]),
  backendUrl: z.string().url("Must be a valid URL"),
  requiresAuth: z.boolean(),
  rateLimit: z.number().min(0),
})

interface Gateway {
  id: string
  name: string
  path: string
  httpMethod: string
  backendUrl: string
  requiresAuth: boolean
  rateLimit: number
}

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  gateway: Gateway | null
  onSubmit: (data: z.infer<typeof formSchema>) => void
}

export function GatewayForm({ open, onOpenChange, gateway, onSubmit }: Props) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      path: "",
      httpMethod: "GET",
      backendUrl: "",
      requiresAuth: false,
      rateLimit: 60,
    },
  })

  useEffect(() => {
    if (gateway) {
      form.reset(gateway)
    }
  }, [gateway, form])

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>
            {gateway ? "Edit Gateway Route" : "Create Gateway Route"}
          </DialogTitle>
          <DialogDescription>
            Configure the API gateway route settings below.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="User API" {...field} />
                  </FormControl>
                  <FormDescription>
                    A descriptive name for this route
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="path"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Route Path</FormLabel>
                  <FormControl>
                    <Input placeholder="/api/users" {...field} />
                  </FormControl>
                  <FormDescription>
                    The public path for this route
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="httpMethod"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>HTTP Method</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select a method" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="GET">GET</SelectItem>
                      <SelectItem value="POST">POST</SelectItem>
                      <SelectItem value="PUT">PUT</SelectItem>
                      <SelectItem value="DELETE">DELETE</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="backendUrl"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Backend URL</FormLabel>
                  <FormControl>
                    <Input 
                      placeholder="https://api.backend.com/users" 
                      {...field} 
                    />
                  </FormControl>
                  <FormDescription>
                    The backend service URL this route forwards to
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="requiresAuth"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                  <div className="space-y-0.5">
                    <FormLabel className="text-base">
                      Require Authentication
                    </FormLabel>
                    <FormDescription>
                      Enable if this route requires authentication
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="rateLimit"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Rate Limit (requests/minute)</FormLabel>
                  <FormControl>
                    <Input 
                      type="number" 
                      {...field}
                      onChange={e => field.onChange(parseInt(e.target.value))}
                    />
                  </FormControl>
                  <FormDescription>
                    Maximum number of requests allowed per minute
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="submit">
                {gateway ? "Update Route" : "Create Route"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

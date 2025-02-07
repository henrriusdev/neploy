import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Button} from "@/components/ui/button";
import {z} from "zod";
import {zodResolver} from "@hookform/resolvers/zod"
import {useForm} from "react-hook-form";
import {RadioGroup, RadioGroupItem} from "@/components/ui/radio-group";
import {useGetGatewayConfigQuery, useSaveGatewayConfigMutation} from "@/services/api/gateways";
import {useToast} from "@/hooks";
import {GatewayConfigProps} from "@/types";

const formSchema = z.object({
  defaultVersioning: z.enum(["header", "uri"], {required_error: "You need to select a default versioning type."}),
  defaultVersion: z.enum(["latest", "stable"], {required_error: "You need to select a default version."}),
  loadBalancer: z.enum(["true", "false"], {required_error: "You need to choose if you want to use a load balancer or not."})
})
export const GatewayConfig: React.FC<GatewayConfigProps> = ({config}) => {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      defaultVersion: config.defaultVersion,
      defaultVersioning: config.defaultVersioningType,
      loadBalancer: config.loadBalancer ? "true" : "false",
    },
  });
  const {toast} = useToast();

  const getConfig = useGetGatewayConfigQuery(null, {
    skip: true
  });
  const [setConfig] = useSaveGatewayConfigMutation();

  // 2. Define a submit handler.
  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      await setConfig({
        defaultVersion: values.defaultVersion,
        defaultVersioning: values.defaultVersioning,
        loadBalancer: values.loadBalancer === "true",
      }).unwrap();
      toast({
        title: "Config saved!",
        description: "Gateway Config was saved successfully."
      })
    } catch (e) {
      console.log(e)
      toast({
        title: "Config not saved.",
        description: e.data.message
      })
    }
  }

  return (
    <section>
      <h2 className="text-2xl font-semibold mb-5">
        API Gateway Config
      </h2>
      <p className="text-muted-foreground">
        Here you can define your config with a variety of options to use.
      </p>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 grid grid-cols-3 items-center">
          <FormField
            control={form.control}
            name="defaultVersioning"
            render={({field}) => (
              <FormItem>
                <FormLabel>Default versioning</FormLabel>
                <FormControl>
                  <RadioGroup
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                    className="flex flex-col space-y-1"
                  >
                    <FormItem className="flex items-center space-x-3 space-y-0">
                      <FormControl>
                        <RadioGroupItem value="header"/>
                      </FormControl>
                      <FormLabel className="font-normal">
                        Header
                      </FormLabel>
                    </FormItem>
                    <FormItem className="flex items-center space-x-3 space-y-0">
                      <FormControl>
                        <RadioGroupItem value="uri"/>
                      </FormControl>
                      <FormLabel className="font-normal">
                        URI
                      </FormLabel>
                    </FormItem>
                  </RadioGroup>
                </FormControl>
                <FormMessage/>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="defaultVersion"
            render={({field}) => (
              <FormItem>
                <FormLabel>Default version</FormLabel>
                <FormControl>
                  <RadioGroup
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                    className="flex flex-col space-y-1"
                  >
                    <FormItem className="flex items-center space-x-3 space-y-0">
                      <FormControl>
                        <RadioGroupItem value="latest"/>
                      </FormControl>
                      <FormLabel className="font-normal">
                        Latest
                      </FormLabel>
                    </FormItem>
                    <FormItem className="flex items-center space-x-3 space-y-0">
                      <FormControl>
                        <RadioGroupItem value="stable"/>
                      </FormControl>
                      <FormLabel className="font-normal">
                        Stable
                      </FormLabel>
                    </FormItem>
                  </RadioGroup>
                </FormControl>
                <FormMessage/>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="loadBalancer"
            render={({field}) => (
              <FormItem>
                <FormLabel>You need a Load Balancer?</FormLabel>
                <FormControl>
                  <RadioGroup
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                    className="flex flex-col space-y-1"
                  >
                    <FormItem className="flex items-center space-x-3 space-y-0">
                      <FormControl>
                        <RadioGroupItem value="true"/>
                      </FormControl>
                      <FormLabel className="font-normal">
                        Yes
                      </FormLabel>
                    </FormItem>
                    <FormItem className="flex items-center space-x-3 space-y-0">
                      <FormControl>
                        <RadioGroupItem value="false"/>
                      </FormControl>
                      <FormLabel className="font-normal">
                        No
                      </FormLabel>
                    </FormItem>
                  </RadioGroup>
                </FormControl>
                <FormMessage/>
              </FormItem>
            )}
          />
          <div className="col-span-full flex justify-end items-center gap-x-2">
            <Button type="button" variant="ghost">Cancel</Button>
            <Button type="submit" variant="default">Save</Button>
          </div>
        </form>
      </Form>
    </section>
  )
}
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Button} from "@/components/ui/button";
import {z} from "zod";
import {zodResolver} from "@hookform/resolvers/zod"
import {useForm} from "react-hook-form";
import {RadioGroup, RadioGroupItem} from "@/components/ui/radio-group";
import {useGetGatewayConfigQuery, useSaveGatewayConfigMutation} from "@/services/api/gateways";
import {useToast} from "@/hooks";
import {GatewayConfigProps} from "@/types";
import {useTranslation} from "react-i18next";

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
  const {t} = useTranslation();
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
        {t('dashboard.gateways.config')}
      </h2>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 grid grid-cols-3 items-center">
          <FormField
            control={form.control}
            name="defaultVersioning"
            render={({field}) => (
              <FormItem>
                <FormLabel>{t('dashboard.gateways.defaultVersioning')}</FormLabel>
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
                        {t('dashboard.gateways.defaultVersioningHeader')}
                      </FormLabel>
                    </FormItem>
                    <FormItem className="flex items-center space-x-3 space-y-0">
                      <FormControl>
                        <RadioGroupItem value="uri"/>
                      </FormControl>
                      <FormLabel className="font-normal">
                        {t('dashboard.gateways.defaultVersioningURI')}
                      </FormLabel>
                    </FormItem>
                  </RadioGroup>
                </FormControl>
                <FormMessage/>
              </FormItem>
            )}
          />
          <div className="col-span-full flex justify-end items-center gap-x-2">
            <Button type="button" variant="ghost">{t('actions.cancel')}</Button>
            <Button type="submit" variant="default">{t('actions.save')}</Button>
          </div>
        </form>
      </Form>
    </section>
  )
}
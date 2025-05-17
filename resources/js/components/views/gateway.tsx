import {GatewayProps} from "@/types";
import {Globe, Lock} from "lucide-react";
import {Tabs, TabsContent, TabsList, TabsTrigger} from "../ui/tabs";
import {GatewayTable} from "../gateway-table";
import {GatewayConfig} from "@/components/views/gateway-config";
import {useTranslation} from "react-i18next";

export function Gateways({ gateways, config }: GatewayProps) {
  const {t} = useTranslation()
  return (
    <div className="container mx-auto p-6">
      <div className="flex-1 overflow-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <div>
              <h1 className="text-2xl font-bold">{t('dashboard.gateways.title')}</h1>
            </div>
          </div>

          <Tabs defaultValue="routes" className="space-y-4">
            <TabsList className="grid w-full grid-cols-2">

              <TabsTrigger value="routes">
                <Globe className="h-4 w-4 mr-2" />
                {t('dashboard.gateways.routes')}
              </TabsTrigger>
              <TabsTrigger value="config">
                <Lock className="h-4 w-4 mr-2" />
                {t('dashboard.gateways.config')}
              </TabsTrigger>
            </TabsList>

            <TabsContent value="routes">
              {!gateways || gateways?.length === 0 ? (
                <p className="text-muted-foreground">
                  {t('dashboard.gateways.noGateways')}
                </p>
              ) : (
                <GatewayTable
                  gateways={gateways}
                />
              )}
            </TabsContent>

            <TabsContent value="config">
              <div className="rounded-md border p-4">
                <GatewayConfig config={config} />
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}

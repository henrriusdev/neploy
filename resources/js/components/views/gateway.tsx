import {Gateway, GatewayProps} from "@/types";
import {useState} from "react";
import {useToast} from "@/hooks/use-toast";
import {router} from "@inertiajs/react";
import {BarChart3, Globe, Lock} from "lucide-react";
import {Tabs, TabsContent, TabsList, TabsTrigger} from "../ui/tabs";
import {GatewayTable} from "../gateway-table";
import {GatewayConfig} from "@/components/views/gateway-config";

export function Gateways({ gateways, config }: GatewayProps) {
  return (
    <div className="container mx-auto p-6">
      <div className="flex-1 overflow-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <div>
              <h1 className="text-2xl font-bold">API Gateway</h1>
            </div>
          </div>

          <Tabs defaultValue="overview" className="space-y-4">
            <TabsList className="grid w-full grid-cols-3">
              <TabsTrigger value="overview">
                <BarChart3 className="h-4 w-4 mr-2" />
                Overview
              </TabsTrigger>
              <TabsTrigger value="routes">
                <Globe className="h-4 w-4 mr-2" />
                Routes
              </TabsTrigger>
              <TabsTrigger value="config">
                <Lock className="h-4 w-4 mr-2" />
                Config
              </TabsTrigger>
            </TabsList>

            <TabsContent value="overview">
              <div className="rounded-md border p-4">
                <h2 className="text-lg font-semibold mb-2">Gateway Overview</h2>
                <p className="text-muted-foreground">
                  Gateway statistics and metrics will be displayed here.
                </p>
              </div>
            </TabsContent>

            <TabsContent value="routes">
              {!gateways || gateways?.length === 0 ? (
                <p className="text-muted-foreground">
                  No routes configured for this application.
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

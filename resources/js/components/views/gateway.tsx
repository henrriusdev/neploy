import { Gateway, GatewayProps } from "@/types";
import { useState } from "react";
import { useToast } from "@/hooks/use-toast";
import { router } from "@inertiajs/react";
import { Button } from "../ui/button";
import { Activity, BarChart3, Globe, Lock, PlusCircle } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { GatewayTable } from "../gateway-table";
import { GatewayForm } from "../forms";
import {GatewayConfig} from "@/components/views/gateway-config";

export function Gateways({ gateways, config }: GatewayProps) {
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingGateway, setEditingGateway] = useState<Gateway | null>(null);
  const { toast } = useToast();

  const handleCreate = (data: Partial<Gateway>) => {
    const { application, ...formData } = data;
    router.post("/gateways", formData, {
      onSuccess: () => {
        setIsFormOpen(false);
        toast({
          title: "Success",
          description: "Gateway route created successfully",
        });
      },
      onError: () => {
        toast({
          title: "Error",
          description: "Failed to create gateway route",
          variant: "destructive",
        });
      },
    });
  };

  const handleUpdate = (id: string, data: Partial<Gateway>) => {
    const { application, ...formData } = data;
    router.put(`/gateways/${id}`, formData, {
      onSuccess: () => {
        setEditingGateway(null);
        toast({
          title: "Success",
          description: "Gateway route updated successfully",
        });
      },
      onError: () => {
        toast({
          title: "Error",
          description: "Failed to update gateway route",
          variant: "destructive",
        });
      },
    });
  };

  const handleDelete = (id: string) => {
    router.delete(`/gateways/${id}`, {
      onSuccess: () => {
        toast({
          title: "Success",
          description: "Gateway route deleted successfully",
        });
      },
      onError: () => {
        toast({
          title: "Error",
          description: "Failed to delete gateway route",
          variant: "destructive",
        });
      },
    });
  };

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
                  onEdit={setEditingGateway}
                  onDelete={handleDelete}
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

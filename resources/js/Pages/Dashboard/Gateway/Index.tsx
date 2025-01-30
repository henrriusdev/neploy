import { useState } from "react";
import { PlusCircle } from "lucide-react";
import { router } from "@inertiajs/react";
import { Globe, Lock, Activity, BarChart3 } from "lucide-react";

import { DashboardLayout } from "@/components/Layouts/DashboardLayout";
import { Button } from "@/components/ui/button";
import { useToast } from "@/hooks/use-toast";
import { GatewayTable } from "./components/gateway-table";
import { GatewayForm } from "./components/gateway-form";
import { Gateway } from "@/types/common";
import { GatewayProps } from "@/types/props";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

export default function Index({ gateways, application }: GatewayProps) {
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingGateway, setEditingGateway] = useState<Gateway | null>(null);
  const { toast } = useToast();

  const handleCreate = (data: Partial<Gateway>) => {
    router.post("/gateways", data, {
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
    router.put(`/gateways/${id}`, data, {
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
    <div className="flex h-screen overflow-hidden">
      <div className="flex-1 overflow-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <div>
              <h1 className="text-2xl font-bold">API Gateway Configuration</h1>
              {application && (
                <p className="text-muted-foreground">
                  Application: {application.name}
                </p>
              )}
            </div>
            <Button onClick={() => setIsFormOpen(true)}>
              <PlusCircle className="mr-2 h-4 w-4" />
              Add Route
            </Button>
          </div>

          <Tabs defaultValue="overview" className="space-y-4">
            <TabsList className="grid w-full grid-cols-4">
              <TabsTrigger value="overview">
                <BarChart3 className="h-4 w-4 mr-2" />
                Overview
              </TabsTrigger>
              <TabsTrigger value="routes">
                <Globe className="h-4 w-4 mr-2" />
                Routes
              </TabsTrigger>
              <TabsTrigger value="security">
                <Lock className="h-4 w-4 mr-2" />
                Security
              </TabsTrigger>
              <TabsTrigger value="rate-limiting">
                <Activity className="h-4 w-4 mr-2" />
                Rate Limiting
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

            <TabsContent value="security">
              <div className="rounded-md border p-4">
                <h2 className="text-lg font-semibold mb-2">
                  Security Settings
                </h2>
                <p className="text-muted-foreground">
                  Configure security settings for your gateway here.
                </p>
              </div>
            </TabsContent>

            <TabsContent value="rate-limiting">
              <div className="rounded-md border p-4">
                <h2 className="text-lg font-semibold mb-2">Rate Limiting</h2>
                <p className="text-muted-foreground">
                  Configure rate limiting rules for your gateway endpoints.
                </p>
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </div>

      <GatewayForm
        open={isFormOpen || !!editingGateway}
        onOpenChange={(open) => {
          setIsFormOpen(open);
          if (!open) setEditingGateway(null);
        }}
        gateway={editingGateway}
        onSubmit={
          editingGateway
            ? (data) => handleUpdate(editingGateway.id, data)
            : handleCreate
        }
      />
    </div>
  );
}

Index.layout = (page: any) => {
  return <DashboardLayout>{page}</DashboardLayout>;
};

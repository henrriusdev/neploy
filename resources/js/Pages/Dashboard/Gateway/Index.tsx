import { useState } from "react";
import { router } from "@inertiajs/react";
import {
  Globe,
  Lock,
  Activity,
  BarChart3,
  PlusCircle,
  Settings,
} from "lucide-react";

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
            <div className="space-x-2">
              <Button onClick={() => setIsFormOpen(true)}>
                <PlusCircle className="mr-2 h-4 w-4" />
                Add Route
              </Button>
            </div>
          </div>

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
  const { user, teamName, logoUrl } = page.props;
  return (
    <DashboardLayout user={user} teamName={teamName} logoUrl={logoUrl}>
      {page}
    </DashboardLayout>
  );
};

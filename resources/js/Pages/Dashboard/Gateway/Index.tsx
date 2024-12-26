import { useState } from "react"
import { PlusCircle } from "lucide-react"
import { router } from "@inertiajs/react"

import { DashboardLayout } from "@/components/Layouts/DashboardLayout"
import { Button } from "@/components/ui/button"
import { useToast } from "@/hooks/use-toast"
import { GatewayTable } from "./components/gateway-table"
import { GatewayForm } from "./components/gateway-form"
import { GatewaySidebar } from "./components/gateway-sidebar"

interface Gateway {
  id: string
  name: string
  path: string
  httpMethod: string
  backendUrl: string
  requiresAuth: boolean
  rateLimit: number
  applicationId: string
  application: {
    id: string
    name: string
  }
}

interface Props {
  gateways: Gateway[]
  application?: {
    id: string
    name: string
  }
  user: {
    name: string
    email: string
    username: string
    provider: string
  }
  teamName: string
  logoUrl: string
}

export default function Index({ gateways, application }: Props) {
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingGateway, setEditingGateway] = useState<Gateway | null>(null)
  const { toast } = useToast()

  const handleCreate = (data: Partial<Gateway>) => {
    router.post("/gateways", data, {
      onSuccess: () => {
        setIsFormOpen(false)
        toast({
          title: "Success",
          description: "Gateway route created successfully",
        })
      },
      onError: () => {
        toast({
          title: "Error",
          description: "Failed to create gateway route",
          variant: "destructive",
        })
      },
    })
  }

  const handleUpdate = (id: string, data: Partial<Gateway>) => {
    router.put(`/gateways/${id}`, data, {
      onSuccess: () => {
        setEditingGateway(null)
        toast({
          title: "Success",
          description: "Gateway route updated successfully",
        })
      },
      onError: () => {
        toast({
          title: "Error",
          description: "Failed to update gateway route",
          variant: "destructive",
        })
      },
    })
  }

  const handleDelete = (id: string) => {
    router.delete(`/gateways/${id}`, {
      onSuccess: () => {
        toast({
          title: "Success",
          description: "Gateway route deleted successfully",
        })
      },
      onError: () => {
        toast({
          title: "Error",
          description: "Failed to delete gateway route",
          variant: "destructive",
        })
      },
    })
  }

  return (
    <div className="flex h-screen overflow-hidden">
      <GatewaySidebar />
      
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
          setIsFormOpen(open)
          if (!open) setEditingGateway(null)
        }}
        gateway={editingGateway}
        onSubmit={editingGateway ? 
          (data) => handleUpdate(editingGateway.id, data) : 
          handleCreate
        }
      />
    </div>
  )
}

Index.layout = (page: any) => {
  return (
    <DashboardLayout>
      {page}
    </DashboardLayout>
  );
};

import { Link } from "@inertiajs/react"
import { Edit2, Trash2 } from "lucide-react"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { GatewayTableProps } from "@/types/props"

const methodColors = {
  GET: "bg-green-100 text-green-800",
  POST: "bg-blue-100 text-blue-800",
  PUT: "bg-yellow-100 text-yellow-800",
  DELETE: "bg-red-100 text-red-800",
} as const

export function GatewayTable({ gateways, onEdit, onDelete }: GatewayTableProps) {
  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Route Path</TableHead>
            <TableHead>Method</TableHead>
            <TableHead>Backend URL</TableHead>
            <TableHead>Auth Required</TableHead>
            <TableHead>Rate Limit</TableHead>
            <TableHead>Application</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {gateways.map((gateway) => (
            <TableRow key={gateway.id}>
              <TableCell>{gateway.path}</TableCell>
              <TableCell>
                <Badge 
                  variant="outline"
                  className={methodColors[gateway.httpMethod as keyof typeof methodColors]}
                >
                  {gateway.httpMethod}
                </Badge>
              </TableCell>
              <TableCell className="font-mono text-sm">
                {gateway.backendUrl}
              </TableCell>
              <TableCell>
                <Badge variant={gateway.requiresAuth ? "default" : "secondary"}>
                  {gateway.requiresAuth ? "Yes" : "No"}
                </Badge>
              </TableCell>
              <TableCell>{gateway.rateLimit} req/min</TableCell>
              <TableCell>
                <Link 
                  href={`/dashboard/applications/${gateway.applicationId}`}
                  className="text-primary hover:underline"
                >
                  {gateway.application.name}
                </Link>
              </TableCell>
              <TableCell className="text-right space-x-2">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => onEdit(gateway)}
                >
                  <Edit2 className="h-4 w-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => {
                    if (confirm("Are you sure you want to delete this route?")) {
                      onDelete(gateway.id)
                    }
                  }}
                >
                  <Trash2 className="h-4 w-4 text-destructive" />
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

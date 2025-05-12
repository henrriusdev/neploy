import {Link} from "@inertiajs/react";

import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from "@/components/ui/table";
import {GatewayTableProps} from "@/types/props";

export function GatewayTable({
  gateways,
}: GatewayTableProps) {
  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Route Path</TableHead>
            <TableHead>Backend URL</TableHead>
            <TableHead>Application</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {gateways.map((gateway) => (
            <TableRow key={gateway.id}>
              <TableCell>{gateway.path}</TableCell>
              <TableCell className="font-mono text-sm">
                {gateway.backendUrl}
              </TableCell>
              <TableCell>
                <Link
                  href={`/dashboard/applications/${gateway.applicationId}`}
                  className="text-primary hover:underline">
                  {gateway.application.appName}
                </Link>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const TraceabilityTab = () => {
  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>System Activity</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-[200px]">
            {/* <LineChart 
              data={{
                labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
                datasets: [{
                  label: 'System Events',
                  data: [65, 59, 80, 81, 56, 55],
                }]
              }}
            /> */}
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex justify-between items-center">
            Audit Logs
            <div className="flex space-x-2">
              <Input placeholder="Search logs..." className="w-[300px]" />
              <Button variant="outline">Filter</Button>
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Timestamp</TableHead>
                <TableHead>User</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Resource</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell>{new Date().toLocaleString()}</TableCell>
                <TableCell>admin@example.com</TableCell>
                <TableCell>Updated API endpoint</TableCell>
                <TableCell>/api/v1/users</TableCell>
                <TableCell>Success</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>{new Date().toLocaleString()}</TableCell>
                <TableCell>user@example.com</TableCell>
                <TableCell>Role modified</TableCell>
                <TableCell>User Management</TableCell>
                <TableCell>Success</TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default TraceabilityTab;

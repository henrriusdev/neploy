import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';

const ApiGatewayTab = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex justify-between items-center">
          API Gateway Configuration
          <Dialog>
            <DialogTrigger asChild>
              <Button>Add Endpoint</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New Endpoint</DialogTitle>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <Input placeholder="Endpoint Path" />
                <Input placeholder="Target URL" />
                <div className="flex items-center justify-between">
                  <span className="text-sm">Enable Authentication</span>
                  <Switch />
                </div>
              </div>
            </DialogContent>
          </Dialog>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-6">
          <div className="grid grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-sm">Authentication</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm">JWT Authentication</span>
                    <Switch />
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm">API Key Authentication</span>
                    <Switch />
                  </div>
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle className="text-sm">Rate Limiting</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <Input placeholder="Requests per minute" type="number" />
                  <div className="flex items-center justify-between">
                    <span className="text-sm">Enable Rate Limiting</span>
                    <Switch />
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Path</TableHead>
                <TableHead>Target</TableHead>
                <TableHead>Auth Required</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell>/api/v1/users</TableCell>
                <TableCell>http://users-service/v1</TableCell>
                <TableCell>Yes</TableCell>
                <TableCell>
                  <Switch defaultChecked />
                </TableCell>
                <TableCell>
                  <Button variant="ghost" size="sm">Edit</Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
};

export default ApiGatewayTab;

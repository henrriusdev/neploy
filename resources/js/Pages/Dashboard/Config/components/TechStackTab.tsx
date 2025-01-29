import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';

const TechStackTab = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex justify-between items-center">
          Tech Stack Management
          <Dialog>
            <DialogTrigger asChild>
              <Button>Add Technology</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New Technology</DialogTitle>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <Input placeholder="Technology Name" />
                <Input placeholder="Version" />
                <Input placeholder="Category (Framework, Language, Tool)" />
              </div>
            </DialogContent>
          </Dialog>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Technology</TableHead>
              <TableHead>Version</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>React</TableCell>
              <TableCell>18.2.0</TableCell>
              <TableCell>Framework</TableCell>
              <TableCell>
                <Badge>Active</Badge>
              </TableCell>
              <TableCell>
                <Button variant="ghost" size="sm">Edit</Button>
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>TypeScript</TableCell>
              <TableCell>5.0.0</TableCell>
              <TableCell>Language</TableCell>
              <TableCell>
                <Badge>Active</Badge>
              </TableCell>
              <TableCell>
                <Button variant="ghost" size="sm">Edit</Button>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};

export default TechStackTab;

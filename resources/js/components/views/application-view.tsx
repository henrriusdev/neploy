"use client"

import { useState } from "react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Progress } from "@/components/ui/progress"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { Edit, StopCircle, Trash2, ChevronDown, Search, Plus } from "lucide-react"

export function ApplicationView() {
  const [isLogsOpen, setIsLogsOpen] = useState(true)
  const [searchLogs, setSearchLogs] = useState("")

  return (
    <div className="p-6 space-y-6 max-w-7xl mx-auto">
      {/* Header Section */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-bold text-white">Hangman Game</h1>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30">
              Running
            </Badge>
            <span className="text-sm text-muted-foreground">ID: app_12345</span>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm">
            <Edit className="w-4 h-4 mr-2" />
            Edit
          </Button>
          <Button variant="outline" size="sm" className="text-yellow-400 border-yellow-400/30 hover:bg-yellow-400/10">
            <StopCircle className="w-4 h-4 mr-2" />
            Stop
          </Button>
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="outline" size="sm" className="text-red-400 border-red-400/30 hover:bg-red-400/10">
                <Trash2 className="w-4 h-4 mr-2" />
                Delete
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                <AlertDialogDescription>
                  This action cannot be undone. This will permanently delete your application and remove all associated
                  data.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction className="bg-red-600 hover:bg-red-700">Delete</AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Overview Section */}
        <Card className="bg-[#0F1631] border-border/50">
          <CardHeader>
            <CardTitle>Overview</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Created At</p>
                <p className="text-sm font-medium">Feb 14, 2025</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Updated At</p>
                <p className="text-sm font-medium">Feb 14, 2025</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Owner</p>
                <p className="text-sm font-medium">user@example.com</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Region</p>
                <p className="text-sm font-medium">us-east-1</p>
              </div>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Description</p>
              <p className="text-sm">A simple hangman game built with Go and React.</p>
            </div>
          </CardContent>
        </Card>

        {/* Metrics Section */}
        <Card className="bg-[#0F1631] border-border/50">
          <CardHeader>
            <CardTitle>Metrics</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm">CPU Usage</span>
                <span className="text-sm font-medium">45%</span>
              </div>
              <Progress value={45} className="h-2" />
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm">Memory Usage</span>
                <span className="text-sm font-medium">1.2GB / 2GB</span>
              </div>
              <Progress value={60} className="h-2" />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Uptime</p>
                <p className="text-sm font-medium">5d 12h 30m</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Requests/min</p>
                <p className="text-sm font-medium">250</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Logs Section */}
        <Card className="md:col-span-2 bg-[#0F1631] border-border/50">
          <Collapsible open={isLogsOpen} onOpenChange={setIsLogsOpen}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle>Logs</CardTitle>
              <CollapsibleTrigger asChild>
                <Button variant="ghost" size="sm">
                  <ChevronDown className={`h-4 w-4 transition-transform ${isLogsOpen ? "transform rotate-180" : ""}`} />
                </Button>
              </CollapsibleTrigger>
            </CardHeader>
            <CollapsibleContent>
              <CardContent>
                <div className="flex items-center gap-2 mb-4">
                  <Search className="w-4 h-4 text-muted-foreground" />
                  <Input
                    placeholder="Search logs..."
                    value={searchLogs}
                    onChange={(e) => setSearchLogs(e.target.value)}
                    className="max-w-sm"
                  />
                </div>
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Timestamp</TableHead>
                        <TableHead>Level</TableHead>
                        <TableHead className="w-full">Message</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      <TableRow>
                        <TableCell className="text-xs">2025-02-14 13:45:30</TableCell>
                        <TableCell>
                          <Badge variant="outline">INFO</Badge>
                        </TableCell>
                        <TableCell className="text-sm">Application started successfully</TableCell>
                      </TableRow>
                      <TableRow>
                        <TableCell className="text-xs">2025-02-14 13:45:29</TableCell>
                        <TableCell>
                          <Badge variant="outline">DEBUG</Badge>
                        </TableCell>
                        <TableCell className="text-sm">Initializing database connection</TableCell>
                      </TableRow>
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </CollapsibleContent>
          </Collapsible>
        </Card>

        {/* Settings Section */}
        <Card className="md:col-span-2 bg-[#0F1631] border-border/50">
          <CardHeader>
            <CardTitle>Settings</CardTitle>
          </CardHeader>
          <CardContent>
            <Tabs defaultValue="general">
              <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger value="general">General</TabsTrigger>
                <TabsTrigger value="webhooks">Webhooks</TabsTrigger>
                <TabsTrigger value="resources">Resources</TabsTrigger>
              </TabsList>
              <TabsContent value="general" className="space-y-4">
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h3 className="text-sm font-medium">Environment Variables</h3>
                    <Button variant="outline" size="sm">
                      <Plus className="w-4 h-4 mr-2" />
                      Add Variable
                    </Button>
                  </div>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Key</TableHead>
                        <TableHead>Value</TableHead>
                        <TableHead className="w-24">Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      <TableRow>
                        <TableCell>DATABASE_URL</TableCell>
                        <TableCell>postgres://user:pass@localhost:5432/db</TableCell>
                        <TableCell>
                          <Button variant="ghost" size="icon" className="h-8 w-8">
                            <Edit className="h-4 w-4" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    </TableBody>
                  </Table>
                </div>
              </TabsContent>
              <TabsContent value="webhooks" className="space-y-4">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-medium">Webhook Endpoints</h3>
                  <Button variant="outline" size="sm">
                    <Plus className="w-4 h-4 mr-2" />
                    Add Webhook
                  </Button>
                </div>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>URL</TableHead>
                      <TableHead>Events</TableHead>
                      <TableHead className="w-24">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow>
                      <TableCell>https://api.example.com/webhook</TableCell>
                      <TableCell>deployment.success, deployment.failure</TableCell>
                      <TableCell>
                        <Button variant="ghost" size="icon" className="h-8 w-8">
                          <Edit className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </TabsContent>
              <TabsContent value="resources" className="space-y-4">
                <div className="space-y-4">
                  <div>
                    <h3 className="text-sm font-medium mb-2">Resource Allocation</h3>
                    <div className="space-y-4">
                      <div className="space-y-2">
                        <div className="flex items-center justify-between">
                          <span className="text-sm">CPU Limit</span>
                          <span className="text-sm font-medium">2 cores</span>
                        </div>
                        <Progress value={50} className="h-2" />
                      </div>
                      <div className="space-y-2">
                        <div className="flex items-center justify-between">
                          <span className="text-sm">Memory Limit</span>
                          <span className="text-sm font-medium">2GB</span>
                        </div>
                        <Progress value={60} className="h-2" />
                      </div>
                    </div>
                  </div>
                </div>
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}


"use client"

import {useEffect, useState} from "react"
import {Badge} from "@/components/ui/badge"
import {Button} from "@/components/ui/button"
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card"
import {Input} from "@/components/ui/input"
import {Progress} from "@/components/ui/progress"
import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs"
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table"
import {Collapsible, CollapsibleContent, CollapsibleTrigger} from "@/components/ui/collapsible"
import {ChevronDown, Edit, Plus, Search, Trash2} from "lucide-react"
import {ApplicationProps} from "@/types";

export const ApplicationView: React.FC<ApplicationProps> = ({application}) => {
  const [isLogsOpen, setIsLogsOpen] = useState(true)
  const [searchLogs, setSearchLogs] = useState("")

  const [filteredLogs, setFilteredLogs] = useState(application.logs?.slice(0, 10) ?? []);

  useEffect(() => {
    setFilteredLogs(
      application.logs?.filter((item) =>
        item.toLowerCase().includes(searchLogs.toLowerCase())
      )?.slice(0, 10) ?? []
    );
  }, [searchLogs, application.logs]);

  return (
    <div className="p-6 space-y-6 max-w-7xl mx-auto">
      {/* Header Section */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-bold ">{application.appName}</h1>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30">
              Running
            </Badge>
            <span className="text-sm text-muted-foreground">ID: {application.id}</span>
          </div>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Overview Section */}
        <Card className="border-border/50">
          <CardHeader>
            <CardTitle>Overview</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Created At</p>
                <p className="text-sm font-medium">{application.createdAt}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Updated At</p>
                <p className="text-sm font-medium">{application.updatedAt}</p>
              </div>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Description</p>
              <p className="text-sm">{application.description}</p>
            </div>
          </CardContent>
        </Card>

        {/* Metrics Section */}
        <Card className=" border-border/50">
          <CardHeader>
            <CardTitle>Metrics</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm">CPU Usage</span>
                <span className="text-sm font-medium">{application.cpuUsage.toFixed(2)}%</span>
              </div>
              <Progress value={application.cpuUsage} className="h-2" />
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm">Memory Usage</span>
                <span className="text-sm font-medium">{application.memoryUsage.toFixed(2)}%</span>
              </div>
              <Progress value={application.memoryUsage} className="h-2" />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Uptime</p>
                <p className="text-sm font-medium">{application.uptime}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Requests/min</p>
                <p className="text-sm font-medium">{application.requestsPerMin}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Logs Section */}
        <Card className="md:col-span-2  border-border/50">
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
                        <TableHead>Log</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {filteredLogs && filteredLogs.map((item, i) => (
                        <TableRow key={i}>
                          <TableCell className="text-sm font-mono font-bold">{item}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </CollapsibleContent>
          </Collapsible>
        </Card>

        {/* API Versions Section */}
        <Card className="md:col-span-2  border-border/50">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>API Versions</CardTitle>
              <Button variant="outline" size="sm">
                <Plus className="w-4 h-4 mr-2" />
                New Version
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Version</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created At</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {application.versions?.length ? (
                  application.versions.map((version, i) => (
                    <TableRow key={i}>
                      <TableCell className="font-mono">{version.versionTag}</TableCell>
                      <TableCell>{version.description}</TableCell>
                      <TableCell>
                        <Badge variant="outline" className="text-xs capitalize">{version.status}</Badge>
                      </TableCell>
                      <TableCell>{new Date(version.createdAt).toLocaleDateString()}</TableCell>
                      <TableCell className="text-right">
                        <Button variant="ghost" size="icon" className="h-8 w-8 text-red-400 hover:bg-red-400/10">
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center text-muted-foreground">
                      No versions found.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}


import * as React from "react";
import { DashboardLayout } from "@/components/Layouts/DashboardLayout";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Settings, UsersRound, Boxes, GitBranch } from "lucide-react";

import GeneralTab from "./components/GeneralTab";
import RolesTab from "./components/RolesTab";
import TechStackTab from "./components/TechStackTab";
import TraceabilityTab from "./components/TraceabilityTab";
import { SettingsProps } from "@/types/props";

const Config = ({
  user,
  teamName,
  logoUrl,
  language,
  roles = [],
  techStacks = [],
}: SettingsProps) => {
  return (
    <div className="container mx-auto p-6">
      <h1 className="text-2xl font-bold mb-6">System Configuration</h1>

      <Tabs defaultValue="general" className="space-y-4">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="general">
            <Settings className="h-4 w-4 mr-2" />
            General
          </TabsTrigger>
          <TabsTrigger value="roles">
            <UsersRound className="h-4 w-4 mr-2" />
            Roles
          </TabsTrigger>
          <TabsTrigger value="techstack">
            <Boxes className="h-4 w-4 mr-2" />
            Tech Stack
          </TabsTrigger>
          <TabsTrigger value="traceability">
            <GitBranch className="h-4 w-4 mr-2" />
            Traceability
          </TabsTrigger>
        </TabsList>

        <TabsContent value="general">
          <GeneralTab
            teamName={teamName}
            logoUrl={logoUrl}
            language={language}
          />
        </TabsContent>

        <TabsContent value="roles">
          <RolesTab roles={roles} />
        </TabsContent>

        <TabsContent value="techstack">
          <TechStackTab techStacks={techStacks} />
        </TabsContent>

        <TabsContent value="traceability">
          <TraceabilityTab />
        </TabsContent>
      </Tabs>
    </div>
  );
};

Config.layout = (page: any) => {
  const { user, teamName, logoUrl } = page.props;
  return (
    <DashboardLayout user={user} teamName={teamName} logoUrl={logoUrl}>
      {page}
    </DashboardLayout>
  );
};

export default Config;

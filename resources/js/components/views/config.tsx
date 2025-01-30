import { SettingsProps } from "@/types";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { Boxes, GitBranch, Settings, UsersRound } from "lucide-react";
import GeneralTab from "./config-general";
import RolesTab from "./config-roles";
import TechStackTab from "./config-techstack";
import TraceabilityTab from "./config-traceability";

export const Config = ({
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

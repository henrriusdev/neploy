import React from "react";
import { DashboardLayout } from "@/components/Layouts/DashboardLayout";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import GeneralTab from "./components/GeneralTab";
import RolesTab from "./components/RolesTab";
import TechStackTab from "./components/TechStackTab";
import ApiGatewayTab from "./components/ApiGatewayTab";
import TraceabilityTab from "./components/TraceabilityTab";

const Config = () => {
  return (
    <DashboardLayout>
      <div className="container mx-auto p-6">
        <h1 className="text-2xl font-bold mb-6">System Configuration</h1>

        <Tabs defaultValue="general" className="space-y-4">
          <TabsList className="grid w-full grid-cols-5">
            <TabsTrigger value="general">General</TabsTrigger>
            <TabsTrigger value="roles">Roles</TabsTrigger>
            <TabsTrigger value="techstack">Tech Stack</TabsTrigger>
            <TabsTrigger value="apigateway">API Gateway</TabsTrigger>
            <TabsTrigger value="traceability">Traceability</TabsTrigger>
          </TabsList>

          <TabsContent value="general">
            <GeneralTab />
          </TabsContent>

          <TabsContent value="roles">
            <RolesTab />
          </TabsContent>

          <TabsContent value="techstack">
            <TechStackTab />
          </TabsContent>

          <TabsContent value="apigateway">
            <ApiGatewayTab />
          </TabsContent>

          <TabsContent value="traceability">
            <TraceabilityTab />
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  );
};

export default Config;

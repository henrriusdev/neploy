import * as React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Switch } from "@/components/ui/switch";
import { Input } from "@/components/ui/input";
import { GeneralSettingsProps } from "@/types/props";

const GeneralTab: React.FC<GeneralSettingsProps> = ({
  teamName,
  logoUrl,
  language,
}) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>General Settings</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Team Name</label>
            <Input placeholder="Team Name" value={teamName} />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Default Language</label>
            <Input placeholder="en-US" value={language} />
          </div>

          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-medium">Dark Mode</h3>
              <p className="text-sm text-gray-500">
                Enable dark mode by default
              </p>
            </div>
            <Switch />
          </div>

          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-medium">Email Notifications</h3>
              <p className="text-sm text-gray-500">
                Enable email notifications
              </p>
            </div>
            <Switch />
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default GeneralTab;

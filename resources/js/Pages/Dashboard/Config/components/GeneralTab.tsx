import * as React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Switch } from "@/components/ui/switch";
import { Input } from "@/components/ui/input";
import { GeneralSettingsProps } from "@/types/props";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";

const GeneralTab: React.FC<GeneralSettingsProps> = ({
  teamName: originalTeamName,
  logoUrl: originalLogoUrl,
  language: originalLanguage,
}) => {
  const [teamName, setTeamName] = React.useState(originalTeamName);
  const [logoUrl, setLogoUrl] = React.useState(originalLogoUrl);
  const [language, setLanguage] = React.useState(originalLanguage);
  const [darkMode, setDarkMode] = React.useState(false);
  const [emailNotifications, setEmailNotifications] = React.useState(false);
  const [isEdited, setIsEdited] = React.useState(false);

  const handleTeamNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setTeamName(event.target.value);
    if (event.target.value !== originalTeamName) {
      setIsEdited(true);
    } else {
      setIsEdited(false);
    }
  };

  const handleLogoUrlChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setLogoUrl(event.target.value);
    if (event.target.value !== originalLogoUrl) {
      setIsEdited(true);
    } else {
      setIsEdited(false);
    }
  };

  const handleLanguageChange = (value: string) => {
    setLanguage(value);
    if (value !== originalLanguage) {
      setIsEdited(true);
    } else {
      setIsEdited(false);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>General Settings</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 grid-cols-6">
          <div className="space-y-2 col-span-3">
            <label className="text-sm font-medium">Team Name</label>
            <Input
              placeholder="Team Name"
              value={teamName}
              onChange={handleTeamNameChange}
            />
          </div>
          <div className="space-y-2 col-span-3">
            <label className="text-sm font-medium">Logo URL</label>
            <Input
              placeholder="Logo URL"
              value={logoUrl}
              onChange={handleLogoUrlChange}
            />
          </div>

          <div className="space-y-2 col-span-2">
            <label className="text-sm font-medium">Default Language</label>
            <Select onValueChange={handleLanguageChange} value={language}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Select a language" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="en">English</SelectItem>
                <SelectItem value="es">Spanish</SelectItem>
                <SelectItem value="fr">French</SelectItem>
                <SelectItem value="pt">Portuguese</SelectItem>
                <SelectItem value="zh">Chinese</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-center justify-start gap-x-14 col-span-2">
            <div>
              <h3 className="font-medium">Dark Mode</h3>
              <p className="text-sm text-gray-500">
                Enable dark mode by default
              </p>
            </div>
            <Switch />
          </div>

          <div className="flex items-center justify-start gap-x-14 col-span-2">
            <div>
              <h3 className="font-medium">Email Notifications</h3>
              <p className="text-sm text-gray-500">
                Enable email notifications
              </p>
            </div>
            <Switch />
          </div>
        </div>
        {isEdited && (
          <div className="flex items-center justify-end gap-x-4">
            <Button variant="ghost">Cancel</Button>
            <Button>Save</Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default GeneralTab;

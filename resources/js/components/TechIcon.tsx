import * as React from "react";
import { IconContext } from "react-icons";
import * as Si from "react-icons/si"; // Simple Icons
import * as Di from "react-icons/di"; // Devicons

interface TechIconProps {
  name: string;
  size?: number;
}

export function TechIcon({ name, size = 75 }: TechIconProps) {
  // Convert name to PascalCase and try both Si and Di collections
  const iconName = name
    .split(/[-_\s]+/)
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join("");
  
  const SiIcon = (Si as any)[`Si${iconName}`];
  const DiIcon = (Di as any)[`Di${iconName}`];
  const Icon = SiIcon || DiIcon;

  console.log(`Trying icon: ${iconName}`);
  if (!Icon) {
    console.warn(`Icon not found for: ${name}`);
    return <div className="p-2" style={{ width: size, height: size }} />;
  }

  return (
    <IconContext.Provider value={{ size: `${size}px` }}>
      <div className="p-2">
        <Icon />
      </div>
    </IconContext.Provider>
  );
}

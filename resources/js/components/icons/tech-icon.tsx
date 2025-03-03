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
  return (
    <img className={`w-[${size}] h-${size}`} src={`https://cdn.simpleicons.org/${name.toLowerCase()}`} />
  )
}

import * as LucideIcons from "lucide-react";
import * as RadixIcons from "@radix-ui/react-icons";

// Get all Lucide icons
const lucideIconNames = Object.keys(LucideIcons).filter(
  (key) => typeof LucideIcons[key as keyof typeof LucideIcons] === "function"
);

// Get all Radix icons and remove the 'Icon' suffix
const radixIconNames = Object.keys(RadixIcons)
  .filter((key) => typeof RadixIcons[key as keyof typeof RadixIcons] === "function")
  .map((name) => name.replace(/Icon$/, ""));

// Combine both sets of icons and sort alphabetically
export const icons = [...new Set([...lucideIconNames, ...radixIconNames])].sort();

// Create a type union of all icon names
export type Icon = typeof icons[number];
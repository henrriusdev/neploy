import { icons as lucideIconsList } from 'lucide-react';
import * as RadixIcons from "@radix-ui/react-icons";
import { DiJavascript1, DiPython, DiGo, DiJava, DiReact, DiAngularSimple, DiNodejsSmall, DiDjango, DiDocker, DiGit, DiMysql, DiPostgresql, DiMongodb, DiNginx, DiRedis } from 'react-icons/di';

// Get Lucide icons
const lucideIconNames = Object.keys(lucideIconsList);

// Get Radix icons and remove the 'Icon' suffix
const radixIconNames = Object.keys(RadixIcons)
  .filter((key) => typeof RadixIcons[key as keyof typeof RadixIcons] === "function")
  .map((name) => name.replace(/Icon$/, ""));

// Combine both sets of icons and sort alphabetically
export const icons = [...new Set([...lucideIconNames, ...radixIconNames])].sort();

// Create a type union of all icon names
export type Icon = typeof icons[number];

export const techIcons = [
  { name: 'JavaScript', icon: DiJavascript1 },
  { name: 'Python', icon: DiPython },
  { name: 'Go', icon: DiGo },
  { name: 'Java', icon: DiJava },
  { name: 'React', icon: DiReact },
  { name: 'Angular', icon: DiAngularSimple },
  { name: 'Node.js', icon: DiNodejsSmall },
  { name: 'Django', icon: DiDjango },
  { name: 'Docker', icon: DiDocker },
  { name: 'Git', icon: DiGit },
  { name: 'MySQL', icon: DiMysql },
  { name: 'PostgreSQL', icon: DiPostgresql },
  { name: 'MongoDB', icon: DiMongodb },
  { name: 'Nginx', icon: DiNginx },
  { name: 'Redis', icon: DiRedis },
] as const;
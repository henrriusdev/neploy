import * as React from "react";

interface TechIconProps {
  name: string;
  size?: number;
}

export function TechIcon({ name, size = 75 }: TechIconProps) {
  const [isDarkMode, setIsDarkMode] = React.useState(
    window.matchMedia("(prefers-color-scheme: dark)").matches
  );
  name = name.replace(' ', '')
  if (name.includes('.')){
    name = name.replace('.', 'dot')
  }
  const url = isDarkMode ? `https://cdn.simpleicons.org/${name.toLowerCase()}/white` : `https://cdn.simpleicons.org/${name.toLowerCase()}/black`
  return (
    <img className={`w-[${size}] h-${size}`} src={url} />
  )
}

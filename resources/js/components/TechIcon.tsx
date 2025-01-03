import * as React from "react";

interface TechIconProps {
  name: string;
  size?: number;
}

export function TechIcon({ name, size = 75 }: TechIconProps) {
  return (
    <div className="p-2" style={{ width: size, height: size }}>
      {getIcon(name)}
    </div>
  );
}

function getIcon(name: string): React.JSX.Element {
  switch (name.toLowerCase()) {
    case "go":
      return <img src="https://www.svgrepo.com/show/353795/go.svg" alt="Golang" className="h-full w-full" />;
    default:
      return <img src="https://www.svgrepo.com/show/372737/unknown-status.svg" alt="Unknown" className="h-full w-full text-white stroke-white"/>;
  }
}

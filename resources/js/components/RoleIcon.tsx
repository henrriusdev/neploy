import * as LucideIcons from "lucide-react";
import * as RadixIcons from "@radix-ui/react-icons";

interface RoleIconProps {
  icon: string;
  color?: string;
  size?: number;
}

export function RoleIcon({ icon, color = "white", size = 40 }: RoleIconProps) {
  // Try Lucide icons first
  const LucideIcon = (LucideIcons as any)[icon];
  if (LucideIcon) {
    return <LucideIcon color={"white"} size={size} style={{ backgroundColor: color, borderRadius: "20%" }} className="p-1" />;
  }

  // If not found in Lucide, try Radix icons
  // Radix icons follow the pattern: IconNameIcon
  const radixIconName = `${icon}Icon`;
  const RadixIcon = (RadixIcons as any)[radixIconName];
  if (RadixIcon) {
    return (
      <div style={{ width: size, height: size, backgroundColor: color, color: "white", borderRadius: "20%" }}>
        <RadixIcon className="w-full h-full p-1" />
      </div>
    );
  }

  console.warn(`Icon not found: ${icon} (tried both Lucide and Radix)`);
  return null;
}

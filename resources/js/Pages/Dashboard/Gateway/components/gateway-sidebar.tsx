import { useState } from "react"
import { Link } from "@inertiajs/react"
import { 
  ChevronLeft, 
  Globe, 
  Lock, 
  Activity,
  BarChart3,
} from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"

const navItems = [
  {
    title: "Overview",
    icon: BarChart3,
    href: "#overview",
  },
  {
    title: "Routes",
    icon: Globe,
    href: "#routes",
  },
  {
    title: "Security",
    icon: Lock,
    href: "#security",
  },
  {
    title: "Rate Limiting",
    icon: Activity,
    href: "#rate-limiting",
  },
]

export function GatewaySidebar() {
  const [collapsed, setCollapsed] = useState(false)

  return (
    <div
      className={cn(
        "h-screen border-l bg-background transition-all duration-300",
        collapsed ? "w-16" : "w-64"
      )}
    >
      <div className="flex h-full flex-col">
        <div className="flex items-center justify-between p-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setCollapsed(!collapsed)}
            className={cn("", !collapsed && "rotate-180")}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          {!collapsed && <h2 className="text-lg font-semibold ml-2">Gateway</h2>}
        </div>

        <nav className="flex-1 space-y-2 p-2">
          {navItems.map((item) => (
            <Link
              key={item.title}
              href={item.href}
              className={cn(
                "flex items-center space-x-2 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:bg-accent hover:text-accent-foreground",
                collapsed ? "justify-center" : "justify-start"
              )}
            >
              <item.icon className="h-4 w-4" />
              {!collapsed && <span>{item.title}</span>}
            </Link>
          ))}
        </nav>
      </div>
    </div>
  )
}

import {
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {useTheme} from "@/hooks";
import {Check, Moon, Palette, Sun} from "lucide-react";
import {useTranslation} from "react-i18next";
import {Button} from "@/components/ui/button";
import {cn} from "@/lib/utils";

export const ThemeSwitcher = ({className = ''}: {className?: string}) => {
  const {theme, isDark, changeTheme, toggleDark, themes} = useTheme();
  const {t} = useTranslation()

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className={cn("px-4 py-2 flex justify-center items-center !rounded-md text-lg", className)}>
          <Palette className="mr-2 h-7 w-7"/>
          {t("theme")}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <div className="px-2 py-1.5 text-xs text-muted-foreground font-semibold">
          Tema
        </div>
        {themes.map((t) => (
          <DropdownMenuItem
            key={t}
            onClick={() => changeTheme(t)}
            className="flex items-center justify-between"
          >
            <span className="capitalize">{t}</span>
            {theme === t && <Check className="h-4 w-4"/>}
          </DropdownMenuItem>
        ))}
        <DropdownMenuSeparator/>
        <DropdownMenuItem onClick={toggleDark}>
          {isDark ? (
            <>
              <Sun className="mr-2 h-4 w-4"/> Light
            </>
          ) : (
            <>
              <Moon className="mr-2 h-4 w-4"/> Dark
            </>
          )}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

import {
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
} from "@/components/ui/dropdown-menu";
import {useTheme} from "@/hooks";
import {Check, Moon, Sun} from "lucide-react";

export const ThemeSwitcher = () => {
  const {theme, isDark, changeTheme, toggleDark, themes} = useTheme();

  return (
    <DropdownMenuSub>
      <DropdownMenuSubTrigger>Apariencia</DropdownMenuSubTrigger>
      <DropdownMenuSubContent>
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
      </DropdownMenuSubContent>
    </DropdownMenuSub>
  );
};

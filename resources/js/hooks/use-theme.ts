import { useEffect, useState } from "react";

const themes = ["neploy", "gruvbox", "rosepine", "tokyonight"] as const;
export type Theme = typeof themes[number] | "system";

export const useTheme = () => {
  const [theme, setTheme] = useState<Theme>("neploy");
  const [isDark, setIsDark] = useState(false);

  const applyTheme = (t: Theme, dark: boolean) => {
    document.documentElement.setAttribute("data-theme", t);
    document.documentElement.classList.toggle("dark", dark);
  };

  useEffect(() => {
    const storedTheme = (localStorage.getItem("theme") as Theme) || "neploy";
    const storedDark = localStorage.getItem("darkMode") === "true";
    setTheme(storedTheme);
    setIsDark(storedDark);
    applyTheme(storedTheme, storedDark);
  }, []);

  const changeTheme = (newTheme: Theme) => {
    localStorage.setItem("theme", newTheme);
    setTheme(newTheme);
    applyTheme(newTheme, isDark);
  };

  const toggleDark = () => {
    const nextDark = !isDark;
    localStorage.setItem("darkMode", String(nextDark));
    setIsDark(nextDark);
    applyTheme(theme, nextDark);
  };

  return { theme, isDark, changeTheme, toggleDark, themes, applyTheme };
};

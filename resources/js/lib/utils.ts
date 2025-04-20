import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function sanitizeAppName(appName: string): string {
  // Reemplazar espacios por guiones
  appName = appName.replace(/ /g, "-");

  // Eliminar cualquier carácter que no sea letra, número o guion
  appName = appName.replace(/[^a-zA-Z0-9-]/g, "");

  // Convertir a minúsculas
  return appName.toLowerCase();
}
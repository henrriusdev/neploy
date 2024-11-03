import React from "react";
import { FormControl, FormItem, FormLabel, FormMessage } from "./ui/form";

export const RenderFormItem = ({
  label,
  className,
  children,
}: React.PropsWithChildren<{ label: string, className: string }>) => (
  <FormItem className={className}>
    <FormLabel>{label}</FormLabel>
    <FormControl>{children}</FormControl>
    <FormMessage />
  </FormItem>
);
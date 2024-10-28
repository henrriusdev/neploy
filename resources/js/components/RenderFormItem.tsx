import React from "react";
import { FormControl, FormItem, FormLabel, FormMessage } from "./ui/form";

export const RenderFormItem = ({
  label,
  children,
}: React.PropsWithChildren<{ label: string }>) => (
  <FormItem>
    <FormLabel>{label}</FormLabel>
    <FormControl>{children}</FormControl>
    <FormMessage />
  </FormItem>
);
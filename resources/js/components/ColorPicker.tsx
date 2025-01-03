"use client";

import * as React from "react";
import { Pencil1Icon } from "@radix-ui/react-icons";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { ControllerRenderProps } from "react-hook-form";

export interface ColorPickerProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  className?: string;
  field: ControllerRenderProps<any>;
}

const ColorPicker = React.forwardRef<HTMLInputElement, ColorPickerProps>(
  ({ className, field, ...props }: ColorPickerProps, ref) => {
    const inputRef = React.useRef<HTMLInputElement>(null);
    const [isOpen, setIsOpen] = React.useState(false);
    const handleColorChange = (newColor: string) => {
      field.onChange?.(newColor);
      field.value = newColor;
      if (inputRef.current) {
        inputRef.current.value = newColor;
        const event = new Event("input", { bubbles: true });
        inputRef.current.dispatchEvent(event);
      }
    };

    return (
      <div className="grid gap-2 w-full max-w-full">
        <Popover open={isOpen} onOpenChange={setIsOpen}>
          <PopoverTrigger asChild>
            <Button
              variant={"outline"}
              className={cn(
                "w-full justify-start text-left font-normal",
                !field.value && "text-muted-foreground",
                className,
              )}
            >
              <div className="w-full flex items-center gap-2">
                {field.value && (
                  <div
                    className="h-4 w-4 rounded !bg-center !bg-cover transition-all border"
                    style={{ backgroundColor: field.value }}
                  ></div>
                )}
                <div className="truncate flex-1">
                  {field.value ? field.value.toUpperCase() : "Pick a field.value"}
                </div>
              </div>
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-64">
            <Tabs defaultValue="solid">
              <TabsList className="w-full mb-4">
                <TabsTrigger className="flex-1" value="solid">
                  Solid
                </TabsTrigger>
                <TabsTrigger className="flex-1" value="gradient">
                  Gradient
                </TabsTrigger>
              </TabsList>
              <div className="flex flex-wrap gap-1 mb-4">
                <Button
                  variant={"outline"}
                  className="w-6 h-6 p-0 flex items-center justify-center"
                  onClick={() => {
                    const eyeDropper = new (window as any).EyeDropper();
                    eyeDropper
                      .open()
                      .then((result: { sRGBHex: string }) => {
                        handleColorChange(result.sRGBHex);
                      })
                      .catch((error: any) => {
                        console.log(error);
                      });
                  }}
                >
                  <Pencil1Icon />
                </Button>
                {[
                  "#000000",
                  "#f44336",
                  "#e91e63",
                  "#9c27b0",
                  "#673ab7",
                  "#3f51b5",
                  "#2196f3",
                  "#03a9f4",
                  "#00bcd4",
                  "#009688",
                  "#4caf50",
                  "#8bc34a",
                  "#ffc107",
                  "#ff9800",
                  "#ff5722",
                  "#795548",
                  "#607d8b",
                ].map((color) => (
                  <Button
                    key={color}
                    variant={"outline"}
                    style={{ backgroundColor: color }}
                    className="w-6 h-6 p-0"
                    onClick={() => handleColorChange(color)}
                  />
                ))}
              </div>
              <div>
                <div className="flex items-center gap-2">
                  <div
                    className="w-[200px] h-4 rounded !bg-center !bg-cover transition-all border"
                    style={{
                      backgroundColor: field.value,
                    }}
                  />
                </div>
                <div className="flex items-center gap-2 mt-4">
                  <Input
                    ref={inputRef}
                    type="text"
                    value={field.value}
                    className="w-[200px]"
                    onChange={(e) => handleColorChange(e.target.value)}
                  />
                </div>
              </div>
            </Tabs>
          </PopoverContent>
        </Popover>
        <Input
          type="hidden"
          ref={ref}
          {...props}
          value={field.value}
          className={cn("", className)}
        />
      </div>
    );
  },
);
ColorPicker.displayName = "ColorPicker";

export { ColorPicker };

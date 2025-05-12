import {DialogButtonProps} from "@/types";
import React from "react";
import {Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger,} from "../ui/dialog";
import {TooltipButton} from "../ui/tooltip-button";
import {Button} from "../ui/button";

export const DialogButton = React.forwardRef<
  HTMLButtonElement,
  DialogButtonProps
>(
  (
    {
      buttonText,
      title,
      description,
      open,
      icon: Icon,
      onOpen,
      className,
      children,
      variant,
    },
    ref
  ) => (
    <Dialog open={open} onOpenChange={onOpen}>
      <DialogTrigger asChild>
        {variant === "tooltip" ? (
          <TooltipButton
            ref={ref}
            tooltip={buttonText}
            icon={Icon}
            variant="ghost"
            size="icon"
          />
        ) : (
          <Button
            ref={ref}
            variant="default"
            className={"flex items-center gap-2 " + className}>
            {buttonText} {Icon && <Icon className="h-4 w-4" />}
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className={"sm:max-w-[500px] " + className}>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        {children}
      </DialogContent>
    </Dialog>
  )
);

DialogButton.displayName = "DialogButton";

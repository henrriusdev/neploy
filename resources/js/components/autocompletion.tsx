import {
  CommandGroup,
  CommandItem,
  CommandList,
  CommandInput,
} from "./ui/command";
import { Command as CommandPrimitive } from "cmdk";
import {
  useState,
  useRef,
  useCallback,
  type KeyboardEvent,
  useEffect,
  forwardRef,
} from "react";
import { Skeleton } from "./ui/skeleton";
import { Check } from "lucide-react";
import { cn } from "../lib/utils";
import { ControllerRenderProps } from "react-hook-form";
import { AutoCompleteProps, Option } from "@/types/props";

export const AutoComplete = forwardRef<HTMLInputElement, AutoCompleteProps>(
  (
    {
      options,
      placeholder,
      emptyMessage,
      value,
      onValueChange,
      disabled,
      isLoading = false,
      field,
    },
    ref
  ) => {
    const inputRef = useRef<HTMLInputElement>(null);
    const [isOpen, setOpen] = useState(false);
    const [inputValue, setInputValue] = useState<string>("");

    // Initialize selected value from field or value prop
    const [selected, setSelected] = useState<Option | undefined>(() => {
      if (field?.value) {
        const fieldOption = options.find((opt) => opt.value === field.value);
        return fieldOption;
      }
      return value;
    });

    // Sync input value with selected option
    useEffect(() => {
      if (selected) {
        setInputValue(selected.label);
      }
    }, [selected]);

    // Sync with external value changes
    useEffect(() => {
      if (field?.value) {
        const fieldOption = options.find((opt) => opt.value === field.value);
        if (
          fieldOption &&
          (!selected || selected.value !== fieldOption.value)
        ) {
          setSelected(fieldOption);
          setInputValue(fieldOption.label);
        }
      }
    }, [field?.value, options]);

    const handleKeyDown = useCallback(
      (event: KeyboardEvent<HTMLDivElement>) => {
        const input = inputRef.current;
        if (!input) return;

        if (!isOpen) {
          setOpen(true);
        }

        if (event.key === "Enter" && input.value !== "") {
          const optionToSelect = options.find(
            (option) => option.label.toLowerCase() === input.value.toLowerCase()
          );
          if (optionToSelect) {
            handleSelectOption(optionToSelect);
          }
        }

        if (event.key === "Escape") {
          input.blur();
        }
      },
      [isOpen, options]
    );

    const handleBlur = useCallback(() => {
      setOpen(false);
      if (selected) {
        setInputValue(selected.label);
      } else {
        setInputValue("");
      }
    }, [selected]);

    const handleSelectOption = useCallback(
      (selectedOption: Option) => {
        setSelected(selectedOption);
        setInputValue(selectedOption.label);

        // Update form field
        if (field) {
          field.onChange(selectedOption.value);
        }

        // Call external onChange handler
        if (onValueChange) {
          onValueChange(selectedOption);
        }

        setTimeout(() => {
          inputRef?.current?.blur();
        }, 0);
      },
      [field, onValueChange]
    );

    return (
      <CommandPrimitive onKeyDown={handleKeyDown}>
        <div>
          <CommandInput
            ref={ref || inputRef}
            value={inputValue}
            onValueChange={isLoading ? undefined : setInputValue}
            onBlur={handleBlur}
            onFocus={() => setOpen(true)}
            placeholder={placeholder}
            disabled={disabled}
            className="text-base"
          />
        </div>
        <div className="relative mt-1">
          <div
            className={cn(
              "animate-in fade-in-0 zoom-in-95 absolute top-0 z-10 w-full rounded-xl bg-white outline-none",
              isOpen ? "block" : "hidden"
            )}>
            <CommandList className="rounded-lg ring-1 ring-slate-200">
              {isLoading ? (
                <CommandPrimitive.Loading>
                  <div className="p-1">
                    <Skeleton className="h-8 w-full" />
                  </div>
                </CommandPrimitive.Loading>
              ) : null}
              {options.length > 0 && !isLoading ? (
                <CommandGroup>
                  {options.map((option) => {
                    const isSelected = selected?.value === option.value;
                    return (
                      <CommandItem
                        key={option.value}
                        value={option.label}
                        onMouseDown={(event) => {
                          event.preventDefault();
                          event.stopPropagation();
                        }}
                        onSelect={() => handleSelectOption(option)}
                        className={cn(
                          "flex w-full items-center gap-2",
                          !isSelected ? "pl-8" : null
                        )}>
                        {isSelected ? <Check className="w-4" /> : null}
                        {option.label}
                      </CommandItem>
                    );
                  })}
                </CommandGroup>
              ) : null}
              {!isLoading ? (
                <CommandPrimitive.Empty className="select-none rounded-sm px-2 py-3 text-center text-sm">
                  {emptyMessage}
                </CommandPrimitive.Empty>
              ) : null}
            </CommandList>
          </div>
        </div>
      </CommandPrimitive>
    );
  }
);

AutoComplete.displayName = "AutoComplete";

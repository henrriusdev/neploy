"use client";

import { addDays, format, lastDayOfYear } from "date-fns";
import { Calendar as CalendarIcon } from "lucide-react";
import { DateRange } from "react-day-picker";
import { useFormContext } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { FormControl } from "@/components/ui/form";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { cn } from "@/lib/utils";
import { DatePickerProps } from "@/types/props";
import { forwardRef, useEffect, useState } from "react";

export const DatePicker = forwardRef<HTMLDivElement, DatePickerProps>(({ className, date, onDateChange, isRangePicker = false, minYear = 1900, maxYear = 2100, field }, ref) => {
  const formContext = useFormContext();
  const isFormContext = !!formContext && !!field;

  const [selectedDate, setSelectedDate] = useState<Date | DateRange | undefined>(isFormContext ? field.value : date);
  const [isRange, setIsRange] = useState<boolean>(isRangePicker);
  const [month, setMonth] = useState<Date>(() => {
    if (maxYear) {
      return lastDayOfYear(new Date(maxYear, 0, 1));
    }
    return selectedDate instanceof Date ? selectedDate : selectedDate && "from" in selectedDate ? selectedDate.from : new Date();
  });

  useEffect(() => {
    const currentValue = isFormContext ? field?.value : selectedDate;

    if (currentValue instanceof Date) {
      setMonth(currentValue);
    } else if (currentValue && typeof currentValue === "object" && "from" in currentValue && currentValue.from instanceof Date) {
      setMonth(currentValue.from);
    } else {
      setMonth(new Date()); // fallback para asegurar consistencia
    }

    setSelectedDate(currentValue);
  }, [isFormContext, field?.value, selectedDate]);

  useEffect(() => {
    setIsRange(isRangePicker);
  }, [isRangePicker]);

  const handleDateSelect = (newDate: Date | DateRange | undefined) => {
    setSelectedDate(newDate);

    console.log(newDate);
    if (newDate instanceof Date) {
      setMonth(newDate);
    } else if (newDate && "from" in newDate && newDate.from instanceof Date) {
      setMonth(newDate.from);
    }

    if (isFormContext) {
      field.onChange(newDate);
    } else {
      onDateChange?.(newDate);
    }
  };

  const formatDate = (date: Date | DateRange | undefined) => {
    if (!date) return "Pick a date";
    if (date instanceof Date) return format(date, "yyyy-MM-dd");
    if (date.from) {
      if (date.to) return `${format(date.from, "yyyy-MM-dd")} - ${format(date.to, "yyyy-MM-dd")}`;
      return `${format(date.from, "yyyy-MM-dd")} - `;
    }
    return "Pick a date range";
  };

  const buttonProps = isFormContext ? field : {};

  return (
    <div ref={ref} className={cn("grid gap-2", className)}>
      <Popover>
        <PopoverTrigger asChild>
          {isFormContext ? (
            <FormControl>
              <Button id={field.name} variant={"outline"} className={cn("w-full justify-start text-left font-normal", !selectedDate && "text-muted-foreground")} {...buttonProps}>
                <CalendarIcon className="mr-2 h-4 w-4" />
                {formatDate(selectedDate)}
              </Button>
            </FormControl>
          ) : (
            <Button variant={"outline"} className={cn("w-full justify-start text-left font-normal", !selectedDate && "text-muted-foreground")} {...buttonProps}>
              <CalendarIcon className="mr-2 h-4 w-4" />
              {formatDate(selectedDate)}
            </Button>
          )}
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0 bg-background" align="start">
          <div className="w-[350px] p-3 space-y-3 flex items-center justify-between">
            <div className="flex space-x-1">
              <Select
                value={month instanceof Date && !isNaN(month.getTime()) ? format(month, "MMMM") : ""}
                onValueChange={(value) => {
                  if (month instanceof Date && !isNaN(month.getTime())) {
                    setMonth(new Date(month.getFullYear(), parseInt(value), 1));
                  }
                }}>
                <SelectTrigger className="w-[160px] text-white">
                  <SelectValue>{month instanceof Date && !isNaN(month.getTime()) ? format(month, "MMMM") : ""}</SelectValue>
                </SelectTrigger>
                <SelectContent>
                  {Array.from({ length: 12 }, (_, i) => (
                    <SelectItem key={i} value={i.toString()}>
                      {format(new Date(0, i), "MMMM")}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Select
                value={month?.getFullYear().toString() ?? new Date(maxYear).getFullYear().toString()}
                onValueChange={(value) => {
                  if (month instanceof Date && !isNaN(month.getTime())) {
                    setMonth(new Date(parseInt(value), month.getMonth(), 1));
                  }
                }}>
                <SelectTrigger className="w-[160px] text-white">
                  <SelectValue>{month?.getFullYear().toString() ?? new Date(maxYear, 1, 1).getFullYear().toString()}</SelectValue>
                </SelectTrigger>
                <SelectContent>
                  {Array.from({ length: maxYear - minYear + 1 }, (_, i) => {
                    const year = maxYear - i;
                    return (
                      <SelectItem key={i} value={year.toString()}>
                        {year}
                      </SelectItem>
                    );
                  })}
                </SelectContent>
              </Select>
            </div>
          </div>
          <Calendar
            mode={isRange ? "range" : "single"}
            selected={selectedDate}
            onSelect={handleDateSelect}
            month={month}
            onMonthChange={setMonth}
            numberOfMonths={1}
            fromYear={minYear}
            toYear={maxYear}
            className="p-3"
            classNames={{
              months: "space-y-4",
              month: "space-y-4",
              caption: "flex justify-center pt-1 relative items-center",
              caption_label: "text-sm font-medium text-white",
              nav: "space-x-1 flex items-center text-white",
              nav_button: "h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100 text-white",
              nav_button_previous: "absolute left-1 text-white",
              nav_button_next: "absolute right-1 text-white",
              table: "!w-full border-collapse space-y-1",
              head_row: "flex !justify-center !w-full",
              head_cell: "text-white rounded-md w-10 font-normal text-[0.8rem]",
              row: "flex !w-full !justify-center mt-2",
              cell: "text-center text-sm p-0 relative [&:has([aria-selected])]:bg-accent first:[&:has([aria-selected])]:rounded-l-md last:[&:has([aria-selected])]:rounded-r-md focus-within:relative focus-within:z-20",
              day: "h-10 w-10 p-0 font-normal aria-selected:opacity-100 text-white",
              day_selected: "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground",
              day_today: "bg-secondary text-accent-foreground",
              day_outside: "text-muted-foreground opacity-50",
              day_disabled: "text-muted-foreground opacity-50",
              day_range_middle: "aria-selected:bg-accent aria-selected:text-accent-foreground",
              day_hidden: "invisible",
            }}
          />
          {maxYear >= new Date().getFullYear() && (
            <div className="flex items-center justify-between p-3 border-t">
              <Button variant="ghost" onClick={() => handleDateSelect(new Date())}>
                Today
              </Button>
              {isRange && (
                <div className="flex gap-2">
                  <Button
                    variant="ghost"
                    onClick={() =>
                      handleDateSelect({
                        from: new Date(),
                        to: addDays(new Date(), 7),
                      })
                    }>
                    Next 7 days
                  </Button>
                  <Button
                    variant="ghost"
                    onClick={() =>
                      handleDateSelect({
                        from: new Date(),
                        to: addDays(new Date(), 30),
                      })
                    }>
                    Next 30 days
                  </Button>
                </div>
              )}
            </div>
          )}
        </PopoverContent>
      </Popover>
    </div>
  );
});

DatePicker.displayName = "DatePicker";

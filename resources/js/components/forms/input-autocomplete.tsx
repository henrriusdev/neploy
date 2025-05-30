import {ControllerRenderProps} from "react-hook-form";
import {AutoComplete} from "./autocompletion";
import {useEffect, useState} from "react";
import {Option} from "@/types/props";

export function InputAutoComplete({
  OPTIONS,
  field,
  placeholder,
}: {
  OPTIONS: Option[];
  field: ControllerRenderProps<any>;
  placeholder?: string;
}) {
  const [isLoading, setLoading] = useState(false);
  const [isDisabled, setDisabled] = useState(false);
  const [selectedOption, setSelectedOption] = useState<Option | undefined>(
    () => {
      return OPTIONS.find((opt) => opt.value === field.value);
    }
  );

  // Update the displayed value when field.value changes
  useEffect(() => {
    const option = OPTIONS.find((opt) => opt.value === field.value);
    setSelectedOption(option);
  }, [field.value, OPTIONS]);

  const handleValueChange = (option: Option) => {
    field.onChange(option.value);
    setSelectedOption(option);
    field.value = option.value;
  };

  return (
    <div className="not-prose mt-8 flex flex-col gap-4">
      <AutoComplete
        options={OPTIONS}
        emptyMessage="No results."
        placeholder={placeholder}
        isLoading={isLoading}
        disabled={isDisabled}
        value={selectedOption}
        onValueChange={handleValueChange}
      />
    </div>
  );
}

import { ControllerRenderProps } from "react-hook-form";
import { AutoComplete, type Option } from "./autocompletion";
import { useState } from "react";


export function InputAutoComplete({ OPTIONS, field }: {
  OPTIONS: Option[];
  field: ControllerRenderProps<any>;
}) {
  const [isLoading, setLoading] = useState(false);
  const [isDisabled, setDisbled] = useState(false);

  return (
    <div className="not-prose mt-8 flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <button
          className="inline-flex h-10 items-center justify-center rounded-md border border-stone-200 bg-white px-4 py-2 text-sm font-medium ring-offset-white transition-colors hover:bg-stone-100 hover:text-slate-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
          onClick={() => setLoading((prev) => !prev)}>
          Toggle loading
        </button>
        <button
          className="inline-flex h-10 items-center justify-center rounded-md border border-stone-200 bg-white px-4 py-2 text-sm font-medium ring-offset-white transition-colors hover:bg-stone-100 hover:text-slate-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
          onClick={() => setDisbled((prev) => !prev)}>
          Toggle disabled
        </button>
      </div>
      <AutoComplete
        options={OPTIONS}
        emptyMessage="No resulsts."
        placeholder="Find something"
        isLoading={isLoading}
        onValueChange={field.onChange}
        value={field.value}
        disabled={isDisabled}
        {...field}
      />
      <span className="text-sm">
        Current value: {field.value ? field.value?.label : "No value selected"}
      </span>
      <span className="text-sm">
        Loading state: {isLoading ? "true" : "false"}
      </span>
      <span className="text-sm">Disabled: {isDisabled ? "true" : "false"}</span>
    </div>
  );
}

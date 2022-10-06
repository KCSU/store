import { CreatableSelect, Size } from "chakra-react-select";
import { forwardRef } from "react";

interface TicketOptionsSelectProps {
  value?: string;
  size?: Size;
  onChange?: (value: string) => void;
}

export function TicketOptionsSelect({
  value,
  size,
  onChange,
}: TicketOptionsSelectProps) {
  const options = ["Normal", "Vegetarian", "Vegan", "Pescetarian"];
  return (
    <CreatableSelect
      // TODO: fix long answers
      isClearable
      isValidNewOption={() => true}
      formatCreateLabel={(input) => (
        input.length > 0 ? `Custom: "${input}"` : "Add a custom option..."
      )}
      selectedOptionColor="brand"
      value={{
        label: value,
        value,
      }}
      onChange={(option) => {
        onChange?.(option?.value ?? "");
      }}
      size={size}
      options={options.map((opt) => ({
        label: opt,
        value: opt,
      }))}
    ></CreatableSelect>
  );
}

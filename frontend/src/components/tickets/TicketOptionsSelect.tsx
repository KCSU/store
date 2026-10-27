import { CreatableSelect, Size } from "chakra-react-select";
import { forwardRef } from "react";

interface TicketOptionsSelectProps {
  value?: string;
  size?: Size;
  onChange?: (value: string) => void;
}

// required, at most 100 characters
function isValidOption(input: string) {
  const length = input.trim().length;
  return length > 0 && length <= 100;
}

export function TicketOptionsSelect({
  value,
  size,
  onChange,
}: TicketOptionsSelectProps) {
  const options = ["Normal", "Vegetarian", "Vegan", "Pescetarian"];
  return (
    <CreatableSelect
      isValidNewOption={isValidOption}
      onInputChange={(input, { action, prevInputValue }) => {
        if (action === "input-blur" && isValidOption(prevInputValue)) {
          onChange?.(prevInputValue.trim());
        }
      }}
      formatCreateLabel={(input) => `Custom: "${input.trim()}"`}
      noOptionsMessage={({ inputValue }) =>
        inputValue.trim().length > 100
          ? "Custom options can be at most 100 characters"
          : "Type a custom option"
      }
      getNewOptionData={(input) => ({ label: input.trim(), value: input.trim() })}
      selectedOptionColor="brand"
      value={{
        label: value,
        value,
      }}
      onChange={(option) => {
        if (option?.value) {
          onChange?.(option.value);
        }
      }}
      size={size}
      options={options.map((opt) => ({
        label: opt,
        value: opt,
      }))}
    ></CreatableSelect>
  );
}

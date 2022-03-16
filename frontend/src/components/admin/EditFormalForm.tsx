import {
  Icon,
} from "@chakra-ui/react";
import { FaSave } from "react-icons/fa";
import { useEditFormal } from "../../hooks/admin/useEditFormal";
import { Formal } from "../../model/Formal";
import { FormalDetailsForm } from "./FormalDetailsForm";

interface FormalProps {
  formal: Formal;
}

export function EditFormalForm({ formal }: FormalProps) {
  const mutation = useEditFormal(formal.id);
  return (
    <FormalDetailsForm
      formal={formal}
      onSubmit={async (values) => {
        await mutation.mutateAsync(values);
      }}
      submitIcon={<Icon as={FaSave} />}
    >
      Save Changes
    </FormalDetailsForm>
  );
}

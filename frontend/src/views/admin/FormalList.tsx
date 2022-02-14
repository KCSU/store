import { Box } from "@chakra-ui/react";
import { useAllFormals } from "../../hooks/admin/useAllFormals";

export function FormalList() {
  const { data, isLoading, isError } = useAllFormals();
  if (!data) {
    return <></>;
  }

  return (
    <>
      {data.map((f) => (
        <Box key={f.id}>{f.name}</Box>
      ))}
    </>
  );
}

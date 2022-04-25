import {
  UseRadioProps,
  useRadio,
  useColorModeValue,
  Heading,
  Text,
  Alert,
  useRadioGroup,
  VStack,
} from "@chakra-ui/react";
import { useMemo } from "react";
import { useAllFormals } from "../../hooks/admin/useAllFormals";
import { useDateTime } from "../../hooks/state/useDateTime";
import { Formal } from "../../model/Formal";
import { Card } from "../utility/Card";

type FormalRadioCardProps = UseRadioProps & {
  formal: Formal;
};

function FormalRadioCard({ formal, ...props }: FormalRadioCardProps) {
  const { getInputProps, getCheckboxProps } = useRadio(props);
  const input = getInputProps();
  const checkbox = getCheckboxProps();
  const bg = useColorModeValue("gray.100", "gray.600");
  const checkedBg = useColorModeValue("brand.500", "brand.200");
  const checkedFg = useColorModeValue("white", "black");
  const date = useDateTime(formal.dateTime);
  return (
    <Card
      p={3}
      bg={bg}
      borderRadius="md"
      as="label"
      {...checkbox}
      cursor="pointer"
      transition="all 0.2s"
      _checked={{
        bg: checkedBg,
        color: checkedFg,
      }}
    >
      <input {...input} />
      <Heading as="h4" size="sm">
        {formal.name}
      </Heading>
      <Text fontSize="sm">{date}</Text>
    </Card>
  );
}

interface FormalRadioGroupProps {
  exclude?: string[];
  onChange?: (value: string) => void;
  value?: string;
}

export function FormalRadioGroup({
  exclude,
  onChange,
  value,
}: FormalRadioGroupProps) {
  const { data, isError } = useAllFormals();
  const { getRootProps, getRadioProps } = useRadioGroup({
    name: "formals",
    value,
    onChange,
  });
  const group = getRootProps();
  // TODO: Pagination / limits / search
  const formals = useMemo(() => {
    return data?.filter((f) => !exclude?.includes(f.id)) ?? [];
  }, [data, exclude]);
  return isError ? (
    <Alert status="error">
      There was an error loading the formals. Please try again later.
    </Alert>
  ) : (
    <VStack {...group} gap={2}>
      {formals?.map((formal) => {
        const radio = getRadioProps({ value: formal.id });
        return <FormalRadioCard key={formal.id} formal={formal} {...radio} />;
      })}
    </VStack>
  );
}

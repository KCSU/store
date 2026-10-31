import { Box, Icon, Tooltip } from "@chakra-ui/react";
import { FaQuestionCircle } from "react-icons/fa";

const HINT =
  "Once this formal is visible, buyers can join the pre-sale queue at any " +
  "time before the sale start time. Tickets are allocated randomly to the " +
  "queue when the sale starts. Any remaining tickets will then be " +
  "first-come first-served.";

export function SaleStartHint() {
  return (
    <Tooltip hasArrow label={HINT}>
      <Box as="span" ml={1} tabIndex={0} aria-label={HINT}>
        <Icon as={FaQuestionCircle} boxSize={3} verticalAlign="baseline" />
      </Box>
    </Tooltip>
  );
}

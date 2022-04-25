import { Heading } from "@chakra-ui/react";
import { useContext } from "react";
import { BillContext } from "../../model/Bill";

export function BillFormalsList() {
  const bill = useContext(BillContext);
  return <Heading size="md">{bill.name}</Heading>
}
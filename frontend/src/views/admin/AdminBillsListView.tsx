import {
  Heading,
  LinkBox,
  LinkOverlay,
  SimpleGrid,
  Text,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import { useMemo } from "react";
import { Link } from "react-router-dom";
import { Card } from "../../components/utility/Card";
import { useBills } from "../../hooks/admin/useBills";
import { Bill } from "../../model/Bill";

interface BillProps {
  bill: Bill;
}

function AdminBillCard({ bill }: BillProps) {
  const start = useMemo(() => dayjs(bill.start).format("ll"), [bill]);
  const end = useMemo(() => dayjs(bill.end).format("ll"), [bill]);
  return (
    <LinkBox
      as={Card}
      borderRadius={3}
      p={4}
      transition="box-shadow 0.2s"
      _hover={{ shadow: "lg" }}
      _focusWithin={{ shadow: "lg" }}
    >
      <LinkOverlay as={Link} to={`/admin/bills/${bill.id}`}>
        <Heading as="h5" size="sm" mb={2}>
          {bill.name}
        </Heading>
      </LinkOverlay>
      <Text fontSize="sm">{start} &mdash; {end}</Text>
    </LinkBox>
  );
}

export function AdminBillsListView() {
  const { data, isLoading, isError } = useBills();
  if (!data) {
    return <></>;
  }

  return (
    <>
      <Heading size="xl" mb={5}>
        Manage Bills
      </Heading>
      <SimpleGrid
        templateColumns="repeat(auto-fill, minmax(250px, 1fr))"
        spacing="20px"
      >
        {data.map((bill) => (
          <AdminBillCard key={bill.id} bill={bill} />
        ))}
      </SimpleGrid>
    </>
  );
}
import {
  Button,
  Heading,
  LinkBox,
  LinkOverlay,
  SimpleGrid,
  Text,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import { useMemo } from "react";
import { FaPlus } from "react-icons/fa";
import { Link } from "react-router-dom";
import { Card } from "../../components/utility/Card";
import { useBills } from "../../hooks/admin/useBills";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
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
      <Text fontSize="sm">{start} &ndash; {end}</Text>
    </LinkBox>
  );
}

export function AdminBillsListView() {
  const { data, isLoading, isError } = useBills();
  const canWrite = useHasPermission("billing", "write");
  if (!data) {
    return <></>;
  }

  return (
    <>
      <Heading size="xl" mb={5}>
        Manage Bills
      </Heading>
      {canWrite && (
        <Button
          colorScheme="brand"
          mb={4}
          as={Link}
          to="/admin/bills/create"
          leftIcon={<FaPlus />}
        >
          Create Bill
        </Button>
      )}
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
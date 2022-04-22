import {
  Container,
  Heading,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  Box,
} from "@chakra-ui/react";
import { Navigate, useParams } from "react-router-dom";
import { EditBillForm } from "../../components/admin/EditBillForm";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { useBill } from "../../hooks/admin/useBill";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Bill } from "../../model/Bill";

interface BillProps {
  bill: Bill;
}

function AdminEditBillCard({ bill }: BillProps) {
  const canWrite = useHasPermission("billing", "write");
  const canDelete = useHasPermission("billing", "delete");
  const hasActions = canWrite || canDelete;
  return (
    <Container maxW="container.md" p={0}>
      <BackButton to="/admin/bills">Back Home</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={4}>
          {bill.name}
        </Heading>
        <Tabs colorScheme="brand">
          <TabList flexWrap="wrap">
            <Tab>Bill Details</Tab>
            <Tab>Formals</Tab>
            <Tab>Stats</Tab>
            {hasActions && <Tab>Actions</Tab>}
          </TabList>
          <TabPanels>
            <TabPanel>
              <EditBillForm bill={bill} />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Card>
    </Container>
  );
}

export function AdminEditBillView() {
  const { id } = useParams();
  const { data: bill, isLoading, isError } = useBill(id ?? "");
  if (isError) {
    return <Navigate to="/admin/bills" />;
  }
  // TODO: loading states
  if (isLoading && !bill) {
    return <Box></Box>;
  }
  if (!bill) {
    return <Box></Box>;
  }
  return <AdminEditBillCard bill={bill}/>;
}

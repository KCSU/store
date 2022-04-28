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
import { useContext } from "react";
import { Navigate, useParams } from "react-router-dom";
import { BillFormalsList } from "../../components/admin/BillFormalsList";
import { EditBillForm } from "../../components/admin/EditBillForm";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { useBill } from "../../hooks/admin/useBill";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { BillContext } from "../../model/Bill";
import { BillStats } from "../../components/formals/BillStats";

function AdminEditBillCard() {
  const bill = useContext(BillContext);
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
        <Tabs colorScheme="brand" isLazy>
          <TabList flexWrap="wrap">
            <Tab>Bill Details</Tab>
            <Tab>Formals</Tab>
            <Tab>Stats</Tab>
            {hasActions && <Tab>Actions</Tab>}
          </TabList>
          <TabPanels>
            <TabPanel>
              <EditBillForm />
            </TabPanel>
            <TabPanel>
              <BillFormalsList />
            </TabPanel>
            <TabPanel>
              <BillStats />
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
  return <BillContext.Provider value={bill}>
    <AdminEditBillCard/>
  </BillContext.Provider>;
}
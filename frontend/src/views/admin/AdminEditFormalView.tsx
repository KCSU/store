import {
  Box,
  Container,
  Heading,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
} from "@chakra-ui/react";
import { Navigate, useParams } from "react-router-dom";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { Formal, FormalContext } from "../../model/Formal";
import { EditFormalForm } from "../../components/admin/EditFormalForm";
import { useFormal } from "../../hooks/admin/useFormal";
import { EditFormalGroupsForm } from "../../components/admin/EditFormalGroupsForm";
import { FormalStats } from "../../components/admin/FormalStats";
import { FormalTicketsList } from "../../components/admin/FormalTicketsList";
import { FormalManualTicketsList } from "../../components/admin/FormalManualTicketsList";
import { useContext } from "react";

function AdminEditFormalCard() {
  const formal = useContext(FormalContext);
  return (
    <Container maxW="container.md" p={0}>
      <BackButton to="/admin/formals">Back Home</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={4}>
          {formal.name}
        </Heading>
        <Tabs colorScheme="brand">
          <TabList flexWrap="wrap">
            <Tab>Event Details</Tab>
            <Tab>Groups</Tab>
            <Tab>Manage Tickets</Tab>
            <Tab>Special Tickets</Tab>
            <Tab>Formal Stats</Tab>
          </TabList>
          <TabPanels>
            <TabPanel>
              <EditFormalForm />
            </TabPanel>
            <TabPanel>
              <EditFormalGroupsForm />
            </TabPanel>
            <TabPanel>
              <FormalTicketsList />
            </TabPanel>
            <TabPanel>
              <FormalManualTicketsList />
            </TabPanel>
            <TabPanel>
              <FormalStats />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Card>
    </Container>
  );
}

export function AdminEditFormalView() {
  const { formalId } = useParams();
  const { data: formal, isLoading, isError } = useFormal(formalId ?? "");
  if (isError) {
    // TODO: return an error!
    return <Navigate to="/admin/formals" />;
  }
  if (isLoading && !formal) {
    // TODO: return something better!
    return <Box></Box>;
  }
  if (!formal) {
    // Hmmm...
    return <Box></Box>;
  }
  return (
    <FormalContext.Provider value={formal}>
      <AdminEditFormalCard />
    </FormalContext.Provider>
  );
}

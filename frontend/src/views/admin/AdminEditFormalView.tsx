import {
  Box,
  Container,
  Heading,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
} from "@chakra-ui/react";
import { Navigate, useParams } from "react-router-dom";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { Formal } from "../../model/Formal";
import { EditFormalForm } from "../../components/admin/EditFormalForm";
import { useAllFormals } from "../../hooks/admin/useAllFormals";

interface FormalProps {
  formal: Formal;
}

function AdminEditFormalCard({ formal }: FormalProps) {
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
              <EditFormalForm formal={formal} />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Card>
    </Container>
  );
}

export function AdminEditFormalView() {
  const { formalId } = useParams();
  const formalIdNum = parseInt(formalId ?? "0");
  const { data: formals, isLoading, isError } = useAllFormals();
  const formal = formals?.find((f) => f.id === formalIdNum);
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
  return <AdminEditFormalCard formal={formal} />;
}

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
import { EditGroupForm } from "../../components/admin/EditGroupForm";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { useGroup } from "../../hooks/admin/useGroup";
import { Group } from "../../model/Group";

interface GroupProps {
  group: Group;
}

function AdminEditGroupCard({ group }: GroupProps) {
  return (
    <Container maxW="container.md" p={0}>
      <BackButton to="/admin/groups">Back Home</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={4}>
          {group.name}
        </Heading>
        <Tabs colorScheme="brand">
          <TabList flexWrap="wrap">
            <Tab>Group Details</Tab>
            <Tab>Directory</Tab>
            <Tab>Manual Members</Tab>
          </TabList>
          <TabPanels>
            <TabPanel>
              <EditGroupForm group={group} />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Card>
    </Container>
  );
}

export function AdminEditGroupView() {
  const { groupId } = useParams();
  const groupIdNum = parseInt(groupId ?? "0");
  const { data: group, isLoading, isError } = useGroup(groupIdNum);
  if (isError) {
    // TODO: return an error!
    return <Navigate to="/admin/groups" />;
  }
  if (isLoading && !group) {
    // TODO: return something better!
    return <Box></Box>;
  }
  if (!group) {
    // Hmmm...
    return <Box></Box>;
  }
  return <AdminEditGroupCard group={group} />;
}

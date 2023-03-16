import {
  Container,
  Heading,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  Box,
  Flex,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialog,
  AlertDialogFooter,
  Button,
  Icon,
  IconButton,
  useDisclosure,
} from "@chakra-ui/react";
import { useContext, useRef } from "react";
import { Navigate, useNavigate, useParams } from "react-router-dom";
import { BillFormalsList } from "../../components/admin/BillFormalsList";
import { EditBillForm } from "../../components/admin/EditBillForm";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { useBill } from "../../hooks/admin/useBill";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { BillContext } from "../../model/Bill";
import { BillStats } from "../../components/formals/BillStats";
import { FaTrashAlt } from "react-icons/fa";
import { useDeleteBill } from "../../hooks/admin/useDeleteBill";
import { BillExtrasList } from "../../components/admin/BillExtrasList";

function AdminEditBillCard() {
  const bill = useContext(BillContext);
  const canDelete = useHasPermission("billing", "delete");
  return (
    <Container maxW="container.md" p={0}>
      <BackButton to="/admin/bills">Back Home</BackButton>
      <Card mb={5}>
        <Flex gap={3}>
          <Heading as="h3" size="lg" mb={4} flex="1">
            {bill.name}
          </Heading>
          {canDelete && <DeleteBillButton />}
        </Flex>
        <Tabs colorScheme="brand" isLazy>
          <TabList flexWrap="wrap">
            <Tab>Bill Details</Tab>
            <Tab>Formals</Tab>
            <Tab>Extras</Tab>
            <Tab isDisabled={bill.formals?.length === 0}>Stats</Tab>
          </TabList>
          <TabPanels>
            <TabPanel>
              <EditBillForm />
            </TabPanel>
            <TabPanel>
              <BillFormalsList />
            </TabPanel>
            <TabPanel>
              <BillExtrasList />
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

function DeleteBillButton() {
  const bill = useContext(BillContext);
  const mutation = useDeleteBill(bill.id);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const navigate = useNavigate();
  const cancelRef = useRef(null);
  return (
    <>
      <IconButton
        onClick={onOpen}
        size="sm"
        colorScheme="red"
        aria-label="Delete"
        icon={<Icon as={FaTrashAlt} />}
      ></IconButton>
      <AlertDialog
        isOpen={isOpen}
        onClose={onClose}
        leastDestructiveRef={cancelRef}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Bill
            </AlertDialogHeader>
            <AlertDialogBody>
              Are you sure? This can't be undone.
            </AlertDialogBody>
            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              <Button
                colorScheme="red"
                onClick={async () => {
                  await mutation.mutateAsync();
                  onClose();
                  navigate("/admin/bills");
                }}
                ml={3}
                isLoading={mutation.isLoading}
              >
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
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
  return (
    <BillContext.Provider value={bill}>
      <AdminEditBillCard />
    </BillContext.Provider>
  );
}

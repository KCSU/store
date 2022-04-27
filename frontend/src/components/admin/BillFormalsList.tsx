import {
  Alert,
  AlertIcon,
  Box,
  Button,
  CloseButton,
  Heading,
  Icon,
  LinkBox,
  LinkOverlay,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  SimpleGrid,
  Text,
  useColorModeValue,
  useDisclosure,
  VStack,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import { Link } from "react-router-dom";
import { useContext, useMemo, useState } from "react";
import { BillContext } from "../../model/Bill";
import { Card } from "../utility/Card";
import { FaArrowRight, FaExternalLinkAlt, FaPlus } from "react-icons/fa";
import { FormalRadioGroup } from "./FormalRadioGroup";
import { useAddBillFormals } from "../../hooks/admin/useAddBillFormals";
import { Formal } from "../../model/Formal";
import { useRemoveBillFormal } from "../../hooks/admin/useRemoveBillFormal";
import { useAllFormals } from "../../hooks/admin/useAllFormals";

function AddFormalButton() {
  const bill = useContext(BillContext);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const exclude = useMemo(() => bill.formals?.map((f) => f.id), [bill]);
  const [formalId, setFormalId] = useState("");
  const mutation = useAddBillFormals(bill.id);
  return (
    <>
      <Button
        colorScheme="brand"
        leftIcon={<Icon as={FaPlus} />}
        onClick={onOpen}
      >
        Add Formal
      </Button>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Add Formal</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <FormalRadioGroup
              exclude={exclude}
              value={formalId}
              onChange={setFormalId}
            />
          </ModalBody>
          <ModalFooter>
            <Button variant="ghost" mr={3} onClick={onClose}>
              Cancel
            </Button>
            <Button
              colorScheme="brand"
              isDisabled={formalId === ""}
              isLoading={mutation.isLoading}
              onClick={async () => {
                await mutation.mutateAsync([formalId]);
                onClose();
                setFormalId("");
              }}
            >
              Add
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

function BillFormalsPreview() {
  const bill = useContext(BillContext);
  const { data, isError } = useAllFormals();
  const bg = useColorModeValue("gray.100", "gray.600");
  const formals = useMemo(() => {
    return (
      data?.filter((f) => bill.start <= f.dateTime && f.dateTime <= bill.end) ??
      []
    );
  }, [data, bill]);
  const mutation = useAddBillFormals(bill.id);
  if (isError) {
    return <Alert status="error">Error loading formals</Alert>;
  }
  return (
    <>
      <Alert status="info">
        <AlertIcon />
        No formals have been added to this bill yet.
      </Alert>
      <Heading as="h4" size="md" mt={3}>
        Selected Formals
      </Heading>
      <Text fontStyle="italic">
        {dayjs(bill.start).format("MMM D, YYYY")}&mdash;
        {dayjs(bill.end).format("MMM D, YYYY")}
      </Text>
      <SimpleGrid columns={[1, 1, 2, 3]} gap={3} mt={2}>
        {formals.map((f) => (
          <LinkBox
            as={Card}
            bg={bg}
            p={3}
            borderRadius="md"
            key={f.id}
            _hover={{ shadow: "md" }}
            transition="box-shadow 0.2s"
          >
            <LinkOverlay
              as={Link}
              to={`/admin/formals/${f.id}`}
              target="_blank"
            >
              <Heading as="h5" size="sm">
                {f.name} <Icon as={FaExternalLinkAlt} ml={1} boxSize={3} />
              </Heading>
            </LinkOverlay>
            <Text fontSize="sm">{dayjs(f.dateTime).calendar()}</Text>
          </LinkBox>
        ))}
      </SimpleGrid>
      <Button
        colorScheme="brand"
        rightIcon={<FaArrowRight />}
        isLoading={mutation.isLoading}
        onClick={() => mutation.mutate(formals.map((f) => f.id))}
        mt={3}
      >
        Add formals to bill
      </Button>
    </>
  );
}

export function BillFormalsList() {
  const bill = useContext(BillContext);
  const bg = useColorModeValue("gray.100", "gray.600");
  const mutation = useRemoveBillFormal(bill.id);
  if (bill.formals?.length === 0) {
    return <BillFormalsPreview />;
  }
  return (
    <>
      <AddFormalButton />
      <SimpleGrid columns={[1, 1, 2, 3]} gap={3} mt={4}>
        {bill.formals?.map((f) => (
          <Card
            flexDir="row"
            justify="space-between"
            bg={bg}
            p={3}
            borderRadius="md"
            key={f.id}
          >
            <VStack align="start">
              <Heading as="h4" size="sm">
                {f.name}
              </Heading>
              <Text fontSize="sm">{dayjs(f.dateTime).calendar()}</Text>
              <Button
                as={Link}
                variant="outline"
                target="_blank"
                size="xs"
                to={`/admin/formals/${f.id}`}
                rightIcon={<Icon as={FaArrowRight} />}
              >
                More Info
              </Button>
            </VStack>
            <CloseButton
              aria-label="Remove"
              size="sm"
              onClick={() => mutation.mutate(f.id)}
            />
          </Card>
        ))}
      </SimpleGrid>
    </>
  );
}

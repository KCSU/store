import {
  Button,
  CloseButton,
  Heading,
  Icon,
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
import { FaArrowRight, FaPlus } from "react-icons/fa";
import { FormalRadioGroup } from "./FormalRadioGroup";

function AddFormalButton() {
  const bill = useContext(BillContext);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const exclude = useMemo(() => bill.formals?.map((f) => f.id), [bill]);
  const [id, setId] = useState("");
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
            <FormalRadioGroup exclude={exclude} value={id} onChange={setId} />
          </ModalBody>
          <ModalFooter>
            <Button variant="ghost" mr={3} onClick={onClose}>
              Cancel
            </Button>
            <Button colorScheme="brand" onClick={() => {
              onClose();
              setId("");
            }}>
              Add
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

export function BillFormalsList() {
  const bill = useContext(BillContext);
  const bg = useColorModeValue("gray.100", "gray.600");
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
          >
            <VStack align="start">
              <Heading as="h4" size="sm">
                {f.name}
              </Heading>
              <Text fontSize="sm">{dayjs(f.dateTime).calendar()}</Text>
              <Button
                as={Link}
                variant="outline"
                size="xs"
                to={`/admin/formals/${f.id}`}
                rightIcon={<Icon as={FaArrowRight} />}
              >
                More Info
              </Button>
            </VStack>
            <CloseButton aria-label="Remove" size="sm" />
          </Card>
        ))}
      </SimpleGrid>
    </>
  );
}

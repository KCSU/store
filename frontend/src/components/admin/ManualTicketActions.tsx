import {
  Button,
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Icon,
  IconButton,
  Input,
  InputGroup,
  InputRightAddon,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Select,
  useDisclosure,
  VStack,
} from "@chakra-ui/react";
import { Formik, Form, Field, FieldProps } from "formik";
import { FaPen, FaTrashAlt } from "react-icons/fa";
import { useCancelManualTicket } from "../../hooks/admin/useCancelManualTicket";
import { useEditManualTicket } from "../../hooks/admin/useEditManualTicket";
import { ManualTicket, ManualTicketType } from "../../model/ManualTicket";
import { TicketOptionsSelect } from "../tickets/TicketOptionsSelect";

interface ManualTicketProps {
  ticket: ManualTicket;
}

function EditManualTicketButton({ ticket }: ManualTicketProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const mutation = useEditManualTicket(ticket.id);
  return (
    <>
      <IconButton
        variant="ghost"
        size="xs"
        onClick={onOpen}
        aria-label="Edit"
        icon={<Icon as={FaPen} />}
      ></IconButton>
      <Formik
        initialValues={{
          name: ticket.name,
          crsid: ticket.email.split("@")[0],
          option: ticket.option,
          justification: ticket.justification,
          type: ticket.type,
        }}
        onSubmit={async (values, form) => {
          await mutation.mutateAsync({
            name: values.name,
            email: values.crsid + "@cam.ac.uk",
            option: values.option,
            justification: values.justification,
            type: values.type as ManualTicketType,
          });
          onClose();
          // OnClose?
        }}
      >
        {(props) => (
          <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
              <ModalHeader>Edit Special Ticket</ModalHeader>
              <ModalCloseButton />
              <ModalBody>
                <Form>
                  <VStack gap={2}>
                    <Field name="crsid">
                      {({ field, form }: FieldProps) => (
                        <FormControl
                          isInvalid={
                            !!(form.errors.crsid && form.touched.crsid)
                          }
                        >
                          <FormLabel htmlFor="crsid">CRSID</FormLabel>
                          <InputGroup>
                            <Input {...field} type="text" placeholder="crsid" />
                            <InputRightAddon>@cam.ac.uk</InputRightAddon>
                          </InputGroup>
                          <FormErrorMessage>
                            {form.errors.crsid}
                          </FormErrorMessage>
                        </FormControl>
                      )}
                    </Field>
                    <Field name="name">
                      {({ field, form }: FieldProps) => (
                        <FormControl
                          isInvalid={!!(form.errors.name && form.touched.name)}
                        >
                          <FormLabel htmlFor="name">Name</FormLabel>
                          <Input {...field} type="text" placeholder="Name" />
                          <FormErrorMessage>
                            {form.errors.name}
                          </FormErrorMessage>
                        </FormControl>
                      )}
                    </Field>
                    <Field name="option">
                      {({ field, form }: FieldProps) => (
                        <FormControl
                          isInvalid={
                            !!(form.errors.option && form.touched.option)
                          }
                        >
                          <FormLabel htmlFor="option">Option</FormLabel>
                          <TicketOptionsSelect
                            {...field}
                            onChange={(value) =>
                              form.setFieldValue("option", value)
                            }
                          />
                          <FormErrorMessage>
                            {form.errors.option}
                          </FormErrorMessage>
                        </FormControl>
                      )}
                    </Field>
                    <Field name="justification">
                      {({ field, form }: FieldProps) => (
                        <FormControl
                          isInvalid={
                            !!(
                              form.errors.justification &&
                              form.touched.justification
                            )
                          }
                        >
                          <FormLabel htmlFor="justification">
                            Justification
                          </FormLabel>
                          <Input
                            {...field}
                            type="text"
                            placeholder="Justification"
                          />
                          <FormErrorMessage>
                            {form.errors.justification}
                          </FormErrorMessage>
                        </FormControl>
                      )}
                    </Field>
                    <Field name="type">
                      {({ field, form }: FieldProps) => (
                        <FormControl
                          isInvalid={!!(form.errors.type && form.touched.type)}
                        >
                          <FormLabel htmlFor="type">Type</FormLabel>
                          <Select {...field}>
                            <option value="standard">King's</option>
                            <option value="guest">Guest</option>
                            <option value="ents">Ents</option>
                            <option value="complimentary">Complimentary</option>
                          </Select>
                          <FormErrorMessage>
                            {form.errors.type}
                          </FormErrorMessage>
                        </FormControl>
                      )}
                    </Field>
                  </VStack>
                </Form>
              </ModalBody>
              <ModalFooter>
                <Button
                  colorScheme="brand"
                  mr={3}
                  onClick={props.submitForm}
                  isLoading={props.isSubmitting}
                >
                  Save
                </Button>
                <Button variant="ghost" onClick={onClose}>
                  Cancel
                </Button>
              </ModalFooter>
            </ModalContent>
          </Modal>
        )}
      </Formik>
    </>
  );
}

function CancelManualTicketButton({ ticket }: ManualTicketProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const mutation = useCancelManualTicket(ticket.id);
  return (
    <>
      <IconButton
        variant="ghost"
        colorScheme="red"
        size="xs"
        onClick={onOpen}
        aria-label="Delete"
        icon={<Icon as={FaTrashAlt} />}
      ></IconButton>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Cancel Ticket</ModalHeader>
          <ModalCloseButton />
          <ModalBody>Are you sure you want to cancel this ticket?</ModalBody>
          <ModalFooter>
            <Button variant="ghost" onClick={onClose}>
              Close
            </Button>
            <Button
              colorScheme="red"
              isLoading={mutation.isLoading}
              onClick={async () => {
                await mutation.mutateAsync();
                onClose();
              }}
            >
              Cancel Ticket
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

interface ManualTicketActionsProps {
  ticket: ManualTicket;
  canDelete: boolean;
  canWrite: boolean;
}

export function ManualTicketActions({
  canDelete,
  canWrite,
  ticket,
}: ManualTicketActionsProps) {
  return (
    <Flex align="center" gap={2}>
      {canWrite && <EditManualTicketButton ticket={ticket} />}
      {canDelete && <CancelManualTicketButton ticket={ticket} />}
    </Flex>
  );
}

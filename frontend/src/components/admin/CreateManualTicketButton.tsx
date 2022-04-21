import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Icon,
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
import { Field, FieldProps, Form, Formik } from "formik";
import { FaPlus } from "react-icons/fa";
import { useCreateManualTicket } from "../../hooks/admin/useCreateManualTicket";
import { Formal } from "../../model/Formal";
import { ManualTicketType } from "../../model/ManualTicket";
import { TicketOptionsSelect } from "../tickets/TicketOptionsSelect";

interface FormalProps {
  formal: Formal;
}

export function CreateManualTicketButton({ formal }: FormalProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const mutation = useCreateManualTicket();
  return (
    <>
      <Button
        size="sm"
        colorScheme="brand"
        leftIcon={<Icon as={FaPlus} />}
        onClick={onOpen}
      >
        Add Ticket
      </Button>
      <Formik
        initialValues={{
          name: "",
          crsid: "",
          option: "Normal",
          justification: "",
          type: "standard",
        }}
        onSubmit={async (values, form) => {
          await mutation.mutateAsync({
            name: values.name,
            email: values.crsid + "@cam.ac.uk",
            option: values.option,
            justification: values.justification,
            type: values.type as ManualTicketType,
            formalId: formal.id,
          });
          form.resetForm();
          onClose();
          // OnClose?
        }}
      >
        {(props) => (
          <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
              <ModalHeader>Create Special Ticket</ModalHeader>
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
                  Create
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

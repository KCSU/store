import {
  Button,
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Icon,
  Input,
  InputGroup,
  InputLeftAddon,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  SimpleGrid,
  Switch,
  Textarea,
  VStack,
} from "@chakra-ui/react";
import { Formik, Form, Field, FieldProps } from "formik";
import { useContext } from "react";
import { FaSave } from "react-icons/fa";
import { useEditFormal } from "../../hooks/admin/useEditFormal";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Formal, FormalContext } from "../../model/Formal";
import DatePicker from "../utility/DatePicker";

interface FormalDetailsFormProps {
  formal: Formal;
  submitIcon?: React.ReactElement;
  isDisabled?: boolean;
  onSubmit: (values: Formal) => void | Promise<any>;
}

const FormalDetailsForm: React.FC<FormalDetailsFormProps> = ({
  formal,
  onSubmit,
  isDisabled = false,
  children,
  submitIcon,
}) => {
  return (
    <Formik initialValues={formal} onSubmit={onSubmit}>
      {/* TODO: VALIDATION */}
      {(props) => (
        <Form>
          <VStack gap={2}>
            <Field name="name">
              {({ field, form }: FieldProps) => (
                <FormControl
                  isDisabled={isDisabled}
                  isInvalid={!!(form.errors.name && form.touched.name)}
                >
                  <FormLabel htmlFor="name">Name</FormLabel>
                  <Input {...field} id="name" placeholder="Formal Name" />
                  <FormErrorMessage>{form.errors.name}</FormErrorMessage>
                </FormControl>
              )}
            </Field>
            <Field name="menu">
              {/* TODO: make this rich text */}
              {({ field, form }: FieldProps) => (
                <FormControl
                  isDisabled={isDisabled}
                  isInvalid={!!(form.errors.menu && form.touched.menu)}
                >
                  <FormLabel htmlFor="menu">Menu</FormLabel>
                  <Textarea {...field} id="menu" placeholder="Menu" />
                  <FormErrorMessage>{form.errors.menu}</FormErrorMessage>
                </FormControl>
              )}
            </Field>
            <SimpleGrid columns={[1, null, 2]} alignSelf="stretch" gap={3}>
              <Field name="tickets">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={!!(form.errors.tickets && form.touched.tickets)}
                  >
                    <FormLabel htmlFor="tickets">King's Tickets</FormLabel>
                    <NumberInput
                      {...field}
                      id="tickets"
                      onChange={(_, val) => form.setFieldValue(field.name, val)}
                    >
                      <NumberInputField />
                      <NumberInputStepper>
                        <NumberIncrementStepper />
                        <NumberDecrementStepper />
                      </NumberInputStepper>
                    </NumberInput>
                    <FormErrorMessage>{form.errors.tickets}</FormErrorMessage>
                  </FormControl>
                )}
              </Field>
              <Field name="guestTickets">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={
                      !!(form.errors.guestTickets && form.touched.guestTickets)
                    }
                  >
                    <FormLabel htmlFor="guestTickets">Guest Tickets</FormLabel>
                    <NumberInput
                      {...field}
                      id="guestTickets"
                      onChange={(_, val) => form.setFieldValue(field.name, val)}
                    >
                      <NumberInputField />
                      <NumberInputStepper>
                        <NumberIncrementStepper />
                        <NumberDecrementStepper />
                      </NumberInputStepper>
                    </NumberInput>
                    <FormErrorMessage>
                      {form.errors.guestTickets}
                    </FormErrorMessage>
                  </FormControl>
                )}
              </Field>
              <Field name="price">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={!!(form.errors.price && form.touched.price)}
                  >
                    <FormLabel htmlFor="price">Ticket Price</FormLabel>
                    <InputGroup>
                      <InputLeftAddon>??</InputLeftAddon>
                      <NumberInput
                        width="100%"
                        {...field}
                        precision={2}
                        id="price"
                        onChange={(_, val) =>
                          form.setFieldValue(field.name, val)
                        }
                      >
                        <NumberInputField borderLeftRadius={0} />
                        <NumberInputStepper>
                          <NumberIncrementStepper />
                          <NumberDecrementStepper />
                        </NumberInputStepper>
                      </NumberInput>
                    </InputGroup>
                    <FormErrorMessage>{form.errors.price}</FormErrorMessage>
                  </FormControl>
                )}
              </Field>
              <Field name="guestPrice">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={
                      !!(form.errors.guestPrice && form.touched.guestPrice)
                    }
                  >
                    <FormLabel htmlFor="guestPrice">Guest Price</FormLabel>
                    <InputGroup>
                      <InputLeftAddon>??</InputLeftAddon>
                      <NumberInput
                        width="100%"
                        {...field}
                        precision={2}
                        id="guestPrice"
                        onChange={(_, val) =>
                          form.setFieldValue(field.name, val)
                        }
                      >
                        <NumberInputField borderLeftRadius={0} />
                        <NumberInputStepper>
                          <NumberIncrementStepper />
                          <NumberDecrementStepper />
                        </NumberInputStepper>
                      </NumberInput>
                    </InputGroup>
                    <FormErrorMessage>
                      {form.errors.guestPrice}
                    </FormErrorMessage>
                  </FormControl>
                )}
              </Field>
              <Field name="guestLimit">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={
                      !!(form.errors.guestLimit && form.touched.guestLimit)
                    }
                  >
                    <FormLabel htmlFor="guestLimit">Guest Limit</FormLabel>
                    <NumberInput
                      {...field}
                      id="guestLimit"
                      onChange={(_, val) => form.setFieldValue(field.name, val)}
                    >
                      <NumberInputField />
                      <NumberInputStepper>
                        <NumberIncrementStepper />
                        <NumberDecrementStepper />
                      </NumberInputStepper>
                    </NumberInput>
                    <FormErrorMessage>
                      {form.errors.guestLimit}
                    </FormErrorMessage>
                  </FormControl>
                )}
              </Field>
              <Field name="dateTime">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={
                      !!(form.errors.dateTime && form.touched.dateTime)
                    }
                  >
                    <FormLabel htmlFor="dateTime">Formal Start Time</FormLabel>
                    <DatePicker
                      {...field}
                      disabled={isDisabled}
                      selectedDate={field.value}
                      id="dateTime"
                      onChange={(val) => form.setFieldValue(field.name, val)}
                      showPopperArrow
                      showTimeSelect
                      dateFormat="MMMM d, yyyy HH:mm"
                      timeFormat="HH:mm"
                    />
                    <FormErrorMessage>{form.errors.dateTime}</FormErrorMessage>
                  </FormControl>
                )}
              </Field>
              <Field name="saleStart">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={
                      !!(form.errors.saleStart && form.touched.saleStart)
                    }
                  >
                    <FormLabel htmlFor="saleStart">Sale Start Time</FormLabel>
                    <DatePicker
                      {...field}
                      disabled={isDisabled}
                      selectedDate={field.value}
                      id="saleStart"
                      onChange={(val) => form.setFieldValue(field.name, val)}
                      showPopperArrow
                      showTimeSelect
                      dateFormat="MMMM d, yyyy HH:mm"
                      timeFormat="HH:mm"
                    />
                    <FormErrorMessage>{form.errors.saleStart}</FormErrorMessage>
                  </FormControl>
                )}
              </Field>
              <Field name="saleEnd">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={!!(form.errors.saleEnd && form.touched.saleEnd)}
                  >
                    <FormLabel htmlFor="saleEnd">Sale End Time</FormLabel>
                    <DatePicker
                      {...field}
                      disabled={isDisabled}
                      selectedDate={field.value}
                      id="saleEnd"
                      onChange={(val) => form.setFieldValue(field.name, val)}
                      showPopperArrow
                      showTimeSelect
                      dateFormat="MMMM d, yyyy HH:mm"
                      timeFormat="HH:mm"
                    />
                    <FormErrorMessage>{form.errors.saleEnd}</FormErrorMessage>
                  </FormControl>
                )}
              </Field>
            </SimpleGrid>
            <Field name="isVisible">
              {({ field, form }: FieldProps) => (
                <FormControl
                  isInvalid={!!(form.errors.isVisible && form.touched.isVisible)}
                  as={Flex}
                  mt={2}
                  alignItems="center"
                >
                  <FormLabel mb={0}>Visible:</FormLabel>
                  <Switch
                    colorScheme="brand"
                    isChecked={field.value}
                    id="isVisible"
                    onChange={e => form.setFieldValue(field.name, e.target.checked)}
                  />
                </FormControl>
              )}
            </Field>
            <Field name="hasGuestList">
              {({ field, form }: FieldProps) => (
                <FormControl
                  isInvalid={!!(form.errors.hasGuestList && form.touched.hasGuestList)}
                  as={Flex}
                  mt={2}
                  alignItems="center"
                >
                  <FormLabel mb={0}>Public Guest List:</FormLabel>
                  <Switch
                    colorScheme="brand"
                    isChecked={field.value}
                    id="hasGuestList"
                    onChange={e => form.setFieldValue(field.name, e.target.checked)}
                  />
                </FormControl>
              )}
            </Field>
            {!isDisabled && (
              <Button
                colorScheme="brand"
                alignSelf="start"
                leftIcon={submitIcon}
                isLoading={props.isSubmitting}
                onClick={props.submitForm}
              >
                {children}
              </Button>
            )}
          </VStack>
        </Form>
      )}
    </Formik>
  );
};

export function EditFormalForm() {
  const formal = useContext(FormalContext);
  const mutation = useEditFormal(formal.id);
  const canWrite = useHasPermission("formals", "write");
  return (
    <FormalDetailsForm
      isDisabled={!canWrite}
      formal={formal}
      onSubmit={async (values) => {
        await mutation.mutateAsync(values);
      }}
      submitIcon={<Icon as={FaSave} />}
    >
      Save Changes
    </FormalDetailsForm>
  );
}

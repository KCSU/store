import {
  Button,
  Checkbox,
  CheckboxGroup,
  Container,
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  Input,
  InputGroup,
  InputLeftAddon,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  SimpleGrid,
  Stack,
  Switch,
  Textarea,
  VStack,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import { Formik, Form, Field, FieldProps } from "formik";
import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import DatePicker from "../../components/utility/DatePicker";
import { CreateFormalDto, useCreateFormal } from "../../hooks/admin/useCreateFormal";
import { useGroups } from "../../hooks/admin/useGroups";
import { Group } from "../../model/Group";

interface CreateFormalForm {
  formal: CreateFormalDto;
  submitIcon?: React.ReactElement;
  onSubmit: (values: CreateFormalDto) => void | Promise<any>;
  availableGroups?: Group[];
}

const CreateFormalForm: React.FC<CreateFormalForm> = ({
  formal,
  onSubmit,
  children,
  submitIcon,
  availableGroups = [],
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
                    isInvalid={!!(form.errors.price && form.touched.price)}
                  >
                    <FormLabel htmlFor="price">Ticket Price</FormLabel>
                    <InputGroup>
                      <InputLeftAddon>£</InputLeftAddon>
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
                    isInvalid={
                      !!(form.errors.guestPrice && form.touched.guestPrice)
                    }
                  >
                    <FormLabel htmlFor="guestPrice">Guest Price</FormLabel>
                    <InputGroup>
                      <InputLeftAddon>£</InputLeftAddon>
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
                    isInvalid={
                      !!(form.errors.dateTime && form.touched.dateTime)
                    }
                  >
                    <FormLabel htmlFor="dateTime">Formal Start Time</FormLabel>
                    <DatePicker
                      {...field}
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
                    isInvalid={
                      !!(form.errors.saleStart && form.touched.saleStart)
                    }
                  >
                    <FormLabel htmlFor="saleStart">Sale Start Time</FormLabel>
                    <DatePicker
                      {...field}
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
                    isInvalid={!!(form.errors.saleEnd && form.touched.saleEnd)}
                  >
                    <FormLabel htmlFor="saleEnd">Sale End Time</FormLabel>
                    <DatePicker
                      {...field}
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
            <Field name="groups">
              {({ field, form }: FieldProps) => (
                <FormControl
                  isInvalid={!!(form.errors.groups && form.touched.groups)}
                >
                  <FormLabel>Groups</FormLabel>
                  <CheckboxGroup
                    defaultValue={field.value.map((g: Group) => g.id)}
                    onChange={(val) => {
                      const groups = availableGroups.filter((g) =>
                        val.includes(g.id)
                      ).map(g => g.id);
                      form.setFieldValue(field.name, groups);
                    }}
                    colorScheme="brand"
                  >
                    <Stack>
                      {availableGroups.map((g) => (
                        <Checkbox value={g.id} key={g.id}>
                          {g.name}
                        </Checkbox>
                      ))}
                    </Stack>
                  </CheckboxGroup>
                </FormControl>
              )}
            </Field>
            <Button
              colorScheme="brand"
              alignSelf="start"
              leftIcon={submitIcon}
              isLoading={props.isSubmitting}
              onClick={props.submitForm}
            >
              {children}
            </Button>
          </VStack>
        </Form>
      )}
    </Formik>
  );
};

export function AdminCreateFormalView() {
  const { data: groups, isLoadingError } = useGroups();
  const mutation = useCreateFormal();
  const navigate = useNavigate();
  const defaultFormal = useMemo(() => {
    const currentDate = dayjs().startOf("day").toDate();
    const f: CreateFormalDto = {
      name: "",
      menu: "",
      price: 0,
      guestPrice: 0,
      guestLimit: 0,
      tickets: 0,
      ticketsRemaining: 0,
      guestTickets: 0,
      guestTicketsRemaining: 0,
      saleStart: currentDate,
      saleEnd: currentDate,
      dateTime: currentDate,
      hasGuestList: true,
      groups: [],
    };
    return f;
  }, []);
  return (
    <Container maxW="container.md" p={0}>
      <BackButton to="/admin/formals">Back Home</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={4}>
          Create a Formal
        </Heading>
        <CreateFormalForm
          formal={defaultFormal}
          availableGroups={groups}
          onSubmit={async (values: CreateFormalDto) => {
            await mutation.mutateAsync(values);
            navigate('/admin/formals');
          }}
        >
          Create Formal
        </CreateFormalForm>
      </Card>
    </Container>
  );
}

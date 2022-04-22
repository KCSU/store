import {
  VStack,
  FormControl,
  FormLabel,
  Input,
  FormErrorMessage,
  Button,
  SimpleGrid,
} from "@chakra-ui/react";
import { Formik, Form, Field, FieldProps } from "formik";
import { Bill } from "../../model/Bill";
import DatePicker from "../utility/DatePicker";

interface BillDetailsFormProps {
  initialValues: Bill;
  onSubmit: (values: Bill) => void | Promise<any>;
  isDisabled?: boolean;
}

export const BillDetailsForm: React.FC<BillDetailsFormProps> = ({
  initialValues,
  onSubmit,
  isDisabled = false,
  children,
}) => {
  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit}>
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
                  <Input {...field} id="name" placeholder="Bill Name" />
                  <FormErrorMessage>{form.errors.name}</FormErrorMessage>
                </FormControl>
              )}
            </Field>
            <SimpleGrid columns={[1, null, 2]} alignSelf="stretch" gap={3}>
              <Field name="start">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={!!(form.errors.start && form.touched.start)}
                  >
                    <FormLabel htmlFor="start">Start Date</FormLabel>
                    <DatePicker
                      {...field}
                      id="start"
                      disabled={isDisabled}
                      selectedDate={field.value}
                      onChange={(date) => form.setFieldValue(field.name, date)}
                      // showPopperArrow
                      dateFormat="MMMM d, yyyy"
                    />
                    <FormErrorMessage>{form.errors.start}</FormErrorMessage>
                  </FormControl>
                )}
              </Field>
              <Field name="end">
                {({ field, form }: FieldProps) => (
                  <FormControl
                    isDisabled={isDisabled}
                    isInvalid={!!(form.errors.end && form.touched.end)}
                  >
                    <FormLabel htmlFor="end">End Date</FormLabel>
                    <DatePicker
                      {...field}
                      id="end"
                      disabled={isDisabled}
                      selectedDate={field.value}
                      onChange={(date) => form.setFieldValue(field.name, date)}
                      // showPopperArrow
                      dateFormat="MMMM d, yyyy"
                    />
                    <FormErrorMessage>{form.errors.end}</FormErrorMessage>
                  </FormControl>
                )}
              </Field>
            </SimpleGrid>
            {!isDisabled && (
              <Button
                colorScheme="brand"
                alignSelf="start"
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

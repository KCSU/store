import {
  VStack,
  FormControl,
  FormLabel,
  Input,
  FormErrorMessage,
  Button,
  Icon,
} from "@chakra-ui/react";
import { Select } from "chakra-react-select";
import { Formik, Form, Field, FieldProps } from "formik";
import { FaSave } from "react-icons/fa";
import { Group, groupType } from "../../model/Group";

interface GroupDetailsFormProps {
  group: Group;
  onSubmit: (values: Group) => void | Promise<any>;
  isDisabled?: boolean;
}

export const GroupDetailsForm: React.FC<GroupDetailsFormProps> = ({
  group,
  onSubmit,
  isDisabled = false,
  children,
}) => {
  return (
    <Formik initialValues={group} onSubmit={onSubmit}>
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
                  <Input {...field} id="name" placeholder="Group Name" />
                  <FormErrorMessage>{form.errors.name}</FormErrorMessage>
                </FormControl>
              )}
            </Field>
            <Field name="type">
              {({ field, form }: FieldProps) => (
                <FormControl
                  isDisabled={isDisabled}
                  isInvalid={!!(form.errors.type && form.touched.type)}
                >
                  <FormLabel htmlFor="type">Lookup Type</FormLabel>
                  <Select
                    value={{
                      value: field.value,
                      label: groupType(field.value),
                    }}
                    options={[
                      { value: "inst", label: "Institution" },
                      { value: "group", label: "Group" },
                      { value: "manual", label: "Manual" },
                    ]}
                    selectedOptionColor="brand"
                    onChange={(option) =>
                      form.setFieldValue(field.name, option?.value)
                    }
                  ></Select>
                </FormControl>
              )}
            </Field>
            <Field name="lookup">
              {({ field, form }: FieldProps) => (
                <FormControl
                  isDisabled={isDisabled}
                  isInvalid={!!(form.errors.lookup && form.touched.lookup)}
                >
                  <FormLabel htmlFor="lookup">Lookup Query</FormLabel>
                  <Input {...field} id="lookup" placeholder="Query" />
                  <FormErrorMessage>{form.errors.lookup}</FormErrorMessage>
                </FormControl>
              )}
            </Field>
            {!isDisabled && (
              <Button
                colorScheme="brand"
                alignSelf="start"
                leftIcon={<Icon as={FaSave} />}
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

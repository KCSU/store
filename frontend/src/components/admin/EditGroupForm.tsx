import {
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import { Select } from "chakra-react-select";
import { Field, FieldProps, Form, Formik } from "formik";
import { Group, groupType } from "../../model/Group";

interface GroupProps {
  group: Group;
}

export function EditGroupForm({ group }: GroupProps) {
  // const
  return (
    <Formik initialValues={group} onSubmit={() => {}}>
      {(props) => (
        <Form>
          <VStack gap={2}>
            <Field name="name">
              {({ field, form }: FieldProps) => (
                <FormControl
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
                  isInvalid={!!(form.errors.type && form.touched.type)}
                >
                  <FormLabel htmlFor="type">Lookup Type</FormLabel>
                  <Select
                    value={{
                      value: field.value,
                      label: groupType(field.value)
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
                  isInvalid={!!(form.errors.lookup && form.touched.lookup)}
                >
                  <FormLabel htmlFor="lookup">Lookup Query</FormLabel>
                  <Input {...field} id="lookup" placeholder="Query" />
                  <FormErrorMessage>{form.errors.lookup}</FormErrorMessage>
                </FormControl>
              )}
            </Field>
          </VStack>
        </Form>
      )}
    </Formik>
  );
}

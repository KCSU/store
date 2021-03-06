import {
  Button,
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  Icon,
  IconButton,
  Input,
  Select,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  useColorModeValue,
} from "@chakra-ui/react";
import { Field, FieldProps, Form, Formik } from "formik";
import { useMemo } from "react";
import { FaPlus, FaTrashAlt } from "react-icons/fa";
import { CellProps, Column, useTable } from "react-table";
import { useAddUserRole } from "../../hooks/admin/useAddUserRole";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { useRemoveUserRole } from "../../hooks/admin/useRemoveUserRole";
import { useRoles } from "../../hooks/admin/useRoles";
import { useUserRoles } from "../../hooks/admin/useUserRoles";
import { UserRole } from "../../model/UserRole";

interface UserRolesTableProps {
  userRoles: UserRole[];
}

function UserRolesTable({ userRoles }: UserRolesTableProps) {
  const canWrite = useHasPermission("roles", "write");
  const columns = useMemo<Column<UserRole>[]>(() => {
    const cols: Column<UserRole>[] = [
      {
        accessor: "userName",
        Header: "Name",
      },
      {
        accessor: "userEmail",
        Header: "Email",
      },
      {
        accessor: "roleName",
        Header: "Role",
      },
    ];
    if (canWrite) {
      cols.push({
        Header: "Actions",
        Cell({ row: { original } }: CellProps<UserRole>) {
          const mutation = useRemoveUserRole();
          return (
            <IconButton
              aria-label="Revoke"
              size="xs"
              colorScheme="red"
              isLoading={mutation.isLoading}
              variant="ghost"
              onClick={() => {
                mutation.mutate({
                  roleId: original.roleId,
                  email: original.userEmail,
                });
              }}
            >
              <Icon as={FaTrashAlt} />
            </IconButton>
          );
        },
      });
    }
    return cols;
  }, []);
  const { getTableProps, getTableBodyProps, headerGroups, prepareRow, rows } =
    useTable({
      columns,
      data: userRoles,
    });
  const background = useColorModeValue("white", "gray.750");
  return (
    <>
      <Table variant="striped" size="sm" {...getTableProps()}>
        <Thead>
          {headerGroups.map((headerGroup) => (
            <Tr {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map((column) => (
                <Th {...column.getHeaderProps()} p={1}>
                  {column.render("Header")}
                </Th>
              ))}
            </Tr>
          ))}
        </Thead>
        <Tbody {...getTableBodyProps()} bg={background}>
          {rows.map((row) => {
            prepareRow(row);
            return (
              <Tr {...row.getRowProps()}>
                {row.cells.map((cell) => (
                  <Td {...cell.getCellProps()} p={1}>
                    {cell.render("Cell")}
                  </Td>
                ))}
              </Tr>
            );
          })}
        </Tbody>
      </Table>
      {canWrite && (
        <>
          <Heading as="h4" size="sm">
            Add User
          </Heading>
          <AddUserRoleForm />
        </>
      )}
    </>
  );
}

function AddUserRoleForm() {
  const { data, isLoading, isError } = useRoles();
  const mutation = useAddUserRole();
  if (isLoading || isError || !data) {
    return null;
  }
  return (
    <Formik
      initialValues={{
        email: "",
        roleId: "",
      }}
      onSubmit={async (values, form) => {
        await mutation.mutateAsync(values);
        form.resetForm();
      }}
    >
      {(props) => (
        <Form>
          <Flex gap={3}>
            <Field name="email">
              {({ field, form }: FieldProps) => (
                <FormControl
                  isInvalid={!!(form.errors.name && form.touched.name)}
                >
                  <Input
                    size="sm"
                    {...field}
                    type="email"
                    placeholder="Email"
                  />
                  <FormErrorMessage>{form.errors.name}</FormErrorMessage>
                </FormControl>
              )}
            </Field>
            <Field name="roleId">
              {({ field, form }: FieldProps) => (
                <FormControl
                  isInvalid={!!(form.errors.roleId && form.touched.roleId)}
                >
                  <Select
                    {...field}
                    placeholder="Choose Role"
                    size="sm"
                  >
                    {data.map((role) => (
                      <option key={role.id} value={role.id}>
                        {role.name}
                      </option>
                    ))}
                  </Select>
                </FormControl>
              )}
            </Field>
            <Button
              size="sm"
              colorScheme="brand"
              flexGrow="0"
              flexShrink="0"
              leftIcon={<Icon as={FaPlus} />}
              isLoading={props.isSubmitting}
              onClick={props.submitForm}
            >
              Add
            </Button>
          </Flex>
        </Form>
      )}
    </Formik>
  );
}

export function UserRolesList() {
  const { data, isLoading, isError } = useUserRoles();
  if (isLoading || isError || !data) {
    return <></>;
  }
  return <UserRolesTable userRoles={data} />;
}

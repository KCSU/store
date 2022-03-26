import {
  Button,
  FormControl,
  FormHelperText,
  Icon,
  IconButton,
  Input,
  Table,
  Tbody,
  Td,
  Tfoot,
  Th,
  Thead,
  Tooltip,
  Tr,
  useColorModeValue,
} from "@chakra-ui/react";
import { useMemo, useState } from "react";
import { FaPlus, FaTrashAlt } from "react-icons/fa";
import { CellProps, Column, useTable } from "react-table";
import { useCreatePermission } from "../../hooks/admin/useCreatePermission";
import { useDeletePermission } from "../../hooks/admin/useDeletePermission";
import { Permission } from "../../model/Permission";
import { Role } from "../../model/Role";

interface RoleProps {
  role: Role;
}

export function PermissionsTable({ role }: RoleProps) {
  const inBg = useColorModeValue("white", "gray.600");
  const mutation = useCreatePermission();
  const [resource, setResource] = useState("");
  const [action, setAction] = useState("");
  const resources = [
    'formals', 'tickets', 'groups', 
    'roles', 'permissions', 'billing', '*'
  ];
  const actions = [
    'read', 'write', 'delete', '*'
  ];
  const columns = useMemo<Column<Permission>[]>(
    () => [
      {
        accessor: "resource",
        Header: "Resource",
      },
      {
        accessor: "action",
        Header: "Permission",
      },
      {
        Header: "Actions",
        Cell({row: {original}}: CellProps<Permission>) {
          const mutation = useDeletePermission(original.id);
          return (
            <IconButton
              aria-label="Revoke"
              isLoading={mutation.isLoading}
              onClick={() => mutation.mutate()}
              size="xs"
              colorScheme="red"
              variant="ghost"
            >
              <Icon as={FaTrashAlt} />
            </IconButton>
          );
        },
      },
    ],
    []
  );
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    prepareRow,
    rows,
  } = useTable({
    columns,
    data: role.permissions ?? [],
  });
  const background = useColorModeValue("white", "gray.750");
  return (
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
                <Td {...cell.getCellProps()} p={2} fontFamily="mono">
                  {cell.render("Cell")}
                </Td>
              ))}
            </Tr>
          );
        })}
      </Tbody>
      <Tfoot>
        <Tr>
          <Td p={1} verticalAlign="top" width="50%">
            <FormControl>
            <Input
              size="sm"
              isDisabled={mutation.isLoading}
              placeholder="Resource"
              bg={inBg}
              fontFamily="mono"
              value={resource}
              onChange={(e) => setResource(e.target.value)}
            ></Input>
            <FormHelperText fontSize="xs" fontStyle="italic" mt="0.5">
              {resources.join(', ')}
            </FormHelperText>
            </FormControl>
          </Td>
          <Td p={1} verticalAlign="top" width="50%">
            <FormControl>
            <Input
              size="sm"
              isDisabled={mutation.isLoading}
              placeholder="Permission"
              bg={inBg}
              fontFamily="mono"
              value={action}
              onChange={(e) => setAction(e.target.value)}
            ></Input>
            <FormHelperText fontSize="xs" fontStyle="italic" mt="0.5">
              {actions.join(', ')}
            </FormHelperText>
            </FormControl>
          </Td>
          <Td verticalAlign="top" p={1}>
            <Button
              size="sm"
              isLoading={mutation.isLoading}
              colorScheme="brand"
              leftIcon={<Icon as={FaPlus} />}
              onClick={async () => {
                await mutation.mutateAsync({
                  roleId: role.id,
                  resource,
                  action,
                });
                setResource('');
                setAction('');
              }}
            >
              Add
            </Button>
          </Td>
        </Tr>
      </Tfoot>
    </Table>
  );
}

import {
  Button,
  Icon,
  IconButton,
  Input,
  Table,
  Tbody,
  Td,
  Tfoot,
  Th,
  Thead,
  Tr,
  useColorModeValue,
} from "@chakra-ui/react";
import { useMemo, useState } from "react";
import { FaPlus, FaTrashAlt } from "react-icons/fa";
import { Column, useTable } from "react-table";
import { useCreatePermission } from "../../hooks/admin/useCreatePermission";
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
        Cell() {
          return (
            <IconButton
              aria-label="Revoke"
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
    footerGroups,
  } = useTable({
    columns,
    data: role.permissions ?? [],
  });
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
      <Tbody {...getTableBodyProps()}>
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
          <Th p={1}>
            <Input
              size="sm"
              isDisabled={mutation.isLoading}
              placeholder="Resource"
              bg={inBg}
              fontFamily="mono"
              value={resource}
              onChange={(e) => setResource(e.target.value)}
            ></Input>
          </Th>
          <Th p={1}>
            <Input
              size="sm"
              isDisabled={mutation.isLoading}
              placeholder="Permission"
              bg={inBg}
              fontFamily="mono"
              value={action}
              onChange={(e) => setAction(e.target.value)}
            ></Input>
          </Th>
          <Th p={1}>
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
          </Th>
        </Tr>
      </Tfoot>
    </Table>
  );
}

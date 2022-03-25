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
import { useMemo } from "react";
import { FaPlus, FaTrashAlt } from "react-icons/fa";
import { Column, useTable } from "react-table";
import { Permission } from "../../model/Permission";
import { Role } from "../../model/Role";

interface RoleProps {
  role: Role;
}

export function PermissionsTable({ role }: RoleProps) {
  const inBg = useColorModeValue('white', 'gray.600');
  const columns = useMemo<Column<Permission>[]>(
    () => [
      {
        accessor: "resource",
        Header: "Resource",
        Footer() {
          return <Input size="sm" placeholder="Resource" bg={inBg}></Input>;
        },
      },
      {
        accessor: "action",
        Header: "Permission",
        Footer() {
          return (
            <Input size="sm" placeholder="Permission" bg={inBg}></Input>
          );
        },
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
        Footer() {
          return (
            <Button
              size="sm"
              colorScheme="brand"
              leftIcon={<Icon as={FaPlus} />}
            >
              Add
            </Button>
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
        {footerGroups.map((fg) => (
          <Tr {...fg.getFooterGroupProps()}>
            {fg.headers.map((column) => (
              <Th {...column.getFooterProps()} p={1}>
                {column.render("Footer")}
              </Th>
            ))}
          </Tr>
        ))}
      </Tfoot>
    </Table>
  );
}

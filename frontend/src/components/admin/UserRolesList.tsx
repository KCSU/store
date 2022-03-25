import { Icon, IconButton, Table, Tbody, Td, Th, Thead, Tr } from "@chakra-ui/react";
import { useMemo } from "react";
import { FaTrashAlt } from "react-icons/fa";
import { Column, useTable } from "react-table";
import { useUserRoles } from "../../hooks/admin/useUserRoles";
import { UserRole } from "../../model/UserRole";

interface UserRolesTableProps {
  userRoles: UserRole[];
}

function UserRolesTable({ userRoles }: UserRolesTableProps) {
  const columns = useMemo<Column<UserRole>[]>(
    () => [
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
      {
        Header: "Actions",
        Cell() {
          return <IconButton aria-label="Revoke" size="xs" colorScheme="red" variant="ghost">
            <Icon as={FaTrashAlt}/>
          </IconButton>
        },
      }
    ],
    []
  );
  const { getTableProps, getTableBodyProps, headerGroups, prepareRow, rows } =
    useTable({
      columns,
      data: userRoles,
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
                <Td {...cell.getCellProps()} p={1}>{cell.render("Cell")}</Td>
              ))}
            </Tr>
          );
        })}
      </Tbody>
    </Table>
  );
}

export function UserRolesList() {
  const { data, isLoading, isError } = useUserRoles();
  if (isLoading || isError || !data) {
    return <></>;
  }
  return <UserRolesTable userRoles={data} />;
}

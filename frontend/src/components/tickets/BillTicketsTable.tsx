import {
  useColorModeValue,
  Table,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
  Tfoot,
} from "@chakra-ui/react";
import { useMemo } from "react";
import { Column, useTable } from "react-table";
import { formatMoney } from "../../helpers/formatMoney";
import { Formal } from "../../model/Formal";
import { FormalTicket, Ticket } from "../../model/Ticket";

interface BillTicketsTableProps {
  tickets: FormalTicket[];
}

type BillTicket = Ticket & {
  formal: Formal;
};

export function BillTicketsTable({ tickets }: BillTicketsTableProps) {
  const columns = useMemo<Column<BillTicket>[]>(
    () => [
      {
        Header: "Type",
        accessor: "isGuest",
        Cell: ({ value }) => (value ? "Guest" : "King's"),
        Footer: () => "Total",
      },
      {
        Header: "Formal",
        accessor: "formal.name",
      } as Column<BillTicket>,
      {
        Header: "Meal Option",
        accessor: "option",
      },
      {
        Header: "Price",
        accessor: "formal",
        Cell: ({ value, row: { original } }) =>
          formatMoney(original.isGuest ? value.guestPrice : value.price),
        isNumeric: true,
        Footer: (table) => {
          return formatMoney(table.rows.reduce((acc, row) => {
            return (
              acc +
              (row.original.isGuest
                ? row.original.formal.guestPrice
                : row.original.formal.price)
            );
          }, 0));
        },
      },
    ],
    []
  );
  const data = useMemo(() => {
    const result: BillTicket[] = [];
    tickets.forEach((t) => {
      result.push({
        ...t.ticket,
        formal: t.formal,
      });
      t.guestTickets.forEach((gt) => {
        result.push({
          ...gt,
          formal: t.formal,
        });
      });
    });
    return result;
  }, [tickets]);
  const { 
    getTableBodyProps,
    getTableProps,
    headerGroups,
    rows,
    prepareRow,
    footerGroups
  } = useTable({ columns, data });
  const background = useColorModeValue("white", "gray.750");
  return (
    <Table my={2} variant="striped" size="sm" {...getTableProps()}>
      <Thead>
        {headerGroups.map((headerGroup) => (
          <Tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map((column) => (
              <Th
                {...column.getHeaderProps()}
                isNumeric={(column as any).isNumeric}
              >
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
                <Td
                  {...cell.getCellProps()}
                  isNumeric={(cell.column as any).isNumeric}
                >
                  {cell.render("Cell")}
                </Td>
              ))}
            </Tr>
          );
        })}
      </Tbody>
      <Tfoot>
        {footerGroups.map((footerGroup) => (
          <Tr {...footerGroup.getFooterGroupProps()}>
            {footerGroup.headers.map((column) => (
              <Td
                {...column.getFooterProps()}
                isNumeric={(column as any).isNumeric}
                fontSize="md"
                fontWeight="bold"
              >
                {column.render("Footer")}
              </Td>
            ))}
          </Tr>
        ))}
      </Tfoot>
    </Table>
  );
}

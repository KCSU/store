import {
  Alert,
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  AlertIcon,
  background,
  Button,
  Heading,
  Icon,
  Table,
  Tbody,
  Td,
  Tfoot,
  Th,
  Thead,
  Tr,
  useColorModeValue,
  useDisclosure,
} from "@chakra-ui/react";
import { useContext, useMemo, useRef } from "react";
import { FaTrashAlt } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { Column, TableInstance, useTable } from "react-table";
import { useDeleteFormal } from "../../hooks/admin/useDeleteFormal";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Formal, FormalContext } from "../../model/Formal";

interface Stat {
  type: string;
  count: number;
}

function typeStats(formal: Formal): Stat[] {
  const stats = {
    kings: 0,
    guest: 0,
    kingsManual: 0,
    guestManual: 0,
    complimentary: 0,
    ents: 0,
  };
  formal.ticketSales?.forEach((t) => {
    if (t.isGuest) {
      stats.guest++;
    } else {
      stats.kings++;
    }
  });
  formal.manualTickets?.forEach((t) => {
    switch (t.type) {
      case "standard":
        stats.kingsManual++;
        break;
      case "guest":
        stats.guestManual++;
        break;
      case "complimentary":
        stats.complimentary++;
        break;
      case "ents":
        stats.ents++;
        break;
    }
  });
  return [
    { type: "King's", count: stats.kings },
    { type: "Guests", count: stats.guest },
    { type: "King's (Manual)", count: stats.kingsManual },
    { type: "Guests (Manual)", count: stats.guestManual },
    { type: "Complimentary", count: stats.complimentary },
    { type: "Ents", count: stats.ents },
  ];
}

function mealStats(formal: Formal): Stat[] {
  const stats = new Map<string, number>();
  formal.ticketSales?.forEach(t => {
    const count = stats.get(t.option) || 0;
    stats.set(t.option, count + 1);
  });
  formal.manualTickets?.forEach(t => {
    const count = stats.get(t.option) || 0;
    stats.set(t.option, count + 1);
  });
  return Array.from(stats.entries()).map(([option, count]) => ({
    type: option,
    count,
  }));
}

export function FormalStats() {
  const formal = useContext(FormalContext);
  const canDelete = useHasPermission("formals", "delete");
  const typeColumns = useMemo<Column<Stat>[]>(
    () => [
      {
        Header: "Type",
        accessor: "type",
        Footer: "Total",
        isNumeric: false,
      },
      {
        Header: "Count",
        accessor: "count",
        Footer: (table) => {
          return table.rows.reduce((acc, row) => acc + row.values.count, 0);
        },
        isNumeric: true,
      },
    ],
    []
  );
  const typeData = useMemo(() => typeStats(formal), [formal]);
  const mealColumns = useMemo<Column<Stat>[]>(
    () => [
      {
        Header: "Option",
        accessor: "type",
        Footer: "Total",
        isNumeric: false,
      },
      {
        Header: "Count",
        accessor: "count",
        Footer: (table) => {
          return table.rows.reduce((acc, row) => acc + row.values.count, 0);
        },
        isNumeric: true,
      },
    ],
    []
  );
  const mealData = useMemo(() => mealStats(formal), [formal]);
  const typeTable = useTable({
    columns: typeColumns,
    data: typeData,
  });
  const mealTable = useTable({
    columns: mealColumns,
    data: mealData,
  });
  return (
    <>
      <Heading as="h4" size="md" mb={1}>Ticket Types</Heading>
      <StatTable table={typeTable} />
      <Heading as="h4" size="md" mt={4} mb={1}>Meal Options</Heading>
      <StatTable table={mealTable} />
      {canDelete && <FormalActions formal={formal} />}
    </>
  );
}

interface StatTableProps<T extends object> {
  table: TableInstance<T>;
}

function StatTable<T extends object>({ table }: StatTableProps<T>) {
  const {
    getTableBodyProps,
    getTableProps,
    prepareRow,
    headerGroups,
    footerGroups,
    rows,
  } = table;
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
              <Th
                {...column.getFooterProps()}
                pt={3}
                isNumeric={(column as any).isNumeric}
              >
                {column.render("Footer")}
              </Th>
            ))}
          </Tr>
        ))}
      </Tfoot>
    </Table>
  );
}

interface FormalProps {
  formal: Formal;
}

function FormalActions({ formal }: FormalProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = useRef(null);
  const mutation = useDeleteFormal(formal.id);
  const navigate = useNavigate();
  return (
    <>
      <Heading as="h4" size="md" mt={4}>
        Actions
      </Heading>
      <Alert status="warning" mb={4} mt={2} variant="left-accent">
        <AlertIcon />
        The following actions are potentially destructive! Only proceed if you
        know what you're doing.
      </Alert>
      <Button
        colorScheme="red"
        leftIcon={<Icon as={FaTrashAlt} />}
        onClick={onOpen}
      >
        Delete Formal
      </Button>
      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Formal
            </AlertDialogHeader>
            <AlertDialogBody>
              Are you sure you want to delete this formal? This action cannot be
              undone.
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              <Button
                colorScheme="red"
                onClick={async () => {
                  await mutation.mutateAsync();
                  onClose();
                  navigate("/admin/formals");
                }}
                ml={3}
                isLoading={mutation.isLoading}
              >
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
}

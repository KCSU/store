import {
  Table,
  Thead,
  Tr,
  Th,
  Tbody,
  Text,
  Td,
  Link,
  Icon,
  Flex,
  HStack,
  IconButton,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  Select,
  Tooltip,
  useBreakpointValue,
} from "@chakra-ui/react";
import _ from "lodash";
import { useMemo } from "react";
import {
  FaAngleDoubleLeft,
  FaAngleDoubleRight,
  FaAngleLeft,
  FaAngleRight,
  FaExternalLinkAlt,
} from "react-icons/fa";
import {
  CellProps,
  Column,
  usePagination,
  useSortBy,
  useTable,
} from "react-table";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Formal } from "../../model/Formal";
import { AdminTicket } from "../../model/Ticket";

interface FormalProps {
  formal: Formal;
}

export function FormalTicketsList({ formal }: FormalProps) {
  const canWrite = useHasPermission("tickets", "write");
  const canDelete = useHasPermission("tickets", "delete");
  const columns: Column<AdminTicket>[] = useMemo(() => {
    const base: Column<AdminTicket>[] = [
      {
        accessor: "userName",
        Header: "Name",
      },
      {
        accessor: "userEmail",
        Header: "Crsid",
        Cell: ({ value }) => {
          const crsid = value.split("@")[0];
          return (
            <Link href={`mailto:${value}`} isExternal>
              {crsid} <Icon boxSize={3} as={FaExternalLinkAlt} />
            </Link>
          );
        },
      },
      {
        accessor: "option",
        Header: "Option",
      },
      {
        accessor: "isGuest",
        Header: "Type",
        Cell: ({ row }: CellProps<AdminTicket>) => {
          return row.original.isGuest ? "Guest" : "King's";
        },
      },
    ];
    // TODO: Actions
    return base;
  }, [canDelete, canWrite]);
  const data = useMemo(
    () => _.orderBy(formal.ticketSales ?? [], ["userName", "isGuest"]),
    [formal]
  );
  const {
    getTableProps,
    getTableBodyProps,
    prepareRow,
    headerGroups,
    gotoPage,
    canPreviousPage,
    canNextPage,
    page,
    previousPage,
    nextPage,
    pageOptions,
    state: { pageIndex, pageSize },
    pageCount,
    setPageSize,
  } = useTable(
    {
      columns,
      data,
    },
    usePagination
  );
  const showOptions = useBreakpointValue({ base: false, lg: true });
  // TODO: Search filter
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
        <Tbody {...getTableBodyProps()}>
          {page.map((row) => {
            prepareRow(row);
            return (
              <Tr {...row.getRowProps()}>
                {row.cells.map((cell) => (
                  <Td {...cell.getCellProps()} p={2}>
                    {cell.render("Cell")}
                  </Td>
                ))}
              </Tr>
            );
          })}
        </Tbody>
      </Table>
      <Flex justifyContent="space-between" alignItems="center" mt={2}>
        <Flex>
          <Tooltip label="First Page">
            <IconButton
              mr={1}
              size="sm"
              onClick={() => gotoPage(0)}
              isDisabled={!canPreviousPage}
              aria-label="first page"
              icon={<Icon as={FaAngleDoubleLeft} h={3} w={3} />}
            />
          </Tooltip>
          <Tooltip label="Previous Page">
            <IconButton
              size="sm"
              aria-label="previous page"
              onClick={previousPage}
              isDisabled={!canPreviousPage}
              icon={<Icon as={FaAngleLeft} h={3} w={3} />}
              mr={4}
            />
          </Tooltip>
        </Flex>
        <Text flexShrink="0" mx={1}>
          Page{" "}
          <Text fontWeight="bold" as="span">
            {pageIndex + 1}
          </Text>{" "}
          of{" "}
          <Text fontWeight="bold" as="span">
            {pageOptions.length}
          </Text>
        </Text>
        {showOptions && (
          <>
            <HStack mx={2}>
              <Text flexShrink="0">Go to page:</Text>{" "}
              <NumberInput
                size="sm"
                maxW={28}
                min={1}
                max={pageOptions.length}
                onChange={(value) => {
                  const page = value ? parseInt(value) - 1 : 0;
                  gotoPage(page);
                }}
                defaultValue={pageIndex + 1}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
            </HStack>

            <Select
              maxW={32}
              size="sm"
              value={pageSize}
              onChange={(e) => {
                setPageSize(Number(e.target.value));
              }}
            >
              {[10, 20, 30, 40, 50].map((pageSize) => (
                <option key={pageSize} value={pageSize}>
                  Show {pageSize}
                </option>
              ))}
            </Select>
          </>
        )}

        <Flex>
          <Tooltip label="Next Page">
            <IconButton
              size="sm"
              onClick={nextPage}
              aria-label="next page"
              isDisabled={!canNextPage}
              icon={<Icon as={FaAngleRight} h={3} w={3} />}
              ml={4}
            />
          </Tooltip>
          <Tooltip label="Last Page">
            <IconButton
              size="sm"
              aria-label="last page"
              onClick={() => gotoPage(pageCount - 1)}
              isDisabled={!canNextPage}
              icon={<Icon as={FaAngleDoubleRight} h={3} w={3} />}
              ml={1}
            />
          </Tooltip>
        </Flex>
      </Flex>
    </>
  );
}

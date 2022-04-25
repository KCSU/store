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
  Input,
  InputGroup,
  InputLeftAddon,
  useColorModeValue,
} from "@chakra-ui/react";
import _ from "lodash";
import { useContext, useEffect, useMemo, useState } from "react";
import {
  FaAngleDoubleLeft,
  FaAngleDoubleRight,
  FaAngleLeft,
  FaAngleRight,
  FaExternalLinkAlt,
  FaSearch,
} from "react-icons/fa";
import {
  CellProps,
  Column,
  useGlobalFilter,
  usePagination,
  useTable,
} from "react-table";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Formal, FormalContext } from "../../model/Formal";
import { AdminTicket } from "../../model/Ticket";
import { TicketActions } from "./TicketActions";


export function FormalTicketsList() {
  const formal = useContext(FormalContext);
  const canWrite = useHasPermission("tickets", "write");
  const canDelete = useHasPermission("tickets", "delete");
  const [query, setQuery] = useState("");
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
        Cell: ({ value }) => {
          return value ? "Guest" : "King's";
        },
      },
    ];
    if (canWrite || canDelete) {
      base.push({
        Header: "Actions",
        Cell: ({ row }: CellProps<AdminTicket>) => {
          return (
            <TicketActions
              canDelete={canDelete}
              canWrite={canWrite}
              ticket={row.original}
            />
          );
        },
      });
    }
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
    setGlobalFilter,
    pageOptions,
    state: { pageIndex, pageSize },
    pageCount,
    setPageSize,
  } = useTable(
    {
      columns,
      data,
    },
    useGlobalFilter,
    usePagination
  );
  useEffect(() => setGlobalFilter(query), [query]);
  const background = useColorModeValue("white", "gray.750");
  const showOptions = useBreakpointValue({ base: false, lg: true });
  // TODO: Search filter
  return (
    <>
      <InputGroup size="sm" maxW="500px" mb={2}>
        <InputLeftAddon>
          <Icon as={FaSearch} />
        </InputLeftAddon>
        <Input
          id="query"
          autoComplete="off"
          // maxW="300px"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search tickets..."
        />
      </InputGroup>
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

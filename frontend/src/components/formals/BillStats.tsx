import {
  Flex,
  Heading,
  HStack,
  Icon,
  IconButton,
  Link,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  Select,
  Table,
  Text,
  Tbody,
  Td,
  Tfoot,
  Th,
  Thead,
  Tooltip,
  Tr,
  useColorModeValue,
  useBreakpointValue,
  Button,
} from "@chakra-ui/react";
import { useContext, useMemo } from "react";
import {
  FaAngleDoubleLeft,
  FaAngleDoubleRight,
  FaAngleLeft,
  FaAngleRight,
  FaDownload,
  FaExternalLinkAlt,
  FaTable,
} from "react-icons/fa";
import {
  CellProps,
  Column,
  useGlobalFilter,
  usePagination,
  useTable,
} from "react-table";
import { useBillStats } from "../../hooks/admin/useBillStats";
import { BillContext } from "../../model/Bill";
import { FormalCostBreakdown, UserCostBreakdown } from "../../model/BillStats";

interface BillFormalOverviewProps {
  stats: FormalCostBreakdown[];
}

function BillFormalOverview({ stats }: BillFormalOverviewProps) {
  const columns = useMemo<Column<FormalCostBreakdown>[]>(
    () => [
      {
        Header: "Formal",
        accessor: "formalName",
        Footer: "Total",
      },
      {
        Header: "King's",
        Cell: ({ row: { original } }: CellProps<FormalCostBreakdown>) => {
          return original.standard;
        },
        isNumeric: true,
        Footer: (table) => {
          return table.rows.reduce((acc, row) => {
            return acc + row.original.standard;
          }, 0);
        },
      },
      {
        Header: "Price",
        accessor: "price",
        isNumeric: true,
      },
      {
        Header: "Guests",
        Cell: ({ row: { original } }: CellProps<FormalCostBreakdown>) => {
          return original.guest;
        },
        isNumeric: true,
        Footer: (table) => {
          return table.rows.reduce((acc, row) => {
            return acc + row.original.guest;
          }, 0);
        },
      },
      {
        Header: "Guest Price",
        accessor: "guestPrice",
        isNumeric: true,
      },
      {
        Header: "Total",
        Cell: ({ row: { original } }: CellProps<FormalCostBreakdown>) => {
          return (
            original.standard * original.price +
            original.guest * original.guestPrice
          );
        },
        isNumeric: true,
        Footer: (table) => {
          return table.rows.reduce((acc, row) => {
            return (
              acc +
              row.original.standard * row.original.price +
              row.original.guest * row.original.guestPrice
            );
          }, 0);
        },
      },
    ],
    []
  );
  const {
    getTableBodyProps,
    getTableProps,
    prepareRow,
    headerGroups,
    footerGroups,
    rows,
  } = useTable({ columns, data: stats });
  const background = useColorModeValue("white", "gray.750");
  return (
    <Table my={2} variant="striped" size="sm" {...getTableProps()}>
      <Thead>
        {headerGroups.map((headerGroup) => (
          <Tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map((column) => (
              <Th
                p={1}
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

interface BillUserOverviewProps {
  stats: UserCostBreakdown[];
}

function BillUserOverview({ stats }: BillUserOverviewProps) {
  const columns = useMemo<Column<UserCostBreakdown>[]>(
    () => [
      {
        Header: "User",
        accessor: "userEmail",
        Cell: ({ value }) => {
          const crsid = value.split("@")[0];
          const email =
            value === "ents" ? import.meta.env.VITE_ENTS_EMAIL : value;
          return (
            <Link href={`mailto:${email}`} isExternal>
              {crsid} <Icon boxSize={3} as={FaExternalLinkAlt} />
            </Link>
          );
        },
        Footer: "Total",
      },
      {
        Header: "Total",
        accessor: "cost",
        isNumeric: true,
        Footer: (table) => {
          return table.rows.reduce((acc, row) => {
            return acc + row.original.cost;
          }, 0);
        },
      },
    ],
    []
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
    footerGroups,
    pageOptions,
    state: { pageIndex, pageSize },
    pageCount,
    setPageSize,
  } = useTable({ columns, data: stats }, useGlobalFilter, usePagination);
  // TODO: Filter
  const background = useColorModeValue("white", "gray.750");
  const showOptions = useBreakpointValue({ base: false, lg: true });
  return (
    <>
      <Table my={2} variant="striped" size="sm" {...getTableProps()}>
        <Thead>
          {headerGroups.map((headerGroup) => (
            <Tr {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map((column) => (
                <Th
                  p={1}
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
          {page.map((row) => {
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

export function BillStats() {
  const bill = useContext(BillContext);
  const base = import.meta.env.VITE_API_BASE_URL;
  const { data: stats, isLoading, isError } = useBillStats(bill.id);
  if (isError) {
    return null; // TODO
  }
  if (isLoading && !stats) {
    return null;
  }
  return (
    <>
      <Flex justifyContent="space-between" alignItems="center" mb={2}>
        <Heading as="h4" size="md">
          By Formal
        </Heading>
        <Button
          size="sm"
          as="a"
          href={`${base}admin/bills/${bill.id}/stats/formals.csv`}
          colorScheme="green"
          leftIcon={<Icon as={FaTable} />}
          rightIcon={<Icon as={FaDownload} />}
        >
          Download
        </Button>
      </Flex>
      <BillFormalOverview stats={stats?.formals ?? []} />
      <Flex justifyContent="space-between" alignItems="center" mb={2} mt={4}>
        <Heading as="h4" size="md">
          By User
        </Heading>
        <Button
          size="sm"
          as="a"
          href={`${base}admin/bills/${bill.id}/stats/users.csv`}
          colorScheme="green"
          leftIcon={<Icon as={FaTable} />}
          rightIcon={<Icon as={FaDownload} />}
        >
          Download
        </Button>
      </Flex>
      <BillUserOverview stats={stats?.users ?? []} />
    </>
  );
}

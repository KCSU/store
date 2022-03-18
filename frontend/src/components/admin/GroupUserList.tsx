import {
  Flex,
  Icon,
  Table,
  Tbody,
  Td,
  Tr,
  Text,
  IconButton,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  Select,
  Tooltip,
  HStack,
  useColorModeValue,
  useBreakpointValue,
  Box,
} from "@chakra-ui/react";
import {
  CellProps,
  Column,
  useFilters,
  usePagination,
  useTable,
} from "react-table";
import { useEffect, useMemo } from "react";
import {
  FaAngleDoubleLeft,
  FaAngleDoubleRight,
  FaAngleLeft,
  FaAngleRight,
  FaTrashAlt,
} from "react-icons/fa";
import { GroupUser } from "../../model/Group";

interface GroupUserListProps {
  users: GroupUser[];
  query: string;
  onDelete?: (user: GroupUser) => void | Promise<void>;
}

export function GroupUserList({
  users,
  query,
  onDelete,
}: GroupUserListProps) {
  // const filteredData = useMemo(() => {
  //   return users.filter(u => {
  //     u.userEmail.startsWith
  //   })
  // }, [users, query])
  const columns: Column<GroupUser>[] = useMemo(() => {
    if (onDelete) {
      return [
        {
          accessor: "userEmail",
        },
        {
          id: "actions",
          Cell: ({ row }: CellProps<GroupUser>) => (
            <Box textAlign="right">
              <IconButton
                aria-label="Remove user from group"
                variant="ghost"
                size="xs"
                icon={<Icon as={FaTrashAlt} />}
                onClick={() => onDelete(row.original)}
              />
            </Box>
          ),
        },
      ];
    } else {
      return [
        {
          accessor: "userEmail",
        },
      ];
    }
  }, []);
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    gotoPage,
    canPreviousPage,
    canNextPage,
    page,
    prepareRow,
    setFilter,
    previousPage,
    pageOptions,
    state: { pageIndex, pageSize },
    nextPage,
    pageCount,
    setPageSize,
  } = useTable(
    {
      columns,
      data: users,
      // filterTypes?
    },
    useFilters,
    usePagination
  );
  useEffect(() => setFilter("userEmail", query), [query]);
  const showOptions = useBreakpointValue({ base: false, lg: true });
  const background = useColorModeValue("white", "gray.750");
  return (
    <>
      <Table variant="striped" size="sm" {...getTableProps()} bg={background}>
        {/* <Thead></Thead> */}
        <Tbody {...getTableBodyProps()}>
          {page.map((row) => {
            prepareRow(row);
            return (
              <Tr {...row.getRowProps()}>
                {row.cells.map((cell) => (
                  <Td {...cell.getCellProps()}>{cell.render("Cell")}</Td>
                ))}
              </Tr>
            );
          })}
        </Tbody>
      </Table>
      <Flex justifyContent="space-between" alignItems="center">
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

import {
  IconButton,
  Icon,
  background,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  useColorModeValue,
  Alert,
  AlertIcon,
  useDisclosure,
  VStack,
  InputGroup,
  InputLeftAddon,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  Modal,
  ModalFooter,
  Button,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  Tfoot,
} from "@chakra-ui/react";
import { Field, FieldProps, Form, Formik } from "formik";
import { useContext, useMemo } from "react";
import { FaPlus, FaTrashAlt } from "react-icons/fa";
import { Column, CellProps, useTable } from "react-table";
import { useAddBillExtra } from "../../hooks/admin/useAddBillExtra";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { useRemoveBillExtra } from "../../hooks/admin/useRemoveBillExtra";
import { BillContext } from "../../model/Bill";
import { BillExtra } from "../../model/BillExtra";

function AddExtraButton() {
  const bill = useContext(BillContext);
  const mutation = useAddBillExtra(bill.id);
  const { isOpen, onOpen, onClose } = useDisclosure();
  return (
    <>
      <Button
        colorScheme="brand"
        leftIcon={<Icon as={FaPlus} />}
        mb={2}
        onClick={onOpen}
      >
        Add Extra Ents Charge
      </Button>
      <Formik
        initialValues={{
          description: "",
          amount: 0,
        }}
        onSubmit={async (values, form) => {
          await mutation.mutateAsync(values);
          form.resetForm();
          onClose();
        }}
      >
        {(props) => (
          <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
              <ModalHeader>Add Extra Ents Charge</ModalHeader>
              <ModalCloseButton />
              <ModalBody>
                <Form>
                  <VStack gap={2}>
                    <Field name="description">
                      {({ field, form }: FieldProps) => (
                        <FormControl
                          isInvalid={form.errors.description !== undefined}
                        >
                          <FormLabel htmlFor="description">
                            Description
                          </FormLabel>
                          <Input
                            {...field}
                            type="text"
                            placeholder="Description"
                          />
                          <FormErrorMessage>
                            {form.errors.description}
                          </FormErrorMessage>
                        </FormControl>
                      )}
                    </Field>
                    <Field name="amount">
                      {({ field, form }: FieldProps) => (
                        <FormControl
                          isInvalid={form.errors.amount !== undefined}
                        >
                          <FormLabel htmlFor="amount">Amount</FormLabel>
                          <InputGroup>
                            <InputLeftAddon>£</InputLeftAddon>
                            <NumberInput
                              width="100%"
                              {...field}
                              precision={2}
                              id="amount"
                              onChange={(_, val) =>
                                form.setFieldValue(field.name, val)
                              }
                            >
                              <NumberInputField borderLeftRadius={0} />
                              <NumberInputStepper>
                                <NumberIncrementStepper />
                                <NumberDecrementStepper />
                              </NumberInputStepper>
                            </NumberInput>
                          </InputGroup>
                          <FormErrorMessage>
                            {form.errors.amount}
                          </FormErrorMessage>
                        </FormControl>
                      )}
                    </Field>
                  </VStack>
                </Form>
              </ModalBody>
              <ModalFooter>
                <Button
                  colorScheme="brand"
                  mr={3}
                  onClick={props.submitForm}
                  isLoading={props.isSubmitting}
                >
                  Add
                </Button>
                <Button variant="ghost" onClick={onClose}>
                  Cancel
                </Button>
              </ModalFooter>
            </ModalContent>
          </Modal>
        )}
      </Formik>
    </>
  );
}

interface BillExtrasTableProps {
  billExtras: BillExtra[];
  billId?: string;
  showActions?: boolean;
}

export function BillExtrasTable({
  billExtras,
  billId,
  showActions = false,
}: BillExtrasTableProps) {
  const canWrite = useHasPermission("tickets", "write"); // FIXME: URGENT
  const columns = useMemo<Column<BillExtra>[]>(() => {
    const cols: Column<BillExtra>[] = [
      {
        accessor: "description",
        Header: "Description",
        Footer: "Total",
      },
      {
        accessor: "amount",
        Header: "Total",
        isNumeric: true,
        Footer: ({ rows }) => {
          const total = rows.reduce(
            (sum: number, row: any) => sum + row.values.amount,
            0
          );
          return total;
        },
      },
    ] as Column<BillExtra>[];
    if (canWrite && showActions) {
      cols.push({
        Header: "Actions",
        isNumeric: true,
        Cell({ row: { original } }: CellProps<BillExtra>) {
          const mutation = useRemoveBillExtra(billId ?? "");
          return (
            <IconButton
              aria-label="Revoke"
              size="xs"
              colorScheme="red"
              isLoading={mutation.isLoading}
              variant="ghost"
              onClick={() => {
                mutation.mutate(original.id);
              }}
            >
              <Icon as={FaTrashAlt} />
            </IconButton>
          );
        },
      } as Column<BillExtra>);
    }
    return cols;
  }, []);
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
    footerGroups,
  } = useTable({
    columns,
    data: billExtras,
  });
  const background = useColorModeValue("white", "gray.750");
  return (
    <Table variant="striped" size="sm" {...getTableProps()}>
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

export function BillExtrasList() {
  const bill = useContext(BillContext);
  const canWrite = useHasPermission("tickets", "write"); // FIXME: URGENT
  const extras = bill?.extras ?? [];
  return (
    <>
      {canWrite && <AddExtraButton />}
      {extras.length === 0 ? (
        <Alert status="info">
          <AlertIcon />
          No extra charges
        </Alert>
      ) : (
        <BillExtrasTable billExtras={extras} billId={bill.id} showActions />
      )}
    </>
  );
}

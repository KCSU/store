import { Container, Heading } from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";
import { BillDetailsForm } from "../../components/admin/BillDetailsForm";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { useCreateBill } from "../../hooks/admin/useCreateBill";
import { Bill } from "../../model/Bill";

export function AdminCreateBillView() {
  const mutation = useCreateBill();
  const navigate = useNavigate();
  const defaultBill: Bill = {
    id: '',
    name: '',
    start: new Date(),
    end: new Date(),
  };
  return (
    <Container maxW="container.md" p={0}>
      <BackButton to="/admin/bills">Back Home</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={4}>
          Create a Bill
        </Heading>
        <BillDetailsForm initialValues={defaultBill} onSubmit={async (values) => {
          await mutation.mutateAsync(values);
          navigate('/admin/bills');
        }}>
          Create Bill
        </BillDetailsForm>
      </Card>
    </Container>
  );
}
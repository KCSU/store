import {
  Container,
  Heading,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import { useMemo } from "react";
import { EditFormalForm } from "../../components/admin/EditFormalForm";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { Formal } from "../../model/Formal";

export function AdminCreateFormalView() {
  const defaultFormal: Formal = useMemo(() => {
    const currentDate = dayjs().startOf("day").toDate();
    return {
      id: 0,
      name: "",
      menu: "",
      price: 0,
      guestPrice: 0,
      guestLimit: 0,
      tickets: 0,
      ticketsRemaining: 0,
      guestTickets: 0,
      guestTicketsRemaining: 0,
      saleStart: currentDate,
      saleEnd: currentDate,
      dateTime: currentDate,
    };
  }, []);
  return (
    <Container maxW="container.md" p={0}>
      <BackButton to="/admin/formals">Back Home</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={4}>
          Create a Formal
        </Heading>
        {/* TODO: Groups */}
        <EditFormalForm formal={defaultFormal} />
      </Card>
    </Container>
  );
}

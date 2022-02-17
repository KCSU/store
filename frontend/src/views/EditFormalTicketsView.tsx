import { Container, Heading } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { EditTicketsForm } from "../components/tickets/EditTicketsForm";
import { BackButton } from "../components/utility/BackButton";
import { useTickets } from "../hooks/useTickets";

export function EditFormalTicketsView() {
    const {id} = useParams();
    const ticketId = parseInt(id ?? "0");
    const {data, isLoading, isError} = useTickets();
    if (isError || isLoading || !data) {
        return null;
    }
    const ticket = data.find(t => t.ticket.id === ticketId);
    if (!ticket) {
        return null;
    }
    return <Container maxW="container.md" p={0}>
        <BackButton to="/tickets">Back to Tickets</BackButton>
        <Heading as="h3" size="lg" mb={4}>
            Edit Tickets: {ticket.formal.name}
        </Heading>
        <EditTicketsForm ticket={ticket}></EditTicketsForm>
    </Container>
}
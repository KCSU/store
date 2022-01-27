import { useMutation, useQueryClient } from "react-query";
import { api } from "../config/api";

export function useCancelTicket() {
    const queryClient = useQueryClient();
    return useMutation((ticketId: number) => {
        return api.delete<void>(`tickets/${ticketId}`);
    }, {
        onSuccess() {
            // TODO: just delete? also other instances
            // Also TODO: brief delay??
            queryClient.invalidateQueries('tickets');
        }
    })
}
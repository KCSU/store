import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "./useCustomMutation";

export function useCancelTicket() {
    const queryClient = useQueryClient();
    return useCustomMutation((ticketId: string) => {
        return api.delete<void>(`tickets/${ticketId}`);
    }, {
        onSuccess() {
            // TODO: just delete? also other instances
            // Also TODO: brief delay??
            // Toast?
            queryClient.invalidateQueries('tickets');
        }
    })
}
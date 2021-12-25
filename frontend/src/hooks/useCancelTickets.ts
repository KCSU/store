import { useMutation, useQueryClient } from "react-query";
import { api } from "../config/api";

export function useCancelTickets() {
    const queryClient = useQueryClient();
    return useMutation((formalId: number) => {
        return api.delete<void>(`formals/${formalId}/tickets`);
    }, {
        onSuccess() {
            queryClient.invalidateQueries('tickets');
        }
    })
}
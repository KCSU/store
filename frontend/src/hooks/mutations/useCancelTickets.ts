import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useCancelTickets() {
    const queryClient = useQueryClient();
    return useCustomMutation((formalId: string) => {
        return api.delete<void>(`formals/${formalId}/tickets`);
    }, {
        onSuccess() {
            queryClient.invalidateQueries('tickets');
            queryClient.invalidateQueries('formals');
        }
    })
}
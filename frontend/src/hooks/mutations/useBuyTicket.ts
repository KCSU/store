import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { QueueRequest } from "../../model/QueueRequest";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useBuyTicket() {
  const queryClient = useQueryClient();
  return useCustomMutation(
    async (qr: QueueRequest) => {
      // TODO: return type?
      return api.post<void>("tickets", qr);
    },
    {
      onSuccess() {
        // FIXME: Alternatively, setQueryData?
        queryClient.invalidateQueries("tickets");
        queryClient.invalidateQueries("formals");
      }
    }
  );
}

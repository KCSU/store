import { useQuery } from "react-query";
import { api } from "../../config/api";
import { AccessLog } from "../../model/AccessLog";

type AccessLogDto = Omit<AccessLog, "createdAt"> & { createdAt: string };

export function useAccessLogs(
  page: number,
  size: number
) {
  return useQuery<AccessLog[]>(
    ["accessLogs", { page, size }],
    async () => {
      const response = await api.get<AccessLogDto[]>("admin/access", {
        params: {
          page, size
        }
      });
      return response.data.map((l) => {
        return {
          ...l,
          createdAt: new Date(l.createdAt),
        };
      });
    },
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
      keepPreviousData: true,
    }
  );
}

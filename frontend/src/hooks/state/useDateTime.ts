import dayjs from "dayjs";

export function useDateTime(date: Date): string {
    return dayjs(date).calendar();
}
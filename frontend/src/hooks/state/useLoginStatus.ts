import React from "react";
import { useLocalStorage } from "./useLocalStorage";

export function useLoginStatus(): readonly [Date | null, React.Dispatch<Date | null>] {
  const [d, setD] = useLocalStorage<Date | null>("_loginStatus", null);
  const date = d ? new Date(d) : null;
  return [date, setD];
}
import { useContext } from "react";
import { UserContext } from "../../model/User";

export function useHasPermission(resource: string, action: string): boolean {
  // const user = useContext(UserContext);
  return import.meta.env.DEV;
}

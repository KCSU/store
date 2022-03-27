import { useContext } from "react";
import { UserContext } from "../../model/User";

export function useHasPermission(resource: string, action: string): boolean {
  const user = useContext(UserContext);
  return user?.permissions.some((p) => {
    return (
      (p.resource === resource || p.resource === "*") &&
      (p.action === action || p.action === "*")
    );
  }) ?? false;
}

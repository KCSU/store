import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useHasPermission } from "../../hooks/admin/useHasPermission";

export interface RequireAdminProps {
  resource: string;
  action: string;
}

export const RequireAdmin: React.FC<RequireAdminProps> = ({
  resource,
  action,
  children,
}) => {
  const navigate = useNavigate();
  const hasAccess = useHasPermission(resource, action);
  useEffect(() => {
    // TODO: real admin permissions with useContext
    if (!hasAccess) {
      navigate("/", { replace: true });
    }
  });

  return <>{children}</>;
};

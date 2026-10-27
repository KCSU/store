import dayjs from "dayjs";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuthUser } from "../../hooks/queries/useAuthUser";
import { useLoginStatus } from "../../hooks/state/useLoginStatus";

export function RequireAuth() {
  const { data: user, isLoading } = useAuthUser();
  const location = useLocation();
  const [lastLogin] = useLoginStatus();

  if (lastLogin && lastLogin < dayjs().subtract(1, "day").toDate()) {
    return <Navigate to="/login" state={{from: location?.pathname}} />;
  }

  if (isLoading && !user) {
    return <></>;
  }

  // Redirect to login if not logged in, before rendering
  if (!user) {
    return <Navigate to="/login" replace state={{from: location?.pathname}} />;
  }

  return <Outlet />;
}

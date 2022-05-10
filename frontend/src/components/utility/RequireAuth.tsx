import dayjs from "dayjs";
import { useEffect } from "react";
import { Navigate, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useAuthUser } from "../../hooks/queries/useAuthUser";
import { useLoginStatus } from "../../hooks/state/useLoginStatus";

export function RequireAuth() {
  const { data: user, isLoading, isError } = useAuthUser();
  const location = useLocation();
  const [lastLogin] = useLoginStatus();
  const navigate = useNavigate();
  useEffect(() => {
    if (
      isError ||
      (!user && !isLoading)
    ) {
      navigate("/login", {
        replace: true,
        state: {
          from: location?.pathname,
        },
      });
    }
  });

  if (lastLogin && lastLogin < dayjs().subtract(1, "day").toDate()) {
    return <Navigate to="/login" state={{from: location?.pathname}} />;
  }

  if (isLoading && !user) {
    return <></>;
  }

  return <Outlet />;
}

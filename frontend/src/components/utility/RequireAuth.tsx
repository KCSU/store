import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuthUser } from "../../hooks/useAuthUser";

export function RequireAuth() {
    const {data: user, isLoading, isError} = useAuthUser();
    const location = useLocation();

    if (isError || (!user && !isLoading)) {
        return <Navigate to="/login" state={{from: location}} replace />
    }

    return <Outlet/>
}
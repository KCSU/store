import { useEffect } from "react";
import { Navigate, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useAuthUser } from "../../hooks/useAuthUser";

export function RequireAuth() {
    const {data: user, isLoading, isError} = useAuthUser();
    const location = useLocation();
    const navigate = useNavigate();
    useEffect(() => {
        if (isError || (!user && !isLoading)) {
            navigate('/login', {
                replace: true,
                state: {
                    from: location?.pathname
                }
            })
        }
    })

    if (isLoading && !user) {
        return <></>
    }

    return <Outlet/>
}
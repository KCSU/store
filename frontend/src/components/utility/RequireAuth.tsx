import { useEffect } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { useAuthUser } from "../../hooks/useAuthUser";
import { UserContext } from "../../model/User";

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

    return <UserContext.Provider value={user}>
        <Outlet/>
    </UserContext.Provider>
}
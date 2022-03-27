import { useAuthUser } from "../../hooks/queries/useAuthUser";
import { UserContext } from "../../model/User";

export const AuthProvider: React.FC = ({children}) => {
    const { data: user } = useAuthUser();
    return <UserContext.Provider value={user}>
        {children}
    </UserContext.Provider>
}
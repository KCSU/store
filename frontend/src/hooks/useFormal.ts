import { useQuery } from "react-query";
import { useFormals } from "./useFormals";

export function useFormal(id: number) {
    const {data: formals} = useFormals();
    return useQuery(['formals', id], () => {
        const formal = formals?.find(f => f.id === id);
        return formal;
    }, {
        enabled: formals === undefined
    })
    
}
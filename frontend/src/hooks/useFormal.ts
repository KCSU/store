import { useFormals } from "./useFormals";

export function useFormal(id: number) {
    // TODO: make query more efficient?
    const formals = useFormals();
    return formals.find(f => f.id === id);
}
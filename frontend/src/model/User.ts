import React from "react";
import { Group } from "./Group";

export interface User {
    id: number;
    name: string;
    email: string;
    groups: Group[]
    // Admin permissions
}

export const UserContext = React.createContext<User | undefined>(undefined);
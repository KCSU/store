import React from "react";
import { Group } from "./Group";
import { Permission } from "./Permission";

export interface User {
    id: number;
    name: string;
    email: string;
    groups: Group[]
    // Admin permissions
    permissions: Permission[];
}

export const UserContext = React.createContext<User | undefined>(undefined);
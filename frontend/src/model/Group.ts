export interface Group {
    id: string;
    name: string;
    type: string;
    lookup: string;
    users?: GroupUser[];
}

export function groupType(type: string) {
    switch (type) {
        case "inst":
            return "Institution";
        case "group":
            return "Group";
        case "manual":
            return "Manual";
        default:
            return type;
    }
        
}

export interface GroupUser {
    userEmail: string;
    groupId: string;
    isManual: boolean;
}
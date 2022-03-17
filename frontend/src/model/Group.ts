export interface Group {
    id: number;
    name: string;
    type: string;
    lookup: string;
    users?: GroupUser[];
}

export function groupType(group: Group) {
    switch (group.type) {
        case "inst":
            return "Institution";
        case "group":
            return "Group";
        case "manual":
            return "Manual";
        default:
            return group.type;
    }
        
}

export interface GroupUser {
    userEmail: string;
    groupId: number;
    isManual: boolean;
}
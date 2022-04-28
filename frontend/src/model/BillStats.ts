export interface FormalCostBreakdown {
  formalId: string;
  formalName: string;
  price: number;
  guestPrice: number;
  dateTime: string;
  standard: number;
  guest: number;
  standardManual: number;
  guestManual: number;
}

export interface UserCostBreakdown {
  userEmail: string;
  cost: number;
}

export interface BillStats {
  formals: FormalCostBreakdown[];
  users: UserCostBreakdown[];
}
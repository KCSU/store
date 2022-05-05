export interface AccessLog {
  id: string;
  email: string;
  message: string;
  createdAt: Date;
  metadata: Record<string, string>;
}

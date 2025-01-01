import { Gateway, TeamMember, User } from "./common";

export interface AcceptInviteProps {
  token: string;
  expired?: boolean;
  alreadyAccepted?: boolean;
  provider?: string;
}

export interface CompleteInviteProps {
  token: string;
  email?: string;
  username?: string;
  provider?: string;
  error?: string;
  status?: "valid" | "expired" | "accepted" | "invalid";
}

export interface SummaryStepProps {
  data: {
    firstName: string;
    lastName: string;
    email: string;
    username: string;
    dob: Date;
    phone: string;
    address: string;
  };
  onBack: () => void;
  onSubmit: () => void;
}

export interface RequestData{
  name: string;
  successful: number;
  errors: number;
}

export interface StackData {
  name: string;
  value: number;
}

export interface VisitorData {
  name: string;
  visitors: number;
}

export interface DashboardProps {
  teamName?: string;
  requestData?: RequestData[];
  techStack?: StackData[];
  user?: User;
  visitorData?: VisitorData[];
  health?: string;
  logoUrl?: string;
}

export interface TeamProps {
  user?: {
    name: string;
    email: string;
    avatar: string;
  };
  teamName?: string;
  logoUrl?: string;
  team?: TeamMember[];
  roles?: Array<{
    name: string;
    description: string;
    icon: string;
    color: string;
  }>;
}

export interface GatewayProps {
  gateways: Gateway[];
  application?: {
    id: string;
    name: string;
  };
  user: {
    name: string;
    email: string;
    username: string;
    provider: string;
  };
  teamName: string;
  logoUrl: string;
}

export interface GatewayTableProps {
  gateways: Gateway[];
  onEdit: (gateway: Gateway) => void;
  onDelete: (id: string) => void;
}
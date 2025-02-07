import { LucideIcon } from "lucide-react";
import { DateRange } from "react-day-picker";
import { ControllerRenderProps } from "react-hook-form";
import {
  Application,
  Gateway, GatewayConfig,
  Roles,
  RoleWithUsers,
  TeamMember,
  TechStack,
  TechStackWithApplications,
  User,
} from "./common";

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
  data: User;
  onBack: () => void;
  onSubmit: () => void;
}

export interface RequestData {
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
  stats?: any; // TODO: add type
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
  config?: GatewayConfig;
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

export interface Option {
  value: string;
  label: string;
  icon?: React.ReactElement;
}

export interface AutoCompleteProps {
  options: Option[];
  emptyMessage: string;
  value?: Option;
  onValueChange?: (value: Option) => void;
  isLoading?: boolean;
  disabled?: boolean;
  placeholder?: string;
  field?: ControllerRenderProps<any>;
}

export interface DatePickerProps {
  className?: string;
  date?: Date | DateRange | undefined;
  onDateChange?: (date: Date | DateRange | undefined) => void;
  isRangePicker?: boolean;
  minYear?: number;
  maxYear?: number;
  field?: ControllerRenderProps<any, any>;
}

export interface ApplicationsProps {
  user?: User;
  teamName: string;
  logoUrl: string;
  applications?: Application[] | null;
}

export interface SettingsProps {
  teamName: string;
  logoUrl: string;
  language: string;
  roles: RoleWithUsers[];
  techStacks: TechStackWithApplications[];
  user?: User;
}

export interface GeneralSettingsProps {
  teamName: string;
  logoUrl: string;
  language: string;
}

export interface RolesSettingsProps {
  roles: RoleWithUsers[];
}

export interface TechStacksSettingsProps {
  techStacks: TechStackWithApplications[];
}

export interface DialogButtonProps {
  buttonText: string;
  title?: string;
  description?: string;
  onOpen?: (boolean) => void;
  open?: boolean;
  icon?: LucideIcon;
  className?: string;
  variant?: "tooltip" | "text";
  children?: React.ReactNode;
}

export interface GatewayConfigProps {
  config?: GatewayConfig
}
export interface User {
  firstName?: string;
  lastName?: string;
  dob?: string;
  phone?: string;
  address?: string;
  email?: string;
  username?: string;
  avatar?: string;
  provider?: string;
}

export interface ApplicationStat {
  id: string;
  applicationId: string;
  environmentId: string;
  date: string;
  requests: number;
  errors: number;
  averageResponseTime: number;
  dataTransfered: number;
  uniqueVisitors: number;
  healthy: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface TechStack {
  id?: string;
  name: string;
  description: string;
}

export interface Application {
  id: string;
  appName: string;
  storageLocation: string;
  deployLocation: string;
  techStackId: string;
  status: "Building" | "Running" | "Stopped" | "Error";
  language?: string;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
  stats: ApplicationStat[];
  techStack: TechStack;
}

export interface Roles {
  name: string;
  description: string;
  icon: string;
  color: string;
}

export interface TeamMember {
  id: string;
  username: string;
  email: string;
  firstName: string;
  lastName: string;
  provider: string;
  roles: Roles[];
  techStacks: TechStack[];
}

export interface Gateway {
  id: string;
  name: string;
  path: string;
  httpMethod: string;
  backendUrl: string;
  requiresAuth: boolean;
  rateLimit: number;
  applicationId: string;
  application: Application;
}
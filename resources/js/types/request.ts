export interface CreateUserRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  username: string;
  dob: Date;
  address: string;
  phone: string;
  roles?: string[];
}

export interface CreateRoleRequest {
  name: string;
  description: string;
  icon: string;
  color: string;
}

export interface OnboardRequest {
  adminUser: CreateUserRequest;
  users: CreateUserRequest[];
  roles: CreateRoleRequest[];
  metadata: MetadataRequest;
}

export interface MetadataRequest {
  name: string;
  logoURL: string;
}

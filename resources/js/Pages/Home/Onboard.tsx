"use client";

import { Option } from "@/components/autocompletion";
import { ColorPicker } from "@/components/ColorPicker";
import { DatePicker } from "@/components/DatePicker";
import { InputAutoComplete } from "@/components/InputAutoComplete";
import { RenderFormItem } from "@/components/RenderFormItem";
import { RenderStepIndicators } from "@/components/RenderIndicators";
import { RoleIcon } from "@/components/RoleIcon";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { icons } from "@/lib/icons";
import {
  CreateRoleRequest,
  CreateUserRequest,
  MetadataRequest,
} from "@/lib/types";
import { zodResolver } from "@hookform/resolvers/zod";
import axios from "axios";
import {
  Check,
  GitBranch,
  Github,
  Trash2
} from "lucide-react";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { withMask } from "use-mask-input";
import * as z from "zod";

const adminSchema = z.object({
  firstName: z.string().min(2, "First name must be at least 2 characters"),
  lastName: z.string().min(2, "Last name must be at least 2 characters"),
  dob: z.date({ required_error: "Date of birth is required" }),
  address: z.string().min(5, "Address must be at least 5 characters"),
  phone: z.string().min(10, "Phone number must be at least 10 characters"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

const roleSchema = z.object({
  name: z.string().min(2, "Role name must be at least 2 characters"),
  description: z.string().min(5, "Description must be at least 5 characters"),
  icon: z.string().min(1, "Icon is required"),
  color: z.string().min(4, "Color must be a valid hex code"),
});

const serviceSchema = z.object({
  teamName: z.string().min(2, "Team name must be at least 2 characters"),
  logo: z.string().min(1, "Logo URL is required"),
  primaryColor: z.string().min(4, "Primary color must be a valid hex code"),
  secondaryColor: z.string().min(4, "Secondary color must be a valid hex code"),
});

export default function Onboarding() {
  const [step, setStep] = useState(1);
  const [adminData, setAdminData] = useState<CreateUserRequest | null>(null);
  const [roles, setRoles] = useState<CreateRoleRequest[]>([]);
  const [serviceData, setServiceData] = useState<MetadataRequest | null>(null);
  const totalSteps = 4;

  const iconNames: Option[] = icons.map((icon) => ({
    value: icon,
    label: icon,
  }));

  const adminForm = useForm<z.infer<typeof adminSchema>>({
    resolver: zodResolver(adminSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      dob: undefined,
      address: "",
      phone: "",
      password: "",
    },
  });

  const roleForm = useForm<z.infer<typeof roleSchema>>({
    resolver: zodResolver(roleSchema),
    defaultValues: {
      name: "",
      description: "",
      icon: "",
      color: "#000000",
    },
  });

  const serviceForm = useForm<z.infer<typeof serviceSchema>>({
    resolver: zodResolver(serviceSchema),
    defaultValues: {
      teamName: "",
      logo: "",
      primaryColor: "#000000",
      secondaryColor: "#ffffff",
    },
  });

  const onAdminSubmit = (data: z.infer<typeof adminSchema>) => {
    setAdminData(data);
    setStep(2);
    adminForm.reset();
  };

  const onRoleSubmit = (data: z.infer<typeof roleSchema>) => {
    setRoles([...roles, data]);
    roleForm.reset();
  };
  const onServiceSubmit = (data: z.infer<typeof serviceSchema>) => {
    setServiceData(data);

    const payload = {
      adminUser: adminData,
      roles: roles,
      metadata: serviceData,
    };

    const response = axios
      .post("/onboard", payload)
      .then((response) => {
        if (response.status === 200) {
          setStep(4);
        }
      })
      .catch((error) => {
        console.error(error);
      });
  };

  const handleAuthProvider = (provider: string) => {
    console.log(`Authenticating with ${provider}`);
    // In a real application, you would handle the authentication process here
  };

  const steps = [1, 2, 3, 4];

  const renderSidebar = () => (
    <div className="hidden w-1/4 bg-primary text-primary-foreground p-6 h-screen fixed left-0 top-0 overflow-y-auto lg:flex flex-col justify-center">
      <h2 className="text-2xl font-bold mb-4">Welcome to Neploy</h2>
      <p className="mb-4">
        We're excited to have you join us! This onboarding process will guide
        you through setting up your account and organization.
      </p>
      <div className="mb-6">
        <h3 className="text-xl font-semibold mb-2">Need Help?</h3>
        <p>Email: support@neploy.dev</p>
        <p>Phone: (123) 456-7890</p>
      </div>
      <div>
        <h3 className="text-xl font-semibold mb-2">Our Address</h3>
        <p>123 Service Street</p>
        <p>Tech City, TC 12345</p>
      </div>
    </div>
  );

  const renderStep = () => {
    switch (step) {
      case 1:
        return (
          <Card className="w-full max-w-screen-md mx-auto">
            <CardHeader>
              <CardTitle>Create Administrator User</CardTitle>
              <CardDescription>Set up the Super Dev account</CardDescription>
            </CardHeader>
            <Form {...adminForm}>
              <form onSubmit={adminForm.handleSubmit(onAdminSubmit)}>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <FormField
                      control={adminForm.control}
                      name="firstName"
                      render={({ field }) => (
                        <RenderFormItem label="First Name">
                          <Input {...field} />
                        </RenderFormItem>
                      )}
                    />
                    <FormField
                      control={adminForm.control}
                      name="lastName"
                      render={({ field }) => (
                        <RenderFormItem label="Last Name">
                          <Input {...field} />
                        </RenderFormItem>
                      )}
                    />
                  </div>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <FormField
                      control={adminForm.control}
                      name="dob"
                      render={({ field }) => (
                        <RenderFormItem
                          label="Date of Birth"
                          className="flex flex-col">
                          <Controller
                            control={adminForm.control}
                            name="dob"
                            render={({ field }) => (
                              <DatePicker
                                field={field}
                                maxYear={new Date().getFullYear() - 18}
                                {...field}
                                minYear={new Date().getFullYear() - 90}
                              />
                            )}
                          />
                        </RenderFormItem>
                      )}
                    />
                    <FormField
                      control={adminForm.control}
                      name="phone"
                      render={({ field }) => (
                        <FormItem className="flex flex-col">
                          <FormLabel>Phone</FormLabel>
                          <FormControl ref={withMask("(9999) 999-99-99")}>
                            <Input {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <FormField
                    control={adminForm.control}
                    name="address"
                    render={({ field }) => (
                      <RenderFormItem label="Address">
                        <Input {...field} />
                      </RenderFormItem>
                    )}
                  />
                  <FormField
                    control={adminForm.control}
                    name="password"
                    render={({ field }) => (
                      <RenderFormItem label="Password">
                        <Input type="password" {...field} />
                      </RenderFormItem>
                    )}
                  />
                </CardContent>
                <CardFooter className="flex flex-col items-center space-y-4">
                  <div className="flex space-x-2">
                    <Button
                      variant="outline"
                      type="button"
                      onClick={() => handleAuthProvider("GitHub")}>
                      <Github className="mr-2 h-4 w-4" />
                      GitHub
                    </Button>
                    <Button
                      variant="outline"
                      type="button"
                      onClick={() => handleAuthProvider("GitLab")}>
                      <GitBranch className="mr-2 h-4 w-4" />
                      GitLab
                    </Button>
                  </div>
                  <Button type="submit" className="w-full">
                    Next
                  </Button>
                </CardFooter>
              </form>
            </Form>
          </Card>
        );
      case 2:
        return (
          <Card className="w-full max-w-screen-md mx-auto">
            <CardHeader>
              <CardTitle>Create Roles</CardTitle>
              <CardDescription>
                Define roles for your organization
              </CardDescription>
            </CardHeader>
            <Form {...roleForm}>
              <form onSubmit={roleForm.handleSubmit(onRoleSubmit)}>
                <CardContent className="space-y-4">
                  <FormField
                    control={roleForm.control}
                    name="name"
                    render={({ field }) => (
                      <RenderFormItem label="Role Name">
                        <Input {...field} />
                      </RenderFormItem>
                    )}
                  />
                  <FormField
                    control={roleForm.control}
                    name="description"
                    render={({ field }) => (
                      <RenderFormItem label="Description">
                        <Textarea {...field} />
                      </RenderFormItem>
                    )}
                  />
                  <FormField
                    control={roleForm.control}
                    name="icon"
                    render={({ field }) => (
                      <RenderFormItem label="Icon">
                        <Controller
                          control={roleForm.control}
                          name="icon"
                          render={({ field }) => (
                            <InputAutoComplete
                              field={field}
                              OPTIONS={iconNames}
                            />
                          )}
                        />
                      </RenderFormItem>
                    )}
                  />
                  <FormField
                    control={roleForm.control}
                    name="color"
                    render={({ field }) => (
                      <RenderFormItem label="Color">
                        <Controller
                          control={roleForm.control}
                          name="color"
                          render={({ field }) => <ColorPicker field={field} />}
                        />
                      </RenderFormItem>
                    )}
                  />
                  {/* show the icon preview */}
                  <div className="flex items-start flex-col space-y-2">
                    <h2 className="font-semibold text-lg">Icon Preview</h2>
                    <RoleIcon
                      icon={roleForm.watch("icon")}
                      color={roleForm.watch("color")}
                    />
                  </div>

                  {roles.length > 0 && (
                    <div>
                      <h3 className="font-semibold text-xl">Roles Selected</h3>
                      <ul>
                        {roles.map((role, index) => (
                          <li
                            key={index}
                            className="flex justify-between items-center space-x-2 my-1">
                            <div className="flex items-center space-x-2 space-y-2">
                              <RoleIcon icon={role.icon} color={role.color} />
                              <div>
                                <p className="font-semibold">{role.name}</p>
                                <p>{role.description}</p>
                              </div>
                            </div>
                            <Button
                              type="button"
                              variant="destructive"
                              size="icon"
                              className="w-12 h-12"
                              onClick={() => {
                                setRoles(roles.filter((r, i) => i !== index));
                              }}>
                              <Trash2 className="!w-7 !h-7" />
                            </Button>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </CardContent>
                <CardFooter className="flex justify-between">
                  <Button type="submit">Add Role</Button>
                  <Button type="button" onClick={() => setStep(3)}>
                    Next
                  </Button>
                </CardFooter>
              </form>
            </Form>
          </Card>
        );
      case 3:
        return (
          <Card className="w-full max-w-screen-md mx-auto">
            <CardHeader>
              <CardTitle>Service Metadata</CardTitle>
              <CardDescription>Set up your team information</CardDescription>
            </CardHeader>
            <Form {...serviceForm}>
              <form onSubmit={serviceForm.handleSubmit(onServiceSubmit)}>
                <CardContent className="space-y-4">
                  <FormField
                    control={serviceForm.control}
                    name="teamName"
                    render={({ field }) => (
                      <RenderFormItem label="Team Name">
                        <Controller
                          control={serviceForm.control}
                          name="teamName"
                          render={({ field }) => <Input {...field} />}
                        />
                      </RenderFormItem>
                    )}
                  />
                  <FormField
                    control={serviceForm.control}
                    name="logo"
                    render={({ field }) => (
                      <RenderFormItem label="Logo URL">
                        <Controller
                          control={serviceForm.control}
                          name="logo"
                          render={({ field }) => <Input {...field} />}
                        />
                      </RenderFormItem>
                    )}
                  />
                  <FormField
                    control={serviceForm.control}
                    name="primaryColor"
                    render={({ field }) => (
                      <RenderFormItem label="Primary Color">
                        <Controller
                          control={serviceForm.control}
                          name="primaryColor"
                          render={({ field }) => <ColorPicker field={field} />}
                        />
                      </RenderFormItem>
                    )}
                  />
                  <FormField
                    control={serviceForm.control}
                    name="secondaryColor"
                    render={({ field }) => (
                      <RenderFormItem label="Secondary Color">
                        <Controller
                          control={serviceForm.control}
                          name="secondaryColor"
                          render={({ field }) => <ColorPicker field={field} />}
                        />
                      </RenderFormItem>
                    )}
                  />
                </CardContent>
                <CardFooter>
                  <Button type="submit">Finish</Button>
                </CardFooter>
              </form>
            </Form>
          </Card>
        );
      case 4:
        return (
          <Card>
            <CardHeader>
              <CardTitle>All Done!</CardTitle>
              <CardDescription>
                Your onboarding process is complete
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col items-center justify-center space-y-4">
                <Check className="w-16 h-16 text-green-500" />
                <p className="text-lg text-center">
                  Congratulations! You've successfully completed the onboarding
                  process. Your account and organization are now set up.
                </p>
              </div>
            </CardContent>
            <CardFooter>
              <Button
                onClick={() =>
                  window.location.replace("/")
                }>
                Go to Login
              </Button>
            </CardFooter>
          </Card>
        );
      default:
        return null;
    }
  };

  return (
    <div className="flex min-h-[900px] h-screen w-full">
      {renderSidebar()}
      <div className="flex-1 lg:ml-[25%] p-3 lg:p-10">
        <h1 className="text-3xl font-bold mb-6">Onboarding</h1>
        <RenderStepIndicators
          step={step}
          totalSteps={totalSteps}
          steps={steps}
        />
        {renderStep()}
      </div>
    </div>
  );
}

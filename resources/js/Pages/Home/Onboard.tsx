"use client";

import { useState } from "react";
import * as React from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { format } from "date-fns";
import {
  Calendar as CalendarIcon,
  Github,
  GitBranch,
  Check,
  ChevronsUpDown,
} from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  CardFooter,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { ColorPicker } from "@/components/ColorPicker";
import { AutoComplete, Option } from "@/components/autocompletion";
import { InputAutoComplete } from "@/components/InputAutoComplete";
import { icons } from "@/lib/arrays";
import { withMask } from "use-mask-input";

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

const userSchema = z.object({
  firstName: z.string().min(2, "First name must be at least 2 characters"),
  lastName: z.string().min(2, "Last name must be at least 2 characters"),
  email: z.string().email("Invalid email address"),
  role: z.string().min(1, "Role is required"),
});

const serviceSchema = z.object({
  teamName: z.string().min(2, "Team name must be at least 2 characters"),
  logo: z.string().min(1, "Logo URL is required"),
  primaryColor: z.string().min(4, "Primary color must be a valid hex code"),
  secondaryColor: z.string().min(4, "Secondary color must be a valid hex code"),
});

export default function Onboarding() {
  const [step, setStep] = useState(1);
  const [adminData, setAdminData] = useState(null);
  const [roles, setRoles] = useState([]);
  const [users, setUsers] = useState([]);
  const [serviceData, setServiceData] = useState(null);
  const totalSteps = 5;

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
      icon: "ChevronLeft",
      color: "#000000",
    },
  });

  const userForm = useForm<z.infer<typeof userSchema>>({
    resolver: zodResolver(userSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      email: "",
      role: "",
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

  const onUserSubmit = (data: z.infer<typeof userSchema>) => {
    setUsers([...users, data]);
    userForm.reset();
  };

  const onServiceSubmit = (data: z.infer<typeof serviceSchema>) => {
    setServiceData(data);
    setStep(5);
  };

  const handleAuthProvider = (provider: string) => {
    console.log(`Authenticating with ${provider}`);
    // In a real application, you would handle the authentication process here
  };

  const steps = [1, 2, 3, 4, 5];

  const renderStepIndicators = () => (
    <div className="flex justify-center mb-8">
      {steps.map((_, index) => (
        <div key={index} className="flex items-center">
          <div
            className={`w-8 h-8 rounded-full flex items-center justify-center ${
              step > index + 1
                ? "bg-lime-700 text-white"
                : step === index + 1
                  ? "bg-primary text-white"
                  : "bg-gray-200"
            }`}
          >
            {step > index + 1 ? <Check className="w-4 h-4" /> : index + 1}
          </div>
          {index < totalSteps - 1 && (
            <div
              className={`w-6 sm:w-8 md:w-10 lg:w-12 h-1 ${
                step - 1 === index + 1
                  ? "bg-gradient-to-r from-lime-700 to-primary from-40% to-90%"
                  : step > index + 1
                    ? "bg-lime-700"
                    : "bg-gray-200"
              }`}
            />
          )}
        </div>
      ))}
    </div>
  );

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
                        <FormItem>
                          <FormLabel>First Name</FormLabel>
                          <FormControl>
                            <Input {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={adminForm.control}
                      name="lastName"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Last Name</FormLabel>
                          <FormControl>
                            <Input {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <FormField
                      control={adminForm.control}
                      name="dob"
                      render={({ field }) => (
                        <FormItem className="flex flex-col">
                          <FormLabel>Date of birth</FormLabel>
                          <Popover>
                            <PopoverTrigger asChild>
                              <FormControl>
                                <Button
                                  variant={"outline"}
                                  className={cn(
                                    "w-full pl-3 text-left font-normal",
                                    !field.value && "text-muted-foreground",
                                  )}
                                >
                                  {field.value ? (
                                    format(field.value, "PPP")
                                  ) : (
                                    <span>Pick a date</span>
                                  )}
                                  <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                </Button>
                              </FormControl>
                            </PopoverTrigger>
                            <PopoverContent
                              className="w-auto p-0"
                              align="start"
                            >
                              <Calendar
                                mode="single"
                                selected={field.value}
                                onSelect={field.onChange}
                                disabled={(date) =>
                                  date > new Date() ||
                                  date < new Date("1900-01-01")
                                }
                                initialFocus
                              />
                            </PopoverContent>
                          </Popover>
                          <FormMessage />
                        </FormItem>
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
                      <FormItem>
                        <FormLabel>Address</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={adminForm.control}
                    name="password"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Password</FormLabel>
                        <FormControl>
                          <Input type="password" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </CardContent>
                <CardFooter className="flex flex-col items-center space-y-4">
                  <div className="flex space-x-2">
                    <Button
                      variant="outline"
                      type="button"
                      onClick={() => handleAuthProvider("GitHub")}
                    >
                      <Github className="mr-2 h-4 w-4" />
                      GitHub
                    </Button>
                    <Button
                      variant="outline"
                      type="button"
                      onClick={() => handleAuthProvider("GitLab")}
                    >
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
                  <div className="grid grid-cols-2 md:grid-cols-3 gap-4 w-11/12">
                    <FormField
                      control={roleForm.control}
                      name="name"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Role Name</FormLabel>
                          <FormControl>
                            <Input
                              className="col-span-3 md:col-span-1"
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={roleForm.control}
                      name="icon"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Icon</FormLabel>
                          <FormControl>
                            {/* when the field value changes, always it would filter the iconNames and get the first 50 items */}
                            <InputAutoComplete
                              field={field}
                              OPTIONS={iconNames}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={roleForm.control}
                      name="color"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Color</FormLabel>
                          <FormControl>
                            <ColorPicker {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <FormField
                    control={roleForm.control}
                    name="description"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Description</FormLabel>
                        <FormControl>
                          <Textarea {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
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
              <CardTitle>Create Users</CardTitle>
              <CardDescription>Add users to your organization</CardDescription>
            </CardHeader>
            <Form {...userForm}>
              <form onSubmit={userForm.handleSubmit(onUserSubmit)}>
                <CardContent className="space-y-4">
                  <FormField
                    control={userForm.control}
                    name="firstName"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>First Name</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={userForm.control}
                    name="lastName"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Last Name</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={userForm.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Email</FormLabel>
                        <FormControl>
                          <Input type="email" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={userForm.control}
                    name="role"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Role</FormLabel>
                        <Select
                          onValueChange={field.onChange}
                          defaultValue={field.value}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Select a role" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {roles.map((role, index) => (
                              <SelectItem key={index} value={role.name}>
                                {role.name}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </CardContent>
                <CardFooter className="flex justify-between">
                  <Button type="submit">Add User</Button>
                  <Button type="button" onClick={() => setStep(4)}>
                    Next
                  </Button>
                </CardFooter>
              </form>
            </Form>
          </Card>
        );
      case 4:
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
                      <FormItem>
                        <FormLabel>Team Name</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={serviceForm.control}
                    name="logo"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Logo URL</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={serviceForm.control}
                    name="primaryColor"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Primary Color</FormLabel>
                        <FormControl>
                          <ColorPicker {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={serviceForm.control}
                    name="secondaryColor"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Secondary Color</FormLabel>
                        <FormControl>
                          <ColorPicker {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
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
      case 5:
        return (
          <Card className="w-full max-w-screen-md mx-auto">
            <CardHeader>
              <CardTitle>Overview</CardTitle>
              <CardDescription>Review all registered data</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {adminData && (
                <div>
                  <h3 className="font-semibold">Administrator</h3>
                  <p>
                    Name: {adminData.firstName} {adminData.lastName}
                  </p>
                  <p>
                    Date of Birth:{" "}
                    {adminData.dob
                      ? format(adminData.dob, "PPP")
                      : "Not provided"}
                  </p>
                  <p>Address: {adminData.address}</p>
                  <p>Phone: {adminData.phone}</p>
                </div>
              )}
              <div>
                <h3 className="font-semibold">Roles</h3>
                <ul>
                  {roles.map((role, index) => (
                    <li key={index}>
                      {role.name} - {role.description}
                    </li>
                  ))}
                </ul>
              </div>
              <div>
                <h3 className="font-semibold">Users</h3>
                <ul>
                  {users.map((user, index) => (
                    <li key={index}>
                      {user.firstName} {user.lastName} - {user.email} (
                      {user.role})
                    </li>
                  ))}
                </ul>
              </div>
              {serviceData && (
                <div>
                  <h3 className="font-semibold">Service Metadata</h3>
                  <p>Team Name: {serviceData.teamName}</p>
                  <p>Logo URL: {serviceData.logo}</p>
                  <p>Primary Color: {serviceData.primaryColor}</p>
                  <p>Secondary Color: {serviceData.secondaryColor}</p>
                </div>
              )}
            </CardContent>
            <CardFooter>
              <Button onClick={() => setStep(6)}>Complete</Button>
            </CardFooter>
          </Card>
        );
      case 6:
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
                  console.log("Redirect to dashboard or home page")
                }
              >
                Go to Dashboard
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
        {renderStepIndicators()}
        {renderStep()}
      </div>
    </div>
  );
}

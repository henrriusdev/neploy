"use client";

import { useState } from "react";
import { format } from "date-fns";
import { Calendar as CalendarIcon } from "lucide-react";
import { Separator } from "@/components/ui/separator";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  CardFooter,
} from "@/components/ui/card";

export default function Onboarding() {
  const [step, setStep] = useState(1);
  const [adminData, setAdminData] = useState(null);

  const onAdminSubmit = (event) => {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = Object.fromEntries(formData);
    setAdminData(data);
    setStep(2);
  };

  const handleAuthProvider = (provider) => {
    console.log(`Authenticating with ${provider}`);
    console.log(`Received email and username from ${provider}`);
  };

  const steps = [
    { title: "Administrator", description: "Create the Super Dev account" },
    { title: "Roles", description: "Define organization roles" },
    { title: "Users", description: "Add users to your organization" },
    { title: "Service", description: "Set up team information" },
    { title: "Overview", description: "Review all data" },
  ];

  const renderStepIndicators = () => (
    <div className="flex justify-center space-x-2 mb-6">
      {steps.map((_, index) => (
        <div
          key={index}
          className={`w-8 h-2 rounded-full ${
            step > index ? "bg-primary" : "bg-muted"
          }`}
        />
      ))}
    </div>
  );

  return (
    <div className="flex min-h-screen">
      {/* Side Section */}
      <div className="w-1/4 bg-black text-white p-8">
        <h2 className="text-2xl font-bold mb-4">Company Name</h2>
        <p className="mb-4">
          Welcome to our onboarding process. We're excited to have you join us!
        </p>
        <div className="mb-6">
          <h3 className="text-xl font-semibold mb-2">Contact Information</h3>
          <p>Email: support@company.com</p>
          <p>Phone: (123) 456-7890</p>
        </div>
        <div>
          <h3 className="text-xl font-semibold mb-2">Office Address</h3>
          <p>123 Tech Street</p>
          <p>Silicon Valley, CA 94000</p>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 p-10 bg-gray-100">
        <h1 className="text-3xl font-bold mb-6">Onboarding</h1>
        {renderStepIndicators()}
        <Card>
          <CardHeader>
            <CardTitle>Create Administrator User</CardTitle>
            <CardDescription>Set up the Super Dev account</CardDescription>
          </CardHeader>
          <form onSubmit={onAdminSubmit}>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="firstName">First Name</Label>
                  <Input id="firstName" name="firstName" />
                </div>
                <div>
                  <Label htmlFor="lastName">Last Name</Label>
                  <Input id="lastName" name="lastName" />
                </div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="dob">Date of birth</Label>
                  <Input id="dob" name="dob" type="date" />
                </div>
                <div>
                  <Label htmlFor="phone">Phone</Label>
                  <Input id="phone" name="phone" />
                </div>
              </div>
              <div>
                <Label htmlFor="address">Address</Label>
                <Input id="address" name="address" />
              </div>
              <div>
                <Label htmlFor="password">Password</Label>
                <Input id="password" name="password" type="password" />
              </div>
            </CardContent>
            <CardFooter className="flex flex-col items-center space-y-4">
              <div className="w-full flex flex-col items-center space-y-2">
                <Separator className="my-4" />
                <p className="text-sm text-muted-foreground">
                  Or sign up with:
                </p>
                <div className="flex space-x-4">
                  <Button
                    variant="outline"
                    type="button"
                    onClick={() => handleAuthProvider("GitHub")}>
                    GitHub
                  </Button>
                  <Button
                    variant="outline"
                    type="button"
                    onClick={() => handleAuthProvider("GitLab")}>
                    GitLab
                  </Button>
                </div>
              </div>
              <Button type="submit" className="w-full">
                Next
              </Button>
            </CardFooter>
          </form>
        </Card>
      </div>
    </div>
  );
}

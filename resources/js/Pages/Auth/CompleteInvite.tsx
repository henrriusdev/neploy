import * as React from "react";
import { router } from "@inertiajs/react";
import { useToast } from "@/hooks/use-toast";
import ProviderStep from "./InviteSteps/ProviderStep";
import UserDataStep from "./InviteSteps/UserDataStep";
import SummaryStep from "./InviteSteps/SummaryStep";
import axios from "axios";
import { CompleteInviteProps } from "@/types/props";
import { User } from "@/types/common";

type Step = "provider" | "data" | "summary";

export default function CompleteInvite({
  token,
  email,
  username,
  provider,
  error,
  status,
}: CompleteInviteProps) {
  const [step, setStep] = React.useState<Step>(() => {
    // If we have provider info, start at data step
    return provider ? "data" : "provider";
  });
  const [userData, setUserData] = React.useState<User>({
    firstName: "",
    lastName: "",
    dob: "",
    phone: "",
    address: "",
    email: email || "",
    username: username || "",
    password: "",
  });
  const { toast } = useToast();

  React.useEffect(() => {
    if (status === "expired") {
      toast({
        title: "Invitation Expired",
        description: "This invitation has expired. Please request a new one.",
        variant: "destructive",
      });
    } else if (status === "accepted") {
      toast({
        title: "Already Accepted",
        description: "This invitation has already been accepted.",
        variant: "destructive",
      });
    } else if (status === "invalid") {
      toast({
        title: "Invalid Invitation",
        description: error || "This invitation is invalid or has been revoked.",
        variant: "destructive",
      });
    }
  }, [status, error]);

  // Don't show the form if the invitation is not valid
  if (status !== "valid") {
    router.visit("/login");
    return null;
  }

  const handleProviderNext = () => {
    setStep("data");
  };

  const handleDataNext = (data: User) => {
    setUserData(data);
    setStep("summary");
  };

  const handleDataBack = () => {
    setStep("provider");
  };

  const handleSummaryBack = () => {
    setStep("data");
  };

  const handleSubmit = () => {
    axios
      .post("/users/complete-invite", {
        token,
        ...userData,
      })
      .then((response) => {
        if (response.status === 200) {
          toast({
            title: "Success",
            description: "Your account has been created successfully!",
          });
          window.location.replace("/");
        }
      })
      .catch((error) => {
        const errorMessage =
          error.response?.data?.error ||
          error.message ||
          "Failed to complete registration";
        toast({
          title: "Error",
          description: errorMessage,
          variant: "destructive",
        });
      });
  };

  const renderStep = () => {
    switch (step) {
      case "provider":
        return <ProviderStep onNext={handleProviderNext} token={token} />;
      case "data":
        return (
          <UserDataStep
            email={email}
            username={username}
            onNext={handleDataNext}
            onBack={handleDataBack}
          />
        );
      case "summary":
        return (
          <SummaryStep
            data={userData}
            onBack={handleSummaryBack}
            onSubmit={handleSubmit}
          />
        );
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary to-accent from-40% to-95% py-12 px-4 sm:px-6 lg:px-8">
      <div className="w-full">{renderStep()}</div>
    </div>
  );
}

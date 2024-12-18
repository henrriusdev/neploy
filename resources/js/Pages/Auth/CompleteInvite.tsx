import * as React from "react";
import { router } from "@inertiajs/react";
import { useToast } from "@/hooks/use-toast";
import ProviderStep from "./InviteSteps/ProviderStep";
import UserDataStep from "./InviteSteps/UserDataStep";
import SummaryStep from "./InviteSteps/SummaryStep";
import axios from "axios";

interface Props {
  token: string;
  email?: string;
  username?: string;
  provider?: string;
  error?: string;
  status?: "valid" | "expired" | "accepted" | "invalid";
}

interface UserData {
  firstName: string;
  lastName: string;
  dob: string;
  phone: string;
  address: string;
  email: string;
  username: string;
  password: string;
}

type Step = "provider" | "data" | "summary";

export default function CompleteInvite({
  token,
  email,
  username,
  provider,
  error,
  status,
}: Props) {
  const [step, setStep] = React.useState<Step>(() => {
    // If we have provider info, start at data step
    return provider ? "data" : "provider";
  });
  const [userData, setUserData] = React.useState<UserData>({
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

  const handleDataNext = (data: UserData) => {
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

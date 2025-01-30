import * as React from "react";
import { router } from "@inertiajs/react";
import { useToast } from "@/hooks/use-toast";
import { ProviderStep, UserDataStep, SummaryStep } from "@/components/steps";
import { CompleteInviteProps } from "@/types/props";
import { User } from "@/types/common";
import { useCompleteInviteMutation } from "@/services/api/auth";

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
  const [completeInvite, { isLoading }] = useCompleteInviteMutation();

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

  if (status !== "valid") {
    router.visit("/login");
    return null;
  }

  const handleProviderNext = () => {
    setStep("data");
  };

  const handleDataNext = (data: User) => {
    const formattedData = {
      firstName: data.firstName || "",
      lastName: data.lastName || "",
      dob: data.dob || "",
      phone: data.phone || "",
      address: data.address || "",
      email: data.email || "",
      username: data.username || "",
    };
    setUserData(formattedData);
    setStep("summary");
  };

  const handleDataBack = () => {
    setStep("provider");
  };

  const handleSummaryBack = () => {
    setStep("data");
  };

  const handleSubmit = () => {
    if (
      !userData.firstName ||
      !userData.lastName ||
      !userData.phone ||
      !userData.address ||
      !userData.email ||
      !userData.username
    ) {
      toast({
        title: "Error",
        description: "Please fill in all required fields",
      });
      return;
    }

    const submitData = {
      token,
      firstName: userData.firstName,
      lastName: userData.lastName,
      phone: userData.phone,
      address: userData.address,
      email: userData.email,
      username: userData.username,
    };

    completeInvite(submitData)
      .unwrap()
      .then(() => {
        toast({
          title: "Success",
          description: "Your account has been created successfully!",
        });
        window.location.replace("/");
      })
      .catch((error) => {
        const errorMessage =
          error.data?.error || "Failed to complete registration";
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

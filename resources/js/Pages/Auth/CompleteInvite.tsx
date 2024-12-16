import * as React from 'react'
import { router } from '@inertiajs/react'
import { useToast } from '@/hooks/use-toast'
import ProviderStep from './InviteSteps/ProviderStep'
import UserDataStep from './InviteSteps/UserDataStep'
import SummaryStep from './InviteSteps/SummaryStep'

interface Props {
    token: string;
    email: string;
    username?: string;
}

type Step = 'provider' | 'data' | 'summary';

export default function CompleteInvite({ token, email, username }: Props) {
    const [step, setStep] = React.useState<Step>('provider')
    const [userData, setUserData] = React.useState<any>(null)
    const { toast } = useToast()

    const handleProviderNext = () => {
        setStep('data')
    }

    const handleDataNext = (data: any) => {
        setUserData(data)
        setStep('summary')
    }

    const handleDataBack = () => {
        setStep('provider')
    }

    const handleSummaryBack = () => {
        setStep('data')
    }

    const handleSubmit = () => {
        router.post('/users/complete-invite', 
            { 
                token,
                ...userData
            },
            {
                onSuccess: () => {
                    toast({
                        title: "Success",
                        description: "Your account has been created successfully!",
                    })
                    router.visit('/dashboard')
                },
                onError: (errors) => {
                    toast({
                        title: "Error",
                        description: errors.message || "Failed to complete registration",
                        variant: "destructive",
                    })
                }
            }
        )
    }

    const renderStep = () => {
        switch (step) {
            case 'provider':
                return <ProviderStep onNext={handleProviderNext} />
            case 'data':
                return (
                    <UserDataStep 
                        email={email}
                        username={username}
                        onNext={handleDataNext}
                        onBack={handleDataBack}
                    />
                )
            case 'summary':
                return (
                    <SummaryStep
                        data={userData}
                        onBack={handleSummaryBack}
                        onSubmit={handleSubmit}
                    />
                )
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
            <div className="w-full">
                {renderStep()}
            </div>
        </div>
    )
}

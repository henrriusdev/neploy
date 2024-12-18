import * as React from 'react'
import { useState, useEffect } from 'react'
import axios from 'axios'
import { router } from '@inertiajs/react'
import { useToast } from '@/hooks/use-toast'
import ProviderStep from '../Auth/InviteSteps/ProviderStep'
import UserDataStep from '../Auth/InviteSteps/UserDataStep'
import RolesStep from './Steps/RolesStep'
import ServiceStep from './Steps/ServiceStep'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Check } from "lucide-react"
import { OnboardingSidebar } from '@/components/OnboardingSidebar'

interface Props {
    email?: string;
    username?: string;
}

type Step = 'provider' | 'data' | 'roles' | 'service' | 'summary';

export default function Onboard({ email, username }: Props) {
    const [step, setStep] = useState<Step>('provider')
    const [adminData, setAdminData] = useState<any>(null)
    const [roles, setRoles] = useState<any[]>([])
    const [serviceData, setServiceData] = useState<any>(null)
    const { toast } = useToast()

    useEffect(() => {
        const params = new URLSearchParams(window.location.search)
        const provider = params.get("provider")
        const username = params.get("username")
        const email = params.get("email")

        if (provider && username && email) {
            setAdminData({ provider, username, email })
            setStep('data')
        }
    }, [])

    const handleProviderNext = () => {
        setStep('data')
    }

    const handleDataNext = (data: any) => {
        setAdminData(data)
        setStep('roles')
    }

    const handleRolesNext = () => {
        setStep('service')
    }

    const handleServiceNext = (data: any) => {
        setServiceData(data)
        setStep('summary')
    }

    const handleDataBack = () => {
        setStep('provider')
    }

    const handleRolesBack = () => {
        setStep('data')
    }

    const handleServiceBack = () => {
        setStep('roles')
    }

    const handleSummaryBack = () => {
        setStep('service')
    }

    const handleSubmit = () => {
        const payload = {
            adminUser: adminData,
            roles: roles,
            metadata: serviceData,
        }

        axios.post('/onboard', payload)
            .then(response => {
                if (response.status === 200) {
                    toast({
                        title: "Success",
                        description: "Your account has been set up successfully!",
                    })
                    router.visit('/dashboard')
                }
            })
            .catch(error => {
                toast({
                    title: "Error",
                    description: error.response?.data?.message || "Failed to complete setup",
                    variant: "destructive",
                })
            })
    }

    const renderStep = () => {
        switch (step) {
            case 'provider':
                return <ProviderStep onNext={handleProviderNext} />
            case 'data':
                return (
                    <UserDataStep 
                        email={adminData?.email || email}
                        username={adminData?.username || username}
                        onNext={handleDataNext}
                        onBack={handleDataBack}
                    />
                )
            case 'roles':
                return (
                    <RolesStep
                        roles={roles}
                        setRoles={setRoles}
                        onNext={handleRolesNext}
                        onBack={handleRolesBack}
                    />
                )
            case 'service':
                return (
                    <ServiceStep
                        onNext={handleServiceNext}
                        onBack={handleServiceBack}
                    />
                )
            case 'summary':
                return (
                    <Card className="w-full max-w-screen-md mx-auto">
                        <CardHeader>
                            <CardTitle>Review Your Setup</CardTitle>
                            <CardDescription>Please verify your information before completing</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <div>
                                <h3 className="font-medium text-lg">Administrator Account</h3>
                                <dl className="mt-2 space-y-1">
                                    <div>
                                        <dt className="text-sm text-muted-foreground">Name</dt>
                                        <dd>{adminData.firstName} {adminData.lastName}</dd>
                                    </div>
                                    <div>
                                        <dt className="text-sm text-muted-foreground">Email</dt>
                                        <dd>{adminData.email}</dd>
                                    </div>
                                </dl>
                            </div>
                            <div>
                                <h3 className="font-medium text-lg">Roles ({roles.length})</h3>
                                <ul className="mt-2 space-y-1">
                                    {roles.map((role, index) => (
                                        <li key={index}>{role.name}</li>
                                    ))}
                                </ul>
                            </div>
                            <div>
                                <h3 className="font-medium text-lg">Team Information</h3>
                                <dl className="mt-2 space-y-1">
                                    <div>
                                        <dt className="text-sm text-muted-foreground">Team Name</dt>
                                        <dd>{serviceData.teamName}</dd>
                                    </div>
                                </dl>
                            </div>
                        </CardContent>
                        <CardFooter className="flex justify-between">
                            <Button type="button" variant="outline" onClick={handleSummaryBack}>
                                Back
                            </Button>
                            <Button onClick={handleSubmit}>
                                Complete Setup
                            </Button>
                        </CardFooter>
                    </Card>
                )
        }
    }

    return (
        <div className="flex min-h-screen">
            <OnboardingSidebar currentStep={step} />
            <div className="flex-1 p-6">
                <h1 className="text-3xl font-bold mb-6 text-center">Setup Your Account</h1>
                {renderStep()}
            </div>
        </div>
    )
}

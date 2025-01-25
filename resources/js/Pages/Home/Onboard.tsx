import { OnboardingSidebar } from '@/components/OnboardingSidebar'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { useToast } from '@/hooks/use-toast'
import { useOnboardMutation } from '@/services/api/onboard'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next';
import '@/i18n';
import ProviderStep from '../Auth/InviteSteps/ProviderStep'
import UserDataStep from '../Auth/InviteSteps/UserDataStep'
import RolesStep from './Steps/RolesStep'
import ServiceStep from './Steps/ServiceStep'

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
    const [onboard] = useOnboardMutation();
    const { t } = useTranslation();

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
        setServiceData({
            ...data,
            language: data.language
        });
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

    const handleSubmit = async () => {
        const payload = {
            adminUser: adminData,
            roles: roles,
            metadata: {
                ...serviceData,
                language: serviceData.language
            },
        }

        try {
            const response = await onboard({ data: payload });
            
            if ('data' in response && response.data.success) {
                toast({
                    title: t('common.success'),
                    description: t('onboarding.success'),
                });
                window.location.replace('/');
            } else {
                toast({
                    title: t('common.error'),
                    description: t('onboarding.error'),
                    variant: "destructive",
                });
            }
        } catch (error: any) {
            toast({
                title: t('common.error'),
                description: error?.data?.message || t('onboarding.error'),
                variant: "destructive",
            });
        }
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
                            <CardTitle>{t('onboarding.reviewSetup')}</CardTitle>
                            <CardDescription>{t('onboarding.verifyInformation')}</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <div>
                                <h3 className="font-medium text-lg">{t('onboarding.administratorAccount')}</h3>
                                <dl className="mt-2 space-y-1">
                                    <div>
                                        <dt className="text-sm text-muted-foreground">{t('onboarding.name')}</dt>
                                        <dd>{adminData.firstName} {adminData.lastName}</dd>
                                    </div>
                                    <div>
                                        <dt className="text-sm text-muted-foreground">{t('onboarding.email')}</dt>
                                        <dd>{adminData.email}</dd>
                                    </div>
                                </dl>
                            </div>
                            <div>
                                <h3 className="font-medium text-lg">{t('onboarding.roles')} ({roles.length})</h3>
                                <ul className="mt-2 space-y-1">
                                    {roles.map((role, index) => (
                                        <li key={index}>{role.name}</li>
                                    ))}
                                </ul>
                            </div>
                            <div>
                                <h3 className="font-medium text-lg">{t('onboarding.teamInformation')}</h3>
                                <dl className="mt-2 space-y-1">
                                    <div>
                                        <dt className="text-sm text-muted-foreground">{t('onboarding.teamName')}</dt>
                                        <dd>{serviceData.teamName}</dd>
                                    </div>
                                    <div>
                                        <dt className="text-sm text-muted-foreground">{t('onboarding.language')}</dt>
                                        <dd>{serviceData.language}</dd>
                                    </div>
                                </dl>
                            </div>
                        </CardContent>
                        <CardFooter className="flex justify-between">
                            <Button type="button" variant="outline" onClick={handleSummaryBack}>
                                {t('common.back')}
                            </Button>
                            <Button onClick={handleSubmit}>
                                {t('onboarding.completeSetup')}
                            </Button>
                        </CardFooter>
                    </Card>
                )
        }
    }

    return (
        <div className="flex min-h-screen">
            <OnboardingSidebar currentStep={step} className='w-1/4' />
            <div className="flex-1 p-6">
                <h1 className="text-3xl font-bold mb-6 text-center">{t('onboarding.setupAccount')}</h1>
                {renderStep()}
            </div>
        </div>
    )
}

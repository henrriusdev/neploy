import * as React from 'react'
import { cn } from '@/lib/utils'

interface StepInfo {
    title: string;
    description: string;
}

interface OnboardingSidebarProps {
    currentStep: 'provider' | 'data' | 'roles' | 'service' | 'summary';
    className?: string;
}

const stepInfo: Record<string, StepInfo> = {
    provider: {
        title: 'Choose Provider',
        description: 'Select your preferred authentication provider or continue with manual setup.',
    },
    data: {
        title: 'Personal Information',
        description: 'Fill in your basic information to set up your account.',
    },
    roles: {
        title: 'Role Management',
        description: 'Define roles and permissions for your organization.',
    },
    service: {
        title: 'Service Setup',
        description: 'Configure your team and service settings.',
    },
    summary: {
        title: 'Review',
        description: 'Review your information before completing the setup.',
    },
}

export function OnboardingSidebar({ currentStep, className }: OnboardingSidebarProps) {
    return (
        <div className={cn("w-80 bg-muted p-6 flex flex-col gap-8", className)}>
            <div>
                <h2 className="text-2xl font-semibold mb-2">Welcome to Neploy</h2>
                <p className="text-muted-foreground">
                    Complete these steps to get started with your account
                </p>
            </div>

            <div className="space-y-6">
                {Object.entries(stepInfo).map(([step, info]) => {
                    const isActive = currentStep === step
                    const isPast = Object.keys(stepInfo).indexOf(currentStep) > 
                                 Object.keys(stepInfo).indexOf(step)

                    return (
                        <div key={step} className="flex gap-4 items-start">
                            <div className={cn(
                                "w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium",
                                isActive ? "bg-primary text-primary-foreground" : 
                                isPast ? "bg-primary/20 text-primary" : 
                                "bg-muted-foreground/20 text-muted-foreground"
                            )}>
                                {isPast ? '✓' : Object.keys(stepInfo).indexOf(step) + 1}
                            </div>
                            <div>
                                <h3 className={cn(
                                    "font-medium mb-1",
                                    isActive ? "text-primary" : 
                                    isPast ? "text-muted-foreground" : 
                                    "text-muted-foreground/70"
                                )}>
                                    {info.title}
                                </h3>
                                <p className="text-sm text-muted-foreground">
                                    {info.description}
                                </p>
                            </div>
                        </div>
                    )
                })}
            </div>

            <div className="mt-auto">
                <h3 className="font-medium mb-2">Need help?</h3>
                <p className="text-sm text-muted-foreground">
                    If you're having trouble setting up your account, please check our{' '}
                    <a href="#" className="text-primary hover:underline">documentation</a>{' '}
                    or contact{' '}
                    <a href="#" className="text-primary hover:underline">support</a>.
                </p>
            </div>
        </div>
    )
}

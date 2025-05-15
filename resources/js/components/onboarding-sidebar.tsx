import * as React from 'react'
import { cn } from '@/lib/utils'
import {LanguageSelector} from "@/components/forms";
import {useTranslation} from "react-i18next";
import {ThemeSwitcher} from "@/components/theme-switcher";

interface StepInfo {
    title: string;
    description: string;
}

interface OnboardingSidebarProps {
    currentStep: 'provider' | 'data' | 'roles' | 'service' | 'summary';
    className?: string;
}


export function OnboardingSidebar({ currentStep, className }: OnboardingSidebarProps) {
    const {t, i18n} = useTranslation();
    const stepInfo: Record<string, StepInfo> = {
        provider: {
            title: t('onboarding.provider.title'),
            description: t('onboarding.provider.description'),
        },
        data: {
            title: t('onboarding.data.title'),
            description: t('onboarding.data.description'),
        },
        roles: {
            title: t('onboarding.roles.title'),
            description: t('onboarding.roles.description'),
        },
        service: {
            title: t('onboarding.service.title'),
            description: t('onboarding.service.description'),
        },
        summary: {
            title: t('onboarding.summary.title'),
            description: t('onboarding.summary.description'),
        },
    }
const handleLanguageChange = (value: string) => {
    i18n.changeLanguage(value);
};
    return (
        <div className={cn("w-80 bg-muted p-6 flex flex-col gap-8", className)}>
            <div className="p-2">
                <h2 className="text-2xl font-semibold mb-2">{t('onboarding.welcome')}</h2>
                <p className="text-muted-foreground">
                    {t('onboarding.description')}
                </p>
            </div>

            <div className="space-y-6 flex flex-col justify-center items-start h-full">
                {Object.entries(stepInfo).map(([step, info]) => {
                    const isActive = currentStep === step
                    const isPast = Object.keys(stepInfo).indexOf(currentStep) > 
                                 Object.keys(stepInfo).indexOf(step)

                    return (
                        <div key={step} className="flex gap-4 items-start">
                            <div className={cn(
                                step === "provider" ? "px-4 py-1" : "",
                                "w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium",
                                isActive ? "bg-primary text-primary-foreground" :
                                isPast ? "bg-primary/20 text-primary" : 
                                "bg-muted-foreground/20 text-muted-foreground"
                            )}>
                                {isPast ? '✓' : Object.keys(stepInfo).indexOf(step) + 1 }
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
                <h3 className="font-medium mb-2">{t('onboarding.help.title')}</h3>
                <p className="text-sm text-muted-foreground">
                    {t('onboarding.help.firstDescription')}{' '}
                    <a href="#" className="text-primary hover:underline">{t('onboarding.help.documentation')}</a>{' '}
                    {t('onboarding.help.secondDescription')}{' '}
                    <a href="#" className="text-primary hover:underline">{t('onboarding.help.support')}</a>.
                </p>
                <p className="text-sm text-muted-foreground capitalize mt-5 mb-1">{t('language')}</p>
                <LanguageSelector />
                <ThemeSwitcher className="mt-2" />
            </div>
        </div>
    )
}

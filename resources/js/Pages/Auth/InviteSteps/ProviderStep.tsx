import * as React from 'react'
import { Button } from "@/components/ui/button"
import { GitHubLogoIcon } from "@radix-ui/react-icons"
import { router } from '@inertiajs/react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface Props {
    onNext: () => void;
    token?: string;
}

export default function ProviderStep({ token, onNext }: Props) {
    const getOAuthUrl = (provider: string) => {
        if (token) {
            return `/auth/${provider}?state=${token}`;
        }
        return `/auth/${provider}`;
    };

    return (
        <Card className="w-full max-w-md mx-auto">
            <CardHeader>
                <CardTitle>Link Your Account</CardTitle>
                <CardDescription>
                    Choose how you want to authenticate with Neploy
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                <Button
                    variant="outline"
                    className="w-full" 
                    onClick={() => window.location.replace(getOAuthUrl('github'))}
                >
                    <GitHubLogoIcon className="mr-2 h-4 w-4" />
                    Continue with GitHub
                </Button>
                <Button 
                    variant="outline" 
                    className="w-full"
                    onClick={() => window.location.replace(getOAuthUrl('gitlab'))}
                >
                    <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24">
                        <path fill="currentColor" d="M21.94 13.11l-1.05-3.22c0-.03-.01-.06-.02-.09l-2.11-6.48a.859.859 0 0 0-.8-.57c-.36 0-.68.25-.79.58l-2.01 6.19H8.84L6.83 3.33a.851.851 0 0 0-.79-.58c-.37 0-.69.25-.8.58L3.13 9.82v.01l-1.05 3.22c-.16.5.01 1.04.44 1.34l9.22 6.71c.17.12.39.12.56 0l9.22-6.71c.43-.3.6-.84.44-1.34M8.15 10.45l2.57 7.91l-6.17-7.91h3.6m1.22 0h5.26l-2.63 8.1l-2.63-8.1m6.48 0h3.6l-6.17 7.91l2.57-7.91m-13.5 2.62l1.17-3.6l5.88 7.54l-7.05-3.94M20.65 13.07l-7.05 3.94l5.88-7.54l1.17 3.6Z"/>
                    </svg>
                    Continue with GitLab
                </Button>
                <div className="relative">
                    <div className="absolute inset-0 flex items-center">
                        <span className="w-full border-t" />
                    </div>
                    <div className="relative flex justify-center text-xs uppercase">
                        <span className="bg-background px-2 text-muted-foreground">
                            Or continue with
                        </span>
                    </div>
                </div>

                <Button 
                    variant="outline" 
                    className="w-full"
                    onClick={onNext}
                >
                    Continue with Email
                </Button>
            </CardContent>
        </Card>
    )
}

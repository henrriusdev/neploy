import * as React from 'react'
import { router } from '@inertiajs/react'
import { Button } from "@/components/ui/button"
import { useToast } from '@/hooks/use-toast'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"
import { AcceptInviteProps } from '@/types/props'

export default function AcceptInvite({ token, expired, alreadyAccepted, provider }: AcceptInviteProps) {
    const { toast } = useToast()
    const [isLoading, setIsLoading] = React.useState(false)

    if (expired) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50">
                <Card className="max-w-md w-full">
                    <CardHeader>
                        <CardTitle>Invitation Expired</CardTitle>
                        <CardDescription>
                            This invitation has expired. Please request a new invitation from your team administrator.
                        </CardDescription>
                    </CardHeader>
                </Card>
            </div>
        )
    }

    if (alreadyAccepted) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50">
                <Card className="max-w-md w-full">
                    <CardHeader>
                        <CardTitle>Invitation Already Accepted</CardTitle>
                        <CardDescription>
                            This invitation has already been used. Please log in to access your account.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <Button 
                            className="w-full"
                            onClick={() => router.visit('/login')}
                        >
                            Go to Login
                        </Button>
                    </CardContent>
                </Card>
            </div>
        )
    }

    const handleAcceptInvite = () => {
        setIsLoading(true)
        router.post('/users/accept-invite', 
            { token },
            {
                onSuccess: () => {
                    // Si no hay provider, redirigir al login para que elija uno
                    if (!provider) {
                        router.visit('/login', {
                            data: { invite_accepted: true }
                        })
                    } else {
                        // Si ya tiene provider, ir directo al onboarding
                        router.visit('/onboard', {
                            data: { invite_accepted: true }
                        })
                    }
                },
                onError: (errors) => {
                    toast({
                        title: "Error",
                        description: errors.message || "Failed to accept invitation",
                        variant: "destructive",
                    })
                },
                onFinish: () => setIsLoading(false)
            }
        )
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
            <Card className="max-w-md w-full">
                <CardHeader>
                    <CardTitle>Accept Team Invitation</CardTitle>
                    <CardDescription>
                        Click below to join your team on Neploy
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <Button
                        className="w-full"
                        onClick={handleAcceptInvite}
                        disabled={isLoading}
                    >
                        {isLoading ? "Accepting..." : "Accept Invitation"}
                    </Button>
                </CardContent>
            </Card>
        </div>
    )
}

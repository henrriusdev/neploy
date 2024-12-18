import * as React from 'react'
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Controller } from "react-hook-form"

const serviceSchema = z.object({
    teamName: z.string().min(1, "Team name is required"),
    logo: z.string().min(1, "Logo URL is required"),
})

interface Props {
    onNext: (data: z.infer<typeof serviceSchema>) => void;
    onBack: () => void;
}

export default function ServiceStep({ onNext, onBack }: Props) {
    const form = useForm<z.infer<typeof serviceSchema>>({
        resolver: zodResolver(serviceSchema),
        defaultValues: {
            teamName: "",
            logo: "",
        },
    })

    return (
        <Card className="w-full max-w-screen-md mx-auto">
            <CardHeader>
                <CardTitle>Service Metadata</CardTitle>
                <CardDescription>Set up your team information</CardDescription>
            </CardHeader>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onNext)}>
                    <CardContent className="space-y-4">
                        <FormField
                            control={form.control}
                            name="teamName"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Team Name</FormLabel>
                                    <FormControl>
                                        <Controller
                                            control={form.control}
                                            name="teamName"
                                            render={({ field }) => (
                                                <Input {...field} />
                                            )}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="logo"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Logo URL</FormLabel>
                                    <FormControl>
                                        <Controller
                                            control={form.control}
                                            name="logo"
                                            render={({ field }) => (
                                                <Input {...field} />
                                            )}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                    </CardContent>
                    <CardFooter className="flex justify-between">
                        <Button type="button" variant="outline" onClick={onBack}>
                            Back
                        </Button>
                        <Button type="submit">
                            Next
                        </Button>
                    </CardFooter>
                </form>
            </Form>
        </Card>
    )
}

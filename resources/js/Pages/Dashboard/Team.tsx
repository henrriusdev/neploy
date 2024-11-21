import * as React from 'react'
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { PlusCircle } from 'lucide-react'
import DashboardLayout from '@/components/Layouts/DashboardLayout'

interface TeamMember {
    id: string
    username: string
    email: string
    firstName: string
    lastName: string
    roles: Array<{
        name: string
        description: string
        icon: string
        color: string
    }>
    techStacks: Array<{
        name: string
        description: string
    }>
}

const defaultTeam: TeamMember[] = [
    {
        id: "1",
        username: "johndoe",
        email: "john@example.com",
        firstName: "John",
        lastName: "Doe",
        roles: [
            {
                name: "Admin",
                description: "Administrator",
                icon: "crown",
                color: "#FFB020"
            }
        ],
        techStacks: [
            {
                name: "React",
                description: "Frontend Framework"
            }
        ]
    }
]

interface TeamProps {
    user?: {
        name: string
        email: string
        avatar: string
    }
    teamName?: string
    logoUrl?: string
    team?: TeamMember[]
}

function Team({
    user,
    teamName,
    logoUrl,
    team = defaultTeam,
}: TeamProps) {
    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Team Members</h2>
                    <p className="text-muted-foreground">
                        Manage your team members and their access levels
                    </p>
                </div>
                <Button>
                    <PlusCircle className="mr-2 h-4 w-4" />
                    Add Member
                </Button>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Team Overview</CardTitle>
                    <CardDescription>
                        A list of all team members including their roles and tech stacks.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Member</TableHead>
                                <TableHead>Roles</TableHead>
                                <TableHead>Tech Stacks</TableHead>
                                <TableHead>Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {team.map((member) => (
                                <TableRow key={member.id}>
                                    <TableCell>
                                        <div className="flex items-center space-x-4">
                                            <Avatar>
                                                <AvatarImage src={`https://avatar.vercel.sh/${member.username}`} />
                                                <AvatarFallback>
                                                    {member.firstName[0]}{member.lastName[0]}
                                                </AvatarFallback>
                                            </Avatar>
                                            <div>
                                                <div className="font-medium">{member.firstName} {member.lastName}</div>
                                                <div className="text-sm text-muted-foreground">{member.email}</div>
                                            </div>
                                        </div>
                                    </TableCell>
                                    <TableCell>
                                        <div className="flex flex-wrap gap-2">
                                            {member.roles.map((role) => (
                                                <Badge
                                                    key={role.name}
                                                    variant="secondary"
                                                    style={{
                                                        backgroundColor: role.color + '20',
                                                        color: role.color
                                                    }}
                                                >
                                                    {role.name}
                                                </Badge>
                                            ))}
                                        </div>
                                    </TableCell>
                                    <TableCell>
                                        <div className="flex flex-wrap gap-2">
                                            {member.techStacks.map((tech) => (
                                                <Badge key={tech.name} variant="outline">
                                                    {tech.name}
                                                </Badge>
                                            ))}
                                        </div>
                                    </TableCell>
                                    <TableCell>
                                        <div className="flex space-x-2">
                                            <Button variant="outline" size="sm">Edit</Button>
                                            <Button variant="destructive" size="sm">Remove</Button>
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>
        </div>
    )
}

// Add the layout property to the component
Team.layout = (page: React.ReactNode, props: TeamProps) => (
    <DashboardLayout {...props}>
        {page}
    </DashboardLayout>
)

export default Team
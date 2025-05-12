"use client"

import {useEffect, useState} from "react"
import {Pencil} from "lucide-react"
import {Button} from "@/components/ui/button"
import {Checkbox} from "@/components/ui/checkbox"
import {Label} from "@/components/ui/label"
import {ScrollArea} from "@/components/ui/scroll-area"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {TechStack} from "@/types";

interface TechAssignmentDialogProps {
  userId: string
  allTechStacks: TechStack[]
  selectedTechIds: string[]
  onSave: (userId: string, techIds: string[]) => void
}

export function TechAssignmentDialog({userId, allTechStacks, selectedTechIds, onSave}: TechAssignmentDialogProps) {
  const [open, setOpen] = useState(false)
  const [selectedTechs, setSelectedTechs] = useState<string[]>([])
  const [initialSelection, setInitialSelection] = useState<string[]>([])

  // Initialize selected techs when dialog opens
  useEffect(() => {
    if (open) {
      setSelectedTechs([...selectedTechIds])
      setInitialSelection([...selectedTechIds])
    }
  }, [open, selectedTechIds])

  // Check if there are any changes compared to initial selection
  const hasChanges = () => {
    if (selectedTechs.length !== initialSelection.length) return true
    return !selectedTechs.every((tech) => initialSelection.includes(tech))
  }

  const handleToggleTech = (techId: string) => {
    setSelectedTechs((prev) => (prev.includes(techId) ? prev.filter((id) => id !== techId) : [...prev, techId]))
  }

  const handleSave = () => {
    onSave(userId, selectedTechs)
    setOpen(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="ghost" size="icon" className="h-8 w-8">
          <Pencil className="h-4 w-4"/>
          <span className="sr-only">Edit Technologies</span>
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Assign Technologies</DialogTitle>
          <DialogDescription>Select the technologies this team member is proficient with.</DialogDescription>
        </DialogHeader>
        <ScrollArea className="max-h-[60vh] mt-4 pr-4">
          <div className="gap-4 grid grid-cols-3 pb-4">
            {allTechStacks.map((tech) => (
              <div key={tech.id} className="flex items-center space-x-2">
                <Checkbox
                  id={`tech-${tech.id}`}
                  checked={selectedTechs.includes(tech.id)}
                  onCheckedChange={() => handleToggleTech(tech.id)}
                />
                <Label
                  htmlFor={`tech-${tech.id}`}
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  {tech.name}
                </Label>
              </div>
            ))}
          </div>
        </ScrollArea>
        <DialogFooter className="flex items-center justify-between sm:justify-between">
          <Button type="button" variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button type="button" onClick={handleSave} disabled={!hasChanges()}>
            Save Changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

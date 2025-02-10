import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { getScheduleTemplate } from "@/lib/api";
import { DayOfWeek } from "@/lib/const";
import { cn } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import { User } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  scheduleTemplateId: number;
}

export default function ShowScheduleTemplateDetailDialog({
  open,
  onOpenChange,
  scheduleTemplateId,
}: Props) {
  const {
    data: scheduleTemplate,
    isPending,
    isError,
    error,
  } = useQuery({
    queryKey: ["schedule-template", scheduleTemplateId],
    queryFn: () =>
      getScheduleTemplate(scheduleTemplateId).then((res) => res.data.data),
  });

  if (isPending) {
    return null;
  }

  if (isError) {
    toast.error(error.message);
    return null;
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogHeader>{scheduleTemplate?.meta.name}</DialogHeader>
          <DialogDescription
            className={cn(
              scheduleTemplate?.meta.description.length <= 0 &&
                "text-muted-foreground"
            )}
          >
            {scheduleTemplate?.meta.description.length > 0
              ? scheduleTemplate?.meta.description
              : "此模板无描述"}
          </DialogDescription>
        </DialogHeader>
        <Tabs defaultValue={DayOfWeek[0].key.toString()}>
          <TabsList className={cn("grid", `grid-cols-${DayOfWeek.length}`)}>
            {DayOfWeek.map((day) => (
              <TabsTrigger key={day.key} value={day.key.toString()}>
                {day.label}
              </TabsTrigger>
            ))}
          </TabsList>
          {DayOfWeek.map((day) => (
            <TabsContent key={day.key} value={day.key.toString()}>
              <div>
                {scheduleTemplate?.shifts.some((shift) =>
                  shift.applicableDays.includes(day.key)
                ) ? (
                  <div className="grid gap-2">
                    {scheduleTemplate.shifts
                      .filter((shift) => shift.applicableDays.includes(day.key))
                      .map((s) => (
                        <div
                          key={s.startTime}
                          className="text-sm text-muted-foreground border-border border rounded-md py-4 flex justify-between px-4"
                        >
                          <span>
                            {s.startTime} - {s.endTime}
                          </span>
                          <span className="flex items-center gap-2">
                            {s.requiredAssistantNumber}名助理
                            <User className="w-4 h-4" />
                          </span>
                        </div>
                      ))}
                  </div>
                ) : (
                  <div className="text-sm text-muted-foreground border-border border rounded-md border-dashed py-4 text-center">
                    当天没有班次
                  </div>
                )}
              </div>
            </TabsContent>
          ))}
        </Tabs>
        <DialogFooter>
          <DialogClose asChild>
            <Button type="button" variant="secondary">
              关闭
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

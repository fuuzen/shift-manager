import { ScheduleTemplate } from "@/lib/types";
import { create } from "zustand";

interface State {
  open: boolean;
  setOpen: (open: boolean) => void;
  scheduleTemplate: ScheduleTemplate | null;
  setScheduleTemplate: (scheduleTemplate: ScheduleTemplate) => void;
}

const useDeleteScheduleTemplateDialogStore = create<State>((set) => ({
  open: false,
  setOpen: (open) => set({ open }),
  scheduleTemplate: null,
  setScheduleTemplate: (scheduleTemplate) => set({ scheduleTemplate }),
}));

export default useDeleteScheduleTemplateDialogStore;

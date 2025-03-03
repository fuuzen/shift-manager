package handler

import (
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sysu-ecnc-dev/shift-manager/backend/internal/domain"
)

func (h *Handler) icalHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.internalServerError(w, r, err)
	}

	plans, err := h.repository.GetAllSchedulePlans()
	if err != nil {
		h.internalServerError(w, r, err)
	}

	var icsContent string = `BEGIN:VCALENDAR
CALSCALE:GREGORIAN
PRODID:SYSU-ECNC
VERSION:2.0
X-APPLE-CALENDAR-COLOR:#CC73E1
X-WR-CALNAME:ECNC值班日程
BEGIN:VTIMEZONE
TZID:Asia/Shanghai
BEGIN:STANDARD
DTSTART:19890917T020000
RRULE:FREQ=YEARLY;UNTIL=19910914T170000Z;BYMONTH=9;BYDAY=3SU
TZNAME:GMT+8
TZOFFSETFROM:+0900
TZOFFSETTO:+0800
END:STANDARD
BEGIN:DAYLIGHT
DTSTART:19910414T020000
RDATE:19910414T020000
TZNAME:GMT+8
TZOFFSETFROM:+0800
TZOFFSETTO:+0900
END:DAYLIGHT
END:VTIMEZONE`

	for _, plan := range plans {
		res, err := h.repository.GetSchedulingResultBySchedulePlanID(plan.ID)
		if err != nil {
			h.internalServerError(w, r, err)
		}

		template, err := h.repository.GetScheduleTemplate(plan.ScheduleTemplateID)
		if err != nil {
			h.internalServerError(w, r, err)
		}

		shifts := template.Shifts
		shiftMap := make(map[int64]*domain.ScheduleTemplateShift)
		for _, shift := range shifts {
			shiftMap[shift.ID] = &shift
		}

		for _, resultShift := range res.Shifts {
			for _, item := range resultShift.Items {
				if (item.PrincipalID != nil && *item.PrincipalID == userID) || slices.Contains(item.AssistantIDs, userID) {
					icsContent += "\nBEGIN:VEVENT"
					if *item.PrincipalID == userID {
						icsContent += "\nDESCRIPTION:负责人岗值班"
					} else {
						icsContent += "\nDESCRIPTION:普通助理岗值班"
					}
					templateShift := shiftMap[resultShift.ShiftID]
					firstTime := plan.ActiveStartTime
					for {
						if firstTime.Weekday() == time.Weekday(item.Day-1) {
							break
						} else {
							firstTime = firstTime.Add(24 * time.Hour)
						}
					}
					startTime, err := time.Parse("15:04:05", templateShift.StartTime)
					if err != nil {
						h.internalServerError(w, r, err)
					}
					endTime, err := time.Parse("15:04:05", templateShift.EndTime)
					if err != nil {
						h.internalServerError(w, r, err)
					}
					icsContent += "\nDTEND;TZID=Asia/Shanghai:" + firstTime.Add(time.Duration(endTime.Hour())*time.Hour).Format("20060102T150405")
					icsContent += "\nDTSTAMP:" + plan.CreatedAt.Format("20060102T150405") + "Z"
					icsContent += "\nDTSTART;TZID=Asia/Shanghai:" + firstTime.Add(time.Duration(startTime.Hour())*time.Hour).Format("20060102T150405")
					icsContent += "\nLAST-MODIFIED:" + time.Now().Format("20060102T150405") + "Z"
					icsContent += "\nLOCATION:中山大学东校园"
					icsContent += "\nRRULE:FREQ=WEEKLY;UNTIL=" + plan.ActiveEndTime.Format("20060102T150405") + "Z"
					icsContent += "\nSEQUENCE:0"
					icsContent += "\nSUMMARY:ECNC值班日程"
					icsContent += "\nTRANSP:OPAQUE"
					icsContent += "\nUID:SYSU-ECNC-" + strconv.FormatInt(userID, 16) + "-" + strconv.FormatInt(resultShift.ShiftID, 16) + "-" + strconv.FormatInt(int64(item.Day), 16)
					icsContent += "\nEND:VEVENT"
				}
			}
		}
	}
	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", "inline; filename=calendar.ics")
	w.Write(([]byte)(icsContent + "\nEND:VCALENDAR"))
}

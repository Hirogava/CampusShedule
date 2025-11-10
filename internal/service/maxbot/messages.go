package maxbot

import (
	"fmt"
	"strings"

	dbModels "github.com/Hirogava/CampusShedule/internal/models/db"
)

func CreateScheduledMessage(days []dbModels.Day) string {
	if len(days) == 0 {
		return "📭 На этой неделе занятий нет!"
	}

	var sb strings.Builder
	sb.WriteString("📅 <b>Расписание на неделю:</b>\n\n")

	for _, day := range days {
		sb.WriteString(fmt.Sprintf("🗓️ <b>%s</b>\n", day.WeekDay))

		if len(day.Lessons) == 0 {
			sb.WriteString("  ❌ Занятий нет\n\n")
			continue
		}

		for _, lesson := range day.Lessons {
			sb.WriteString(fmt.Sprintf(
				"  ⏰ <b>%s–%s</b>\n  📘 %s (%s)\n",
				lesson.StartTime,
				lesson.EndTime,
				lesson.Type.TypeToEmoji(),
				lesson.Type.String(),
			))
			if lesson.Teacher != "" {
				sb.WriteString(fmt.Sprintf("  👨‍🏫 %s\n", lesson.Teacher))
			}
			if lesson.Room != "" {
				sb.WriteString(fmt.Sprintf("  🚪 %s\n", lesson.Room))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

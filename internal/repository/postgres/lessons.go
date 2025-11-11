package postgres

import (
	"time"
	dbModels "github.com/Hirogava/CampusShedule/internal/models/db"
)

func (m *Manager) GetUpcomingLessons(target time.Time) ([]dbModels.LessonNotify, error) {
	rows, err := m.Conn.Query(`
		SELECT l.subject, l.teacher, l.room, l.start_time, l.type, u.chat_id
		FROM schedule l
		JOIN users u ON u.group_id = l.group_id
		WHERE l.start_time BETWEEN $1 AND $2
	`, target.Add(-1*time.Minute), target.Add(1*time.Minute))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []dbModels.LessonNotify
	for rows.Next() {
		var l dbModels.LessonNotify
		if err := rows.Scan(&l.Name, &l.Teacher, &l.Room, &l.StartTime, &l.Type, &l.ChatID); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

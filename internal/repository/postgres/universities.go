package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	dbErrors "github.com/Hirogava/CampusShedule/internal/errors/db"
	dbModels "github.com/Hirogava/CampusShedule/internal/models/db"
)

var weekdays = map[time.Weekday]string{
	time.Monday:    "Понедельник",
	time.Tuesday:   "Вторник",
	time.Wednesday: "Среда",
	time.Thursday:  "Четверг",
	time.Friday:    "Пятница",
	time.Saturday:  "Суббота",
	time.Sunday:    "Воскресенье",
}

func (manager *Manager) GetUniversities() ([]dbModels.University, error) {
	var universities []dbModels.University

	rows, err := manager.Conn.Query("SELECT id, name FROM universities WHERE schedule = true")
	if err != nil {
		logger.Logger.Error("Failed to get universities", "error", err.Error())
		return nil, err
	}

	for rows.Next() {
		var university dbModels.University

		if err := rows.Scan(&university.ID, &university.Name); err != nil {
			logger.Logger.Error("Failed to scan university", "error", err.Error())
			return nil, err
		}

		universities = append(universities, university)
	}

	return universities, nil
}

func (manager *Manager) HasUserUniversity(userID int64) (bool, error) {
	var universityID sql.NullInt64

	if err := manager.Conn.QueryRow("SELECT university_id FROM users WHERE chat_id = $1", userID).Scan(&universityID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			logger.Logger.Error("Failed to check if user has university", "error", err.Error())
			return false, err
		}
	}

	if universityID.Valid {
		return true, nil
	}

	return false, nil
}

func (manager *Manager) GetUniversity(id int) (dbModels.University, error) {
	var university dbModels.University

	if err := manager.Conn.QueryRow("SELECT id, name FROM universities WHERE id = $1", id).Scan(&university.ID, &university.Name); err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Warn("University not found", "id", id)
			return dbModels.University{}, dbErrors.ErrUniversityNotFound
		} else {
			logger.Logger.Error("Failed to get university", "error", err.Error())
			return dbModels.University{}, err
		}
	}

	return university, nil
}

func (manager *Manager) SetUserUniversity(userID int64, universityID int) (string, error) {
	var apiURL string
    
    err := manager.Conn.QueryRow(
        `WITH upsert AS (
			INSERT INTO users (chat_id, university_id)
			VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE SET university_id = EXCLUDED.university_id
			RETURNING university_id
		)
		SELECT api_url FROM universities WHERE id = (SELECT university_id FROM upsert);
		`,userID, universityID,
    ).Scan(&apiURL)
	if err != nil {
		logger.Logger.Error("Failed to set user university", "error", err.Error())
		return "", err
	}

	return apiURL, nil
}

func (manager *Manager) SetUserGroup(userID int64, group int) error {
	_, err := manager.Conn.Exec(`
		UPDATE users SET group_id = $1 WHERE chat_id = $2
	`, group, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Warn("User not found", "id", userID)
			return dbErrors.ErrUserNotFound
		} else {
			logger.Logger.Error("Failed to set user group", "error", err.Error())
			return err
		}
	}

	return nil
}

func (manager *Manager) GetUniversityGroups(universityID int) ([]dbModels.Group, error) {
	rows, err := manager.Conn.Query(`
		SELECT g.id, g.name
		FROM groups g
		JOIN universities_groups ug ON g.id = ug.group_id
		WHERE ug.university_id = $1
		ORDER BY g.name`, universityID)
	if err != nil {
		logger.Logger.Error("Failed to get university groups", "error", err.Error())
		return nil, err
	}

	var groups []dbModels.Group
	for rows.Next() {
		var group dbModels.Group

		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			logger.Logger.Error("Failed to scan group", "error", err.Error())
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (manager *Manager) GetUserGroup(userID int64) (int, error) {
	var groupID sql.NullInt64

	err := manager.Conn.QueryRow("SELECT group_id FROM users WHERE id = $1", userID).Scan(&groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, dbErrors.ErrUserNotFound
		}
		return 0, err
	}

	if !groupID.Valid {
		return 0, dbErrors.ErrUserNotFound
	}

	return int(groupID.Int64), nil
}

func (manager *Manager) GetWeekSchedule(ctx context.Context, groupID int) ([]dbModels.Day, error) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // воскресенье → 7
	}
	startOfWeek := now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	query := `
		SELECT teacher, room, start_time, end_time, date, type
		FROM lessons
		WHERE group_id = $1
		  AND date >= $2
		  AND date < $3
		ORDER BY date, start_time;
	`

	rows, err := manager.Conn.QueryContext(ctx, query, groupID, startOfWeek, endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении расписания: %w", err)
	}
	defer rows.Close()

	daysMap := make(map[string][]dbModels.Lesson)

	for rows.Next() {
		var lesson dbModels.Lesson
		var date time.Time
		var lessonType string

		if err := rows.Scan(
			&lesson.Teacher,
			&lesson.Room,
			&lesson.StartTime,
			&lesson.EndTime,
			&date,
			&lessonType,
		); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}

		lesson.Date = date.Format("2006-01-02")
		lesson.DateOfWeek = weekdays[date.Weekday()]
		lesson.Type = dbModels.LessonType(lessonType)

		daysMap[lesson.DateOfWeek] = append(daysMap[lesson.DateOfWeek], lesson)
	}

	var result []dbModels.Day
	order := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}

	for _, wd := range order {
		wdStr := wd.String()
		if lessons, ok := daysMap[wdStr]; ok {
			result = append(result, dbModels.Day{
				WeekDay: wdStr,
				Lessons: lessons,
			})
		}
	}

	return result, nil
}

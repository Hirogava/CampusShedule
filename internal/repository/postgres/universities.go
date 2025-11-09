package postgres

import (
	"database/sql"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	dbErrors "github.com/Hirogava/CampusShedule/internal/errors/db"
	dbModels "github.com/Hirogava/CampusShedule/internal/models/db"
)

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

	if err := manager.Conn.QueryRow("SELECT university_id FROM users WHERE id = $1", userID).Scan(&universityID); err != nil {
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
        `INSERT INTO users (id, university_id) 
         VALUES ($1, $2) 
         ON CONFLICT (id) DO UPDATE SET university_id = $2
         RETURNING (SELECT api_url FROM universities WHERE id = $2)`,
        userID, universityID,
    ).Scan(&apiURL)
	if err != nil {
		logger.Logger.Error("Failed to set user university", "error", err.Error())
		return "", err
	}

	return apiURL, nil
}

func (manager *Manager) SetUserGroup(userID int64, group string) error {
	_, err := manager.Conn.Exec(`
		UPDATE users SET group_name = $1 WHERE id = $2
	`)
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

package maxbot

import (
	"context"
	"fmt"
	"log"
	"time"

	service "github.com/Hirogava/CampusShedule/internal/service/maxbot"
	"github.com/Hirogava/CampusShedule/internal/repository/postgres"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
)

type Scheduler struct {
	api       *maxbot.Api
	manager   *postgres.Manager
	interval  time.Duration
	ahead     time.Duration
}

func NewScheduler(api *maxbot.Api, manager *postgres.Manager, interval, ahead time.Duration) *Scheduler {
	return &Scheduler{
		api:       api,
		manager:   manager,
		interval:  interval,
		ahead:     ahead,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Println("[Scheduler] Started")

	for {
		select {
		case <-ctx.Done():
			log.Println("[Scheduler] Stopped")
			return
		case <-ticker.C:
			s.checkAndNotify(ctx)
		}
	}
}

func (s *Scheduler) checkAndNotify(ctx context.Context) {
	now := time.Now()
	target := now.Add(s.ahead)

	lessons, err := s.manager.GetUpcomingLessons(target)
	if err != nil {
		log.Printf("[Scheduler] Error getting lessons: %v", err)
		return
	}

	for _, lesson := range lessons {
		diff := lesson.StartTime.Sub(now)
		if diff < s.ahead-1*time.Minute || diff > s.ahead+1*time.Minute {
			continue
		}

		msg := fmt.Sprintf(
			"⏰ Через %d минут начнётся %s по <b>%s</b> в аудитории %s\nПреподаватель: %s",
			int(s.ahead.Minutes()),
			lesson.Type,
			lesson.Name,
			lesson.Room,
			lesson.Teacher,
		)

		kb := service.CreateKeyboardForStart(s.api)
		_, err := s.api.Messages.Send(ctx,
			maxbot.NewMessage().
				AddKeyboard(kb).
				SetChat(lesson.ChatID).
				SetFormat("html").
				SetText(msg),
		)
		if err != nil {
			log.Printf("[Scheduler] Send error: %v", err)
			continue
		}
	}
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

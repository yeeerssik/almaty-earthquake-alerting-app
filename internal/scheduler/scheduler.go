package scheduler

import (
	"almeqapp/config"
	"almeqapp/internal/usecase"
	"time"
)

func StartScheduler(conf *config.Config, service *usecase.EarthquakeService) {
	ticker := time.NewTicker(time.Duration(conf.MinutesCount) * time.Minute)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			_ = service.NotifyAboutEQ()
		}
	}()
}

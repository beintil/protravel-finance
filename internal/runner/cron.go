package runner

import (
	"protravel-finance/internal/cron"
	"protravel-finance/pkg/logger"
)

func InitCronTasks(log logger.Logger, cronTasks ...cron.Cron) {
	log.Info("Starting cron jobs...")
	for _, task := range cronTasks {
		go task.Run()
	}
	log.Info("Cron jobs started")
}

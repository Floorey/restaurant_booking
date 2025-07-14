package job

import (
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"time"
)

func StartBackroundJobs(db *sqlx.DB) {
	c := cron.New()

	// Reminder
	c.AddFunc("0 * * * *", func() {
		sendReminders(db, 24*time.Hour)
	})
	c.AddFunc("30 * * * *", func() {
		sendReminders(db, 8*time.Hour)
	})
	// // Release expired bookings
	c.AddFunc("@every 30m", func() {
		expireUnconfirmed(db)
	})
	c.Start()
}

func sendReminders(db *sqlx.DB, delta time.Duration) {
	// ToDo: SMS sending Twilio
}

func expireUnconfirmed(db *sqlx.DB) {
	// ToDo: add Auto Expire
}

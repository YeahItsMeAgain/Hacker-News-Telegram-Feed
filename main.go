package main

import (
	"hn_feed/bot"
	"hn_feed/config"
	"hn_feed/db"
	"hn_feed/timer"
)

func main() {
	config.Init()
	db.Init()

	tgBot := bot.Init()
	go timer.ScheduleUpdates(tgBot)
	bot.Run(tgBot)
}

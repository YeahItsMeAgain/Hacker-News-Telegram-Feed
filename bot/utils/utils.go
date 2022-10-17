package utils

import (
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

func CreateBot(botToken string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func CreateReplyMarkup(rows ...telebot.Row) [][]telebot.ReplyButton {
	replyKeys := make([][]telebot.ReplyButton, 0, len(rows))
	for _, row := range rows {
		keys := make([]telebot.ReplyButton, 0, len(row))
		for _, btn := range row {
			btn := btn.Reply()
			if btn != nil {
				keys = append(keys, *btn)
			}
		}
		replyKeys = append(replyKeys, keys)
	}
	return replyKeys
}

func SendAndDelete(ctx telebot.Context, msg string) error {
	reply, err := ctx.Bot().Send(ctx.Recipient(), msg)
	if err != nil {
		return err
	}
	time.AfterFunc(time.Duration(10)*time.Second, func() {
		ctx.Bot().Delete(reply)
	})
	return nil
}

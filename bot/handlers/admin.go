package handlers

import (
	"hn_feed/bot/utils"
	"hn_feed/db"

	"gopkg.in/telebot.v3"
)

var (
	AdminBtnList = telebot.Btn{Text: "📚 List"}
	AdminMenu    = &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		ReplyKeyboard: utils.CreateReplyMarkup(
			telebot.Row{AdminBtnList},
		),
	}
)

func HandleAdminListChannels(ctx telebot.Context) error {
	var channels []db.Channel
	db.DB.Find(&channels)
	return ctx.Send(utils.StructsToString(channels))
}

package handlers

import (
	"hn_feed/bot/utils"
	"hn_feed/db"
	"hn_feed/db/models"
	db_utils "hn_feed/db/utils"

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
	var channels []models.Channel
	db.DB.Find(&channels)
	return ctx.Send(db_utils.StructsToString(channels))
}

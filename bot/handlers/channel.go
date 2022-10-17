package handlers

import (
	middlwares "hn_feed/bot/middlewares"

	"gopkg.in/telebot.v3"
)

func OnChannelRegister(ctx telebot.Context) error {
	ctx.Delete()
	return ctx.Send("🚀 Registered!\n Set the number of top posts you want to see with /set <count>.")
}

func OnChannelConfigureCount(ctx telebot.Context) error {
	payload := ctx.Get(middlwares.ChannelCommandPayloadKey)
	return ctx.Send(payload)
}

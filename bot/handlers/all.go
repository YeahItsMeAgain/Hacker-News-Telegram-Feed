package handlers

import (
	"fmt"
	"hn_feed/config"
	"log"

	"golang.org/x/exp/slices"
	"gopkg.in/telebot.v3"
)

func HandleStart(ctx telebot.Context) error {
	log.Printf("[*] %d : %s Started the bot.", ctx.Sender().ID, ctx.Sender().Username)
	if slices.Contains(config.Config.AdminIds, ctx.Sender().ID) {
		return ctx.Send(fmt.Sprintf("👋 Welcome %s!", ctx.Sender().FirstName), AdminMenu)
	}
	return ctx.Send(fmt.Sprintf("👋 Welcome %s!", ctx.Sender().FirstName))
}

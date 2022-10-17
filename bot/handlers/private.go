package handlers

import (
	"fmt"
	"hn_feed/config"
	"log"

	"golang.org/x/exp/slices"
	"gopkg.in/telebot.v3"
)

const START_MESSAGE = `👋 Welcome %s!

Add me to a channel with the following permissions:
- Post messages.
- Delete messages of others.

And send /register in that channel.
`

func HandleStart(ctx telebot.Context) error {
	log.Printf("[*] %d : %s Started the bot.", ctx.Sender().ID, ctx.Sender().Username)
	if slices.Contains(config.Config.AdminIds, ctx.Sender().ID) {
		return ctx.Send(fmt.Sprintf(START_MESSAGE, ctx.Sender().FirstName), AdminMenu)
	}
	return ctx.Send(fmt.Sprintf(START_MESSAGE, ctx.Sender().FirstName))
}

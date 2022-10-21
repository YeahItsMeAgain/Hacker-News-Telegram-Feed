package handlers

import (
	"fmt"
	"hn_feed/config"
	"log"

	"golang.org/x/exp/slices"
	"gopkg.in/telebot.v3"
)

const startMessage = `👋 Welcome %s!

Add me to a channel with the following permissions:
- Post messages.
- Delete messages of others.

And send /register in that channel.
`

func HandleStart(ctx telebot.Context) error {
	sender := ctx.Sender()
	log.Printf("[*] <%d : %s> Started the bot.", sender.ID, sender.Username)
	if slices.Contains(config.Config.AdminIds, sender.ID) {
		return ctx.Send(fmt.Sprintf(startMessage, sender.FirstName), AdminMenu)
	}
	return ctx.Send(fmt.Sprintf(startMessage, sender.FirstName))
}

package middlewares

import (
	"strings"

	"gopkg.in/telebot.v3"
)

const (
	ChannelCommandPayloadKey = "ChannelCommandPayloadKey"
)

func ChannelCommand(command string) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Update().ChannelPost == nil ||
				!strings.HasPrefix(c.Text(), command) {
				return nil
			}

			payload := strings.TrimPrefix(strings.TrimPrefix(c.Text(), command), " ")
			if payload != "" {
				c.Set(ChannelCommandPayloadKey, payload)
			}
			return next(c)
		}
	}
}

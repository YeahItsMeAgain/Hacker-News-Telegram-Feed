package handlers

import (
	"fmt"
	"hn_feed/bot/utils"
	"hn_feed/db/models"
	db_utils "hn_feed/db/utils"
	"log"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

const (
	channelCommandPayloadKey = "ChannelCommandPayloadKey"
)

func CreateChannelCommandsHandler(handlers map[string]telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		text := ctx.Text()
		input := strings.SplitN(text, " ", 2)
		if len(input) > 1 {
			ctx.Set(channelCommandPayloadKey, input[1])
		}

		for command, handler := range handlers {
			if strings.HasPrefix(text, command) {
				ctx.Delete()
				return handler(ctx)
			}
		}
		return nil
	}
}

func OnChannelRegister(ctx telebot.Context) error {
	chat := ctx.Chat()
	db_utils.UpsertByTgId(
		&models.Channel{
			TgId:  chat.ID,
			Title: chat.Title,
		},
	)

	log.Printf("[*] Registered: <%d - %s - %s>.", chat.ID, chat.Title, chat.Username)
	return utils.SilentlySendAndDelete(
		ctx,
		"🚀 Registered!\n\nUse /help to see the configuration options.",
	)
}

func OnChannelHelp(ctx telebot.Context) error {
	return utils.SilentlySendAndDelete(
		ctx,
		`ℹ️ Available Commands:

		- /feed <topstories\newstories\beststories>
		- /count <1-100>
		- /whitelist <keyword\hostname>
		- /blacklist <keyword\hostname>
		`,
	)
}

func OnChannelConfigureCount(ctx telebot.Context) error {
	payload := ctx.Get(channelCommandPayloadKey)
	if payload == nil || payload == "" {
		return utils.SilentlySendAndDelete(ctx, "❗ Specify the number of top posts you want to see!")
	}

	count, err := strconv.Atoi(payload.(string))
	if err != nil ||
		count < 1 || count > 100 {
		return utils.SilentlySendAndDelete(ctx, "❗ The count should be between 1 and 100!")
	}

	chat := ctx.Chat()
	db_utils.UpsertByTgId(
		&models.Channel{
			TgId:          chat.ID,
			Title:         chat.Title,
			TopPostsCount: count,
		},
	)
	log.Printf("[*] Updated count of: <%d - %s - %s> to %d.", chat.ID, chat.Title, chat.Username, count)
	return utils.SilentlySendAndDelete(
		ctx,
		fmt.Sprintf("🚀 Configured to: %d top posts per hour!", count),
	)
}

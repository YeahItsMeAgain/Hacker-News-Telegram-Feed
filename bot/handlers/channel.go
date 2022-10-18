package handlers

import (
	"fmt"
	"hn_feed/bot/utils"
	"hn_feed/db"
	"hn_feed/db/models"
	db_utils "hn_feed/db/utils"
	"log"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
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
	channel := db_utils.GetOrCreateChannel(chat.ID)
	channel.Title = chat.Title
	db.DB.Save(&channel)

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

func OnChannelConfigureFeedType(ctx telebot.Context) error {
	payload := ctx.Get(channelCommandPayloadKey)
	if payload == nil || payload == "" {
		return utils.SilentlySendAndDelete(ctx, "❗ Specify the feed type: <topstories\\newstories\\beststories>!")
	}

	feedType := payload.(string)
	if !slices.Contains(models.FeedTypes, feedType) {
		return utils.SilentlySendAndDelete(ctx, "❗ Specify the feed type: <topstories\\newstories\\beststories>!")
	}

	chat := ctx.Chat()
	channel := db_utils.GetOrCreateChannel(chat.ID)
	channel.FeedType = feedType
	db.DB.Save(&channel)
	log.Printf("[*] Updated feed type of: <%d - %s - %s> to %s.", chat.ID, chat.Title, chat.Username, payload)
	return utils.SilentlySendAndDelete(
		ctx,
		fmt.Sprintf("🚀 Configured feed type to: %s!", payload),
	)
}

func OnChannelConfigureCount(ctx telebot.Context) error {
	payload := ctx.Get(channelCommandPayloadKey)
	if payload == nil || payload == "" {
		return utils.SilentlySendAndDelete(ctx, "❗ Specify the number of posts you want to see!")
	}

	count, err := strconv.Atoi(payload.(string))
	if err != nil ||
		count < 1 || count > 100 {
		return utils.SilentlySendAndDelete(ctx, "❗ The count should be between 1 and 100!")
	}

	chat := ctx.Chat()
	channel := db_utils.GetOrCreateChannel(chat.ID)
	channel.PostsCount = count
	db.DB.Save(&channel)
	log.Printf("[*] Updated count of: <%d - %s - %s> to %d.", chat.ID, chat.Title, chat.Username, count)
	return utils.SilentlySendAndDelete(
		ctx,
		fmt.Sprintf("🚀 Configured count to: %d posts per hour!", count),
	)
}

package handlers

import (
	"fmt"
	"hn_feed/bot/utils"
	"hn_feed/config"
	"hn_feed/db"
	"log"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/telebot.v3"
)

const (
	channelCommandPayloadKey = "ChannelCommandPayloadKey"
)

func ChannelCommandsHandler(handlers map[string]telebot.HandlerFunc) telebot.HandlerFunc {
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
		return OnChannelHelp(ctx)
	}
}

func OnChannelRegister(ctx telebot.Context) error {
	chat := ctx.Chat()
	channel := db.Channel{TgId: chat.ID}
	channel.FirstOrCreate()
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
		fmt.Sprintf(`ℹ️ Available Commands:

		/info
		/feed <topstories\newstories>
		/count <1-%d>
		/whitelist <:empty:\keyword\hostname>
		/blacklist <:empty:\keyword\hostname>
		`, config.Get().MaxPosts),
	)
}

func OnChannelInfo(ctx telebot.Context) error {
	channel := db.Channel{TgId: ctx.Chat().ID}
	channel.FirstOrCreate()
	return utils.SilentlySendAndDelete(
		ctx,
		fmt.Sprintf(`ℹ️ Current Configuration:

		Feed type: %s.
		Max posts per hour: %d.
		`, channel.FeedType, channel.PostsCount,
		),
	)
}

func OnChannelConfigureFeedType(ctx telebot.Context) error {
	payload := ctx.Get(channelCommandPayloadKey)
	if payload == nil || payload == "" {
		return utils.SilentlySendAndDelete(ctx, "❗ Specify the feed type: <topstories\\newstories>!")
	}

	feedType := payload.(string)
	if !slices.Contains(db.FeedTypes, feedType) {
		return utils.SilentlySendAndDelete(ctx, "❗ Specify the feed type: <topstories\\newstories>!")
	}

	chat := ctx.Chat()
	channel := db.Channel{TgId: chat.ID}
	channel.FirstOrCreate()
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
		count < 1 || count > config.Get().MaxPosts {
		return utils.SilentlySendAndDelete(ctx, fmt.Sprintf("❗ The count should be between 1 and %d!", config.Get().MaxPosts))
	}

	chat := ctx.Chat()
	channel := db.Channel{TgId: chat.ID}
	channel.FirstOrCreate()
	channel.PostsCount = count
	db.DB.Save(&channel)
	log.Printf("[*] Updated count of: <%d - %s - %s> to %d.", chat.ID, chat.Title, chat.Username, count)
	return utils.SilentlySendAndDelete(
		ctx,
		fmt.Sprintf("🚀 Configured count to: %d posts per hour!", count),
	)
}

func onChannelConfigureAssociationList(ctx telebot.Context, association string) error {
	chat := ctx.Chat()
	channel := db.Channel{TgId: chat.ID}
	channel.FirstOrCreate()
	payload := ctx.Get(channelCommandPayloadKey)
	if payload == nil || payload == "" {
		return utils.SilentlySendAndDelete(ctx, channel.GetAssociatedKeywords(association))
	}

	keyword := db.Keyword{Keyword: payload.(string)}
	keyword.FirstOrCreate()
	if len(channel.GetAssociatedKeywordsByKeyword(association, keyword)) > 0 {
		db.DB.Model(&channel).Association(association).Delete(&keyword)
		log.Printf("[*] <%d - %s - %s> removed %s: %s.", chat.ID, chat.Title, chat.Username, association, keyword.Keyword)
		return utils.SilentlySendAndDelete(ctx, fmt.Sprintf("🚀 Removed %s: %s", association, keyword.Keyword))
	}

	db.DB.Model(&channel).Association(association).Append(&keyword)
	log.Printf("[*] <%d - %s - %s> added %s: %s.", chat.ID, chat.Title, chat.Username, association, keyword.Keyword)
	return utils.SilentlySendAndDelete(ctx, fmt.Sprintf("🚀 Added %s: %s", association, keyword.Keyword))
}

func OnChannelConfigureWhitelist(ctx telebot.Context) error {
	return onChannelConfigureAssociationList(ctx, "WhitelistedKeywords")
}

func OnChannelConfigureBlacklist(ctx telebot.Context) error {
	return onChannelConfigureAssociationList(ctx, "BlacklistedKeywords")
}

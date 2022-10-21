package timer

import (
	"fmt"
	"hn_feed/config"
	"hn_feed/db"
	"hn_feed/hn_api"
	"html"
	"log"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"github.com/samber/lo"
	"gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

func ScheduleUpdates(bot *telebot.Bot) {
	ticker := time.NewTicker(time.Duration(config.Config.UpdateIntervalMins) * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				var wg sync.WaitGroup
				for _, feedType := range db.FeedTypes {
					wg.Add(1)
					go func(feedType string) {
						defer wg.Done()
						updateChannels(feedType, bot)
					}(feedType)
				}
				wg.Wait()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func updateChannels(feedType string, bot *telebot.Bot) {
	start := time.Now()
	log.Printf("[*] Updating channels with %s.", feedType)
	posts, err := hn_api.GetNewPosts(feedType)
	if err != nil {
		return
	}

	var dbUpdates sync.WaitGroup
	var channels []db.Channel
	db.DB.Preload(clause.Associations).Find(&channels, "feed_type =?", feedType)
	for i, post := range posts {
		// TODO: filter with blacklist\whitelist.
		channelsToUpdate := lo.Filter(channels, func(channel db.Channel, _ int) bool {
			return i < channel.PostsCount &&
				slices.IndexFunc(channel.Posts, func(channelPost *db.Post) bool { return channelPost.PostId == post.PostId }) == -1
		})
		for _, channel := range channelsToUpdate {
			dbUpdates.Add(1)
			go func(channel db.Channel, post db.Post) {
				defer dbUpdates.Done()
				err := db.DB.Model(&channel).Association("Posts").Append(&post)
				if err != nil {
					log.Printf("[!] Error white appending post to channel, %s", err)
				}
			}(channel, post)
			bot.Send(
				&telebot.User{ID: channel.TgId},
				fmt.Sprintf("<b>%s</b>\n\n%s", post.Title, html.EscapeString(post.Url)),
				telebot.NoPreview,
				telebot.ModeHTML,
			)
		}
	}
	dbUpdates.Wait()
	log.Printf("[*] Finished updates of %s, took: %s.", feedType, time.Since(start))
}

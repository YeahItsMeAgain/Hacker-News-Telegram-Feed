package timer

import (
	"fmt"
	"hn_feed/config"
	"hn_feed/db"
	"hn_feed/hn_api"
	"html"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"github.com/samber/lo"
	"gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

func ScheduleUpdates(bot *telebot.Bot) {
	ticker := time.NewTicker(time.Duration(config.Get().UpdateIntervalMins) * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("[*] Updating channels.")
				start := time.Now()
				wg := &sync.WaitGroup{}
				channelUpdatePool := make(chan struct{}, config.Get().ConcurrentChannelUpdateLimit)
				for _, feedType := range db.FeedTypes {
					wg.Add(1)
					go func(feedType string) {
						defer wg.Done()
						updateChannels(feedType, bot, wg, channelUpdatePool)
					}(feedType)
				}
				wg.Wait()
				log.Printf("[*] Finished updates took: %s.", time.Since(start))
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func updateChannels(feedType string, bot *telebot.Bot, wg *sync.WaitGroup, channelUpdatePool chan struct{}) {
	posts, err := hn_api.GetNewPosts(feedType)
	if err != nil {
		return
	}

	var channels []db.Channel
	db.DB.Preload(clause.Associations).Find(&channels, "feed_type =?", feedType)
	if len(channels) == 0 {
		log.Printf("[*] No channels waiting for %s.", feedType)
		return
	}

	for i, post := range posts {
		channelsToUpdate := lo.Filter(channels, func(channel db.Channel, _ int) bool {
			return shouldUpdateChannel(i, post, channel)
		})
		for _, channel := range channelsToUpdate {
			wg.Add(1)
			channelUpdatePool <- struct{}{}
			go func(channelId uint, post db.Post) {
				defer func() {
					wg.Done()
					<-channelUpdatePool
				}()

				// https://github.com/go-gorm/gorm/issues/5801
				channel := db.Channel{}
				db.DB.First(&channel, "id=?", channelId)
				err := db.DB.Model(&channel).Association("Posts").Append(&post)
				if err != nil {
					log.Printf("[!] Error white appending post to channel, %s", err)
					return
				}
				bot.Send(
					&telebot.User{ID: channel.TgId},
					fmt.Sprintf("<b>%s</b>\n\n%s", post.Title, html.EscapeString(post.Url)),
					telebot.NoPreview,
					telebot.ModeHTML,
				)
			}(channel.ID, post)
		}
	}
}

func shouldUpdateChannel(postCount int, post db.Post, channel db.Channel) bool {
	if postCount >= channel.PostsCount {
		return false
	}

	if slices.IndexFunc(channel.Posts, func(channelPost *db.Post) bool {
		return channelPost.PostId == post.PostId
	}) >= 0 {
		return false
	}

	url, err := url.Parse(post.Url)
	if err != nil {
		log.Printf("[!] Error parsing post %d url: %s", post.PostId, post.Url)
		return false
	}

	channelKeywords := append(strings.Split(post.Title, " "), strings.TrimPrefix(url.Hostname(), "www."))
	if len(channel.BlacklistedKeywords) > 0 &&
		lo.Some(channelKeywords, lo.Map(channel.BlacklistedKeywords, func(keyword *db.Keyword, _ int) string {
			return keyword.Keyword
		})) {
		return false
	}

	if len(channel.WhitelistedKeywords) > 0 {
		return lo.Some(channelKeywords, lo.Map(channel.WhitelistedKeywords, func(keyword *db.Keyword, _ int) string {
			return keyword.Keyword
		}))
	}

	return true
}

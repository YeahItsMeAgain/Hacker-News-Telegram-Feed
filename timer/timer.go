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

	"github.com/samber/lo"
	"gopkg.in/telebot.v3"
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
	db.DB.
		Preload("WhitelistedKeywords").
		Preload("BlacklistedKeywords").
		Find(&channels, "feed_type =?", feedType)
	if len(channels) == 0 {
		log.Printf("[*] No channels waiting for %s.", feedType)
		return
	}

	channelsPostCount := make(map[uint]int, len(channels))
	for _, post := range posts {
		channelsToUpdate := lo.Filter(channels, func(channel db.Channel, _ int) bool {
			return shouldUpdateChannel(post, channel)
		})

		for _, channel := range channelsToUpdate {
			if channelsPostCount[channel.ID] > channel.PostsCount {
				continue
			}
			channelsPostCount[channel.ID]++

			wg.Add(1)
			channelUpdatePool <- struct{}{}
			go func(channel db.Channel, post db.Post) {
				defer func() {
					wg.Done()
					<-channelUpdatePool
				}()

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
			}(channel, post)
		}
	}
}

func shouldUpdateChannel(post db.Post, channel db.Channel) bool {
	if post.Score < channel.MinimumScore {
		return false
	}

	if db.DB.Model(&channel).
		Where("channels_posts.post_id = ?", post.ID).
		Association("Posts").Count() > 0 {
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

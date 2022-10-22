package hn_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"hn_feed/config"
	"hn_feed/db"
	"net/http"
	"sync"
)

const baseUrl = "https://hacker-news.firebaseio.com/v0"

var postsLock sync.RWMutex // GetNewPosts gets called in multiple goroutines, FirstOrCreate is not thread safe: https://github.com/go-gorm/gorm/issues/5803.

func GetNewPosts(feedType string) (map[int]db.Post, error) {
	var postsIds []int
	r, err := http.Get(fmt.Sprintf("%s/%s.json", baseUrl, feedType))
	if err != nil {
		return nil, errors.New("can't get new posts")
	}

	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&postsIds)

	var wg sync.WaitGroup
	posts := make(map[int]db.Post, config.Get().MaxPosts)
	for i := 0; i < config.Get().MaxPosts; i++ {
		wg.Add(1)
		go func(postCount int) {
			defer wg.Done()
			r, err := http.Get(fmt.Sprintf("%s/item/%d.json", baseUrl, postsIds[postCount]))
			if err != nil {
				return
			}

			var post db.Post
			json.NewDecoder(r.Body).Decode(&post)
			if post.Title == "" {
				return
			}

			if post.Url == "" {
				post.Url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", post.PostId)
			}

			postsLock.Lock()
			post.FirstOrCreate()
			posts[postCount] = post
			postsLock.Unlock()
		}(i)
	}
	wg.Wait()
	return posts, nil
}

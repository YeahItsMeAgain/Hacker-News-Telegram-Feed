package config

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type configT struct {
	BotToken                     string
	SqliteDb                     string
	AdminIds                     []int64
	UpdateIntervalMins           int
	MaxPosts                     int
	ConcurrentChannelUpdateLimit int
}

var (
	config     *configT
	configLock = new(sync.RWMutex)
)

func Init() {
	load(true)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP)
	go func() {
		for range signals {
			load(false)
		}
	}()
}

func load(fail bool) {
	log.Println("[*] Loading config.")

	var tmpConfig *configT
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&tmpConfig)
	if err != nil {
		log.Printf("[!] Can't read config.json: %s", err)
		if fail {
			os.Exit(1)
		}
	}
	configLock.Lock()
	defer configLock.Unlock()
	config = tmpConfig
}

func Get() *configT {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

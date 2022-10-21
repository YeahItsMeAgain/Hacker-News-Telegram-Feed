package config

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	BotToken           string
	SqliteDb           string
	AdminIds           []int64
	UpdateIntervalMins int
	MaxPosts           int
}

var Config *config

func Init() {
	file, _ := os.Open("config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Config)
	if err != nil {
		log.Fatal("[!] Can't read config.json", err)
	}
}

package utils

import (
	"log"
	"time"

	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/telebot.v3"
)

func StructsToString[E any](elements []E) string {
	if len(elements) == 0 {
		return "The list is empty."
	}

	var res string
	for _, element := range elements {
		val := reflect.ValueOf(element)

		res += "----------\n"
		for i := 0; i < val.NumField(); i++ {
			if strVal := valToString(val.Field(i)); strVal != "" {
				res += fmt.Sprintf("%s: %s\n", val.Type().Field(i).Name, strVal)
			}
		}
		res += "----------\n"
	}
	return res
}

func valToString(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.String:
		return val.String()
	default:
		return ""
	}
}

func CreateBot(botToken string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func CreateReplyMarkup(rows ...telebot.Row) [][]telebot.ReplyButton {
	replyKeys := make([][]telebot.ReplyButton, 0, len(rows))
	for _, row := range rows {
		keys := make([]telebot.ReplyButton, 0, len(row))
		for _, btn := range row {
			btn := btn.Reply()
			if btn != nil {
				keys = append(keys, *btn)
			}
		}
		replyKeys = append(replyKeys, keys)
	}
	return replyKeys
}

func SilentlySendAndDelete(ctx telebot.Context, msg string) error {
	reply, err := ctx.Bot().Send(ctx.Recipient(), msg, telebot.Silent)
	if err != nil {
		return err
	}
	time.AfterFunc(time.Duration(10)*time.Second, func() {
		ctx.Bot().Delete(reply)
	})
	return nil
}

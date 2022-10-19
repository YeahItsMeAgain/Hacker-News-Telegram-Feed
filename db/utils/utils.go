package utils

import (
	"fmt"
	"hn_feed/db"
	"hn_feed/db/models"
	"reflect"
	"strconv"

	"gorm.io/gorm/clause"
)

func StructsToString[E any](elements []E) string {
	if len(elements) == 0 {
		return "The list is empty."
	}

	var res string
	for _, element := range elements {
		res += fmt.Sprintf(
			"----------\n%s----------\n",
			structToString(reflect.ValueOf(element)),
		)
	}
	return res
}

func structToString(val reflect.Value) string {
	var res string
	for i := 0; i < val.NumField(); i++ {
		if strVal := valToString(val.Field(i)); strVal != "" {
			res += fmt.Sprintf("%s: %s\n", val.Type().Field(i).Name, strVal)
		}
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

func GetOrCreateChannel(TgId int64, preload bool) models.Channel {
	dbChannel := models.Channel{TgId: TgId}
	if preload {
		db.DB.Preload(clause.Associations).FirstOrCreate(&dbChannel, "tg_id = ?", TgId)
	} else {
		db.DB.FirstOrCreate(&dbChannel, "tg_id = ?", TgId)
	}
	return dbChannel
}

func GetOrCreateKeyword(keyword string) models.Keyword {
	dbKeyword := models.Keyword{Keyword: keyword}
	db.DB.FirstOrCreate(&dbKeyword, "keyword = ?", keyword)
	return dbKeyword
}

func GetAssociatedKeywords(channel *models.Channel, association string) []string {
	var dbKeywords []models.Keyword
	db.DB.Model(channel).Association(association).Find(&dbKeywords)
	var keywords []string
	for _, keyword := range dbKeywords {
		keywords = append(keywords, keyword.Keyword)
	}
	return keywords
}

package main

import (
	"encoding/json"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

type Config struct {
	TgToken string `mapstructure:"TG_TOKEN"`
	ChatId  string `mapstructure:"CHAT_ID"`
}
type Response struct {
	Type     string `json:"type"`
	Metadata struct {
		Generated int64  `json:"generated"`
		URL       string `json:"url"`
		Title     string `json:"title"`
		Status    int    `json:"status"`
		API       string `json:"api"`
		Count     int    `json:"count"`
	} `json:"metadata"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			Mag     float64 `json:"mag"`
			Place   string  `json:"place"`
			Time    int64   `json:"time"`
			Updated int64   `json:"updated"`
			Tz      any     `json:"tz"`
			URL     string  `json:"url"`
			Detail  string  `json:"detail"`
			Felt    any     `json:"felt"`
			Cdi     any     `json:"cdi"`
			Mmi     any     `json:"mmi"`
			Alert   any     `json:"alert"`
			Status  string  `json:"status"`
			Tsunami int     `json:"tsunami"`
			Sig     int     `json:"sig"`
			Net     string  `json:"net"`
			Code    string  `json:"code"`
			Ids     string  `json:"ids"`
			Sources string  `json:"sources"`
			Types   string  `json:"types"`
			Nst     int     `json:"nst"`
			Dmin    float64 `json:"dmin"`
			Rms     float64 `json:"rms"`
			Gap     int     `json:"gap"`
			MagType string  `json:"magType"`
			Type    string  `json:"type"`
			Title   string  `json:"title"`
		} `json:"properties"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		ID string `json:"id"`
	} `json:"features"`
	Bbox []float64 `json:"bbox"`
}

type Point struct {
	Latitude  float64
	Longitude float64
	Name      string
}

var config Config

var pointA = Point{43.25, 76.9, "Almaty, Kazakhstan"}

func init() {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName("config.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Panic(err)
	}
	if config.TgToken == "" {
		log.Panic("Empty token!")
	}
}

func sendMessageToChannel(bot *tg.BotAPI, msgText string) error {
	msg := tg.NewMessageToChannel(config.ChatId, msgText)
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func NotifyAboutEQ(bot *tg.BotAPI) error {
	var EqData Response
	var message string
	EqData, err := getEqData(4, 1000)
	if err != nil {
		return err
	}

	for _, feature := range EqData.Features {
		distance := CalculateDistanceBetween(
			pointA.Latitude,
			pointA.Longitude,
			feature.Geometry.Coordinates[0],
			feature.Geometry.Coordinates[1])
		if distance <= 200 {
			message += fmt.Sprintf("Place: %s;\nMag: %.1f;\nTime: %s;\n\n\n", feature.Properties.Place, feature.Properties.Mag, timestampToDate(feature.Properties.Time))
		} else {
			log.Printf("This event was occurred %.2f km from %s in the place %s\n", distance, pointA.Name, feature.Properties.Place)
		}
	}
	if message != "" {
		err = sendMessageToChannel(bot, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func CalculateDistanceBetween(LatitudeA, LongitudeA, LatitudeB, LongitudeB float64) float64 {
	const R = 6371e3
	φ1 := LatitudeA * math.Pi / 180 // φ, λ in radians
	φ2 := LatitudeB * math.Pi / 180
	Δφ := (LatitudeB - LatitudeA) * math.Pi / 180
	Δλ := (LongitudeB - LongitudeA) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := (R * c) / 1000
	return d
}

// Перевод из формата unix в нормальную дату
func timestampToDate(timestamp int64) time.Time {
	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		panic(err)
	}
	// Особенность конвертации Unix в Timestamp
	return time.Unix(timestamp/1e3, 0).In(location)
}

func getEqData(minMagnitude int, minutes time.Duration) (Response, error) {
	const ISO_8601 = "2006-01-02T15:04:05"
	var result Response
	startTime := time.Now().Add(-(time.Minute * minutes)).Format(ISO_8601)
	endTime := time.Now().Format(ISO_8601)

	url := fmt.Sprintf("https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&minmag=%d&starttime=%s&endtime=%s",
		minMagnitude, startTime, endTime)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, nil
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, nil
	}

	return result, nil
}

func DoTaskByTime(f func(), minutes time.Duration) {
	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			f()
		}
	}
}

func main() {
	bot, err := tg.NewBotAPI(config.TgToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Successfully authenticated as %s\n", bot.Self.UserName)

	DoTaskByTime(func() {
		err = NotifyAboutEQ(bot)
		if err != nil {
			log.Panic(err)
		}
	}, 1)
}

// DONE 1 через определенный промежуток времени делать запрос в сервис;
// DONE 2 вытащить данные;
// 3 обработать их и проверить близко ли это к координатам Алматы;
// DONE 4 отправить сообщение в канал, где будет бот и другие участники если будет реально землетрясение
// 5 Выложить программу в открытый доступ в какой-то из серверов

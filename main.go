package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	Latitude     float64
	Longitude    float64
	MaxRadius    float64
	MinMagnitude int
	MinutesCount time.Duration
	Name         string
}

var config Config

var pointA = Point{
	43.25,
	76.9,
	800,
	4,
	1,
	"Almaty, Kazakhstan",
}

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

// Notifying about earthquake events using Telegram
func notifyAboutEQ(bot *tg.BotAPI) error {
	var EqData Response
	var message string
	EqData, err := getEqData(pointA.MinMagnitude, pointA.MinutesCount)
	if err != nil {
		return err
	}
	log.Printf("Recieved from request - %d elements", EqData.Metadata.Count)
	for _, feature := range EqData.Features {
		message += fmt.Sprintf("🌍 Earthquake Alert! 🌍\n\n📍 Location: %s\n📏 Magnitude: %.2f\n🕒 Time: %s\n📏 Depth: %.2f km near %s \n\nStay safe, everyone! 🚨\n\n\n",
			feature.Properties.Place,
			feature.Properties.Mag,
			timestampToDate(feature.Properties.Time),
			calculateDistanceBetween(pointA.Longitude, pointA.Latitude, feature.Geometry.Coordinates[1], feature.Geometry.Coordinates[0]),
			pointA.Name)
	}
	if message != "" {
		err = sendMessageToChannel(bot, message)
		if err != nil {
			return err
		}
	}
	return nil
}

// calculateDistanceBetween Calculate distance between two points using Haversine formula
func calculateDistanceBetween(LatitudeA, LongitudeA, LatitudeB, LongitudeB float64) float64 {
	const R = 6371e3
	phi1 := LatitudeA * math.Pi / 180 // φ, λ in radians
	phi2 := LatitudeB * math.Pi / 180
	deltaPhi := (LatitudeB - LatitudeA) * math.Pi / 180
	deltaLambda := (LongitudeB - LongitudeA) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := (R * c) / 1000
	return d
}

// timestampToDate timeConvert from unix epoch to timestamp
func timestampToDate(timestamp int64) time.Time {
	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		panic(err)
	}
	return time.Unix(timestamp/1e3, 0).In(location)
}

// getEqData Sends request to USGS service and receives data about earthquake events
func getEqData(minMagnitude int, minutes time.Duration) (Response, error) {
	const ISO_8601 = "2006-01-02T15:04:05"
	var result Response
	startTime := time.Now().Add(-(time.Minute * minutes)).Format(ISO_8601)
	endTime := time.Now().Format(ISO_8601)
	url := fmt.Sprintf("https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&minmag=%d&starttime=%s&endtime=%s&latitude=%.2f&longitude=%.2f&maxradiuskm=%.1f",
		minMagnitude, startTime, endTime, pointA.Latitude, pointA.Longitude, pointA.MaxRadius)

	req, err := http.NewRequest("GET", url, nil)
	log.Printf("Sending request to %s", url)
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

// doTaskByTime Simple scheduler implementation, that invokes some function once in defined time measure
func doTaskByTime(f func(), minutes time.Duration) {
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
	// logging requests
	//bot.Debug = true
	log.Printf("Successfully authenticated as %s\n", bot.Self.UserName)

	doTaskByTime(func() {
		err = notifyAboutEQ(bot)
		if err != nil {
			log.Panic(err)
		}
	}, 1)
}

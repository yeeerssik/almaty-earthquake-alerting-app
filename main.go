package main

import (
	"almeqapp/config"
	"almeqapp/internal/models"
	"almeqapp/internal/service"
	"almeqapp/utils/helper"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Panic(err)
	}

	tgService, err := service.New(c)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Successfully authenticated as %s\n", tgService.Bot.Self.UserName)

	doTaskByTime(func() {
		err = notifyAboutEQ(c, tgService.Bot)
		if err != nil {
			log.Panic(err)
		}
	}, 1)
}

func sendMessageToChannel(conf *config.Config, bot *tg.BotAPI, msgText string) error {
	msg := tg.NewMessageToChannel(conf.ChatId, msgText)
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// Notifying about earthquake events using Telegram
func notifyAboutEQ(conf *config.Config, bot *tg.BotAPI) error {
	var EqData models.Response
	var message string
	EqData, err := getEqData(conf)
	if err != nil {
		return err
	}
	log.Printf("Recieved from request - %d elements", EqData.Metadata.Count)
	for _, feature := range EqData.Features {
		message += fmt.Sprintf("🌍 Earthquake Alert! 🌍\n\n📍 Location: %s\n📏 Magnitude: %.2f\n🕒 Time: %s\n📏 Depth: %.2f km near %s \n\nStay safe, everyone! 🚨\n\n\n",
			feature.Properties.Place,
			feature.Properties.Mag,
			helper.TimestampToDate(feature.Properties.Time),
			calculateDistanceBetween(conf.TargetPlaceLongitude, conf.TargetPlaceLatitude, feature.Geometry.Coordinates[1], feature.Geometry.Coordinates[0]),
			conf.TargetPlaceName)
	}
	if message != "" {
		err = sendMessageToChannel(conf, bot, message)
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

// getEqData Sends request to USGS service and receives data about earthquake events
func getEqData(conf *config.Config) (result models.Response, err error) {
	const ISO_8601 = "2006-01-02T15:04:05"

	startTime := time.Now().Add(-(time.Minute * time.Duration(conf.MinutesCount))).Format(ISO_8601)
	endTime := time.Now().Format(ISO_8601)
	url := fmt.Sprintf("https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&minmag=%d&starttime=%s&endtime=%s&latitude=%.2f&longitude=%.2f&maxradiuskm=%.1f",
		conf.MinMagnitude, startTime, endTime, conf.TargetPlaceLatitude, conf.TargetPlaceLongitude, conf.MaxRadius)

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

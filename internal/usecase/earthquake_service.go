package usecase

import (
	"almeqapp/config"
	"almeqapp/internal/client"
	"almeqapp/internal/models"
	"almeqapp/internal/service"
	"almeqapp/utils/helper"
	"fmt"
	"log"
)

type TelegramService interface {
	SendMessage(message string) (err error)
}

type EarthquakeService struct {
	conf      *config.Config
	apiClient *client.Client
	tgService *service.TelegramService
}

func New(conf *config.Config, apiClient *client.Client, tgService *service.TelegramService) *EarthquakeService {
	return &EarthquakeService{conf, apiClient, tgService}
}

func (e *EarthquakeService) NotifyAboutEQ() (err error) {
	var EqData models.Response
	var message string
	EqData, err = e.apiClient.GetData()
	if err != nil {
		return err
	}
	log.Printf("Recieved from request - %d elements", EqData.Metadata.Count)
	for _, feature := range EqData.Features {
		message += fmt.Sprintf("🌍 Earthquake Alert! 🌍\n\n📍 Location: %s\n📏 Magnitude: %.2f\n🕒 Time: %s\n📏 Depth: %.2f km near %s \n\nStay safe, everyone! 🚨\n\n\n",
			feature.Properties.Place,
			feature.Properties.Mag,
			helper.TimestampToDate(feature.Properties.Time),
			helper.CalculateDistanceBetween(
				e.conf.TargetPlaceLongitude,
				e.conf.TargetPlaceLatitude,
				feature.Geometry.Coordinates[1],
				feature.Geometry.Coordinates[0]),
			e.conf.TargetPlaceName,
		)
	}
	if message != "" {
		err = e.tgService.SendMessage(message)
		if err != nil {
			return err
		}
	}
	return
}

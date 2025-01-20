package client

import (
	"almeqapp/config"
	"almeqapp/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	conf *config.Config
}

const (
	ISO_8601 = "2006-01-02T15:04:05"
)

func New(conf *config.Config) *Client {
	return &Client{conf}
}

func (c *Client) GetData() (result models.Response, err error) {
	startTime := time.Now().Add(-(time.Minute * time.Duration(c.conf.MinutesCount))).Format(ISO_8601)
	endTime := time.Now().Format(ISO_8601)
	url := fmt.Sprintf("https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&minmag=%d&starttime=%s&endtime=%s&latitude=%.2f&longitude=%.2f&maxradiuskm=%.1f",
		c.conf.MinMagnitude,
		startTime,
		endTime,
		c.conf.TargetPlaceLatitude,
		c.conf.TargetPlaceLongitude,
		c.conf.MaxRadius,
	)

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

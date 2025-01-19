package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TgToken              string  `mapstructure:"TG_TOKEN"`
	ChatId               string  `mapstructure:"CHAT_ID"`
	TargetPlaceLatitude  float64 `mapstructure:"TARGET_PLACE_LATITUDE"`
	TargetPlaceLongitude float64 `mapstructure:"TARGET_PLACE_LONGITUDE"`
	TargetPlaceName      string  `mapstructure:"TARGET_PLACE_NAME"`
	MaxRadius            float64 `mapstructure:"MAX_RADIUS"`
	MinMagnitude         int     `mapstructure:"MIN_MAGNITUDE"`
	MinutesCount         int     `mapstructure:"MINUTES_COUNT"`
}

func New() (c *Config, err error) {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName("config.env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

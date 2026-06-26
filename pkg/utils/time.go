package utils

import (
	"fmt"
	"os"
	"time"
)

func GetCurrentAppTime() (time.Time, error) {
	timeZone := os.Getenv("APP_TIMEZONE")
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load location: %w", err)
	}

	return time.Now().In(location), nil
}

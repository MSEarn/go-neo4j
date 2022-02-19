package util

import (
	"time"

	"go.uber.org/zap"
)

var (
	Location *time.Location
)

func init() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		zap.L().Panic(`Cannot load timezone "Asia/Bangkok"`)
	}
	Location = loc
}

func Now() time.Time {
	return time.Now().In(Location)
}

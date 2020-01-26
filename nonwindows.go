// +build !windows

package main

import (
	"log"
	"time"
)

func setSystemTime(t time.Time) error {
	// NOP
	return nil
}

func systemTray() error {
	// satisfy 'unused' linter
	log.Printf(
		"%s %s %s %s %s %s %s %s",
		gpsdata.formatStatus(),
		gpsdata.formatGridsquare(),
		gpsdata.formatLatitude(),
		gpsdata.formatLongitude(),
		gpsdata.formatTime(),
		gpsdata.formatFixQuality(),
		gpsdata.formatNumSatellites(),
		gpsdata.formatHDOP(),
	)

	// NOP
	return nil
}

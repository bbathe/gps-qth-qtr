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
		"%v %v %v %v %v %v",
		gpsdata.getStatus(),
		gpsdata.getLocation(),
		gpsdata.getTime(),
		gpsdata.getFixQuality(),
		gpsdata.getNumSatellites(),
		gpsdata.getHDOP(),
	)

	// NOP
	return nil
}

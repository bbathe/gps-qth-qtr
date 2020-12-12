package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// our type
// gpsData is a structure to control concurrent access to the data from the gps device.
type gpsData struct {
	s   string
	tm  time.Time
	lat float64
	lon float64
	loc string
	q   string
	n   int
	h   float64
	mu  sync.RWMutex
}

// newGPSData is for initializing a new gpsData.
func newGPSData() *gpsData {
	return &gpsData{
		lat: -91.0,
		lon: -181.0,
		n:   -1,
		h:   -1.0,
	}
}

// copy duplicate values from new.
func (g *gpsData) copy(new *gpsData) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.s = new.s
	g.tm = new.tm
	g.loc = new.loc
	g.lat = new.lat
	g.lon = new.lon
	g.q = new.q
	g.n = new.n
	g.h = new.h
}

// getStatus returns the status.
func (g *gpsData) getStatus() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.s
}

//setStatus sets the status.
func (g *gpsData) setStatus(s string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.s = s
}

// formatStatus returns a string representation of status to show user.
func (g *gpsData) formatStatus() string {
	s := g.getStatus()

	if s == "" {
		return "OK"
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

// getTime returns the time.
func (g *gpsData) getTime() time.Time {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.tm
}

// setTime sets the time.
func (g *gpsData) setTime(t time.Time) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.tm = t
}

// formatTime returns a string representation of time to show user.
func (g *gpsData) formatTime() string {
	tm := g.getTime()

	if tm != (time.Time{}) {
		return tm.Format("02-Jan-2006 15:04:05 UTC")
	}
	return ""
}

// getLocation returns the gridsquare.
func (g *gpsData) getGridsquare() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.loc
}

// setLocation sets the gridsquare.
func (g *gpsData) setGridsquare(l string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.loc = l
}

// formatGridsquare returns a string representation of location as a gridsquare to show user.
func (g *gpsData) formatGridsquare() string {
	return g.getGridsquare()
}

// getLatitude returns the latitude.
func (g *gpsData) getLatitude() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.lat
}

// setLatitude sets the latitude.
func (g *gpsData) setLatitude(l float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.lat = l
}

// formatLatitude returns a string representation of the latitude to show user.
func (g *gpsData) formatLatitude() string {
	lat := g.getLatitude()

	if lat >= -90.0 && lat <= 90.0 {
		return strconv.FormatFloat(lat, 'f', -1, 64)
	}
	return ""
}

// getLongitude returns the longitude.
func (g *gpsData) getLongitude() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.lon
}

// setLongitude sets the longitude.
func (g *gpsData) setLongitude(l float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.lon = l
}

// formatLongitude returns a string representation of the longitude to show user.
func (g *gpsData) formatLongitude() string {
	lon := g.getLongitude()

	if lon >= -180.0 && lon <= 180.0 {
		return strconv.FormatFloat(lon, 'f', -1, 64)
	}
	return ""
}

// getFixQuality returns the fix quality.
func (g *gpsData) getFixQuality() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.q
}

// setFixQuality sets the fix quality.
func (g *gpsData) setFixQuality(q string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.q = q
}

// formatFixQuality returns a string representation of the fixquality to show user.
func (g *gpsData) formatFixQuality() string {
	return g.getFixQuality()
}

// getNumSatellites returns the number of satellites.
func (g *gpsData) getNumSatellites() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.n
}

// setNumSatellites sets the number of satellites.
func (g *gpsData) setNumSatellites(n int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.n = n
}

// formatNumSatellites returns a string representation of the number of satellites to show user.
func (g *gpsData) formatNumSatellites() string {
	n := g.getNumSatellites()

	if n > -1 {
		return fmt.Sprintf("%d", n)
	}
	return ""
}

// getHDOP returns the horizontal dilution of precision.
func (g *gpsData) getHDOP() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.h
}

// setHDOP sets the horizontal dilution of precision.
func (g *gpsData) setHDOP(h float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.h = h
}

// formatHDOP returns a string representation of the horizontal dilution of precision to show user.
func (g *gpsData) formatHDOP() string {
	h := g.getHDOP()

	if h > -1 {
		return strconv.FormatFloat(h, 'f', -1, 64)
	}
	return ""
}

package main

import (
	"sync"
	"time"
)

// our type
// gpsData is a structure to control concurrent access to the data from the gps device
type gpsData struct {
	s   string
	tm  time.Time
	loc string
	q   string
	n   int
	h   float64
	mu  sync.RWMutex
}

// newGPSData is for initializing a new gpsData
func newGPSData() *gpsData {
	return &gpsData{
		n: -1,
		h: -1.0,
	}
}

// copy duplicate values from new
func (g *gpsData) copy(new *gpsData) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.s = new.s
	g.tm = new.tm
	g.loc = new.loc
	g.q = new.q
	g.n = new.n
	g.h = new.h
}

// getStatus returns the status
func (g *gpsData) getStatus() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.s
}

//setStatus sets the status
func (g *gpsData) setStatus(s string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.s = s
}

// getTime returns the time
func (g *gpsData) getTime() time.Time {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.tm
}

// setTime sets the time
func (g *gpsData) setTime(t time.Time) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.tm = t
}

// getLocation returns the location
func (g *gpsData) getLocation() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.loc
}

// setLocation sets the location
func (g *gpsData) setLocation(l string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.loc = l
}

// getFixQuality returns the fix quality
func (g *gpsData) getFixQuality() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.q
}

// setFixQuality sets the fix quality
func (g *gpsData) setFixQuality(q string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.q = q
}

// getNumSatellites returns the number of satellites
func (g *gpsData) getNumSatellites() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.n
}

// setNumSatellites sets the number of satellites
func (g *gpsData) setNumSatellites(n int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.n = n
}

// getHDOP returns the horizontal dilution of precision
func (g *gpsData) getHDOP() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.h
}

// setHDOP sets the horizontal dilution of precision
func (g *gpsData) setHDOP(h float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.h = h
}

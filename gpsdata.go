package main

import (
	"sync"
	"time"
)

// gpsData is a structure to control concurrent access to the data from the gps device
type gpsData struct {
	s   string
	tm  time.Time
	loc string
	q   string
	n   int
	mu  sync.RWMutex
}

func (g *gpsData) getStatus() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.s
}

func (g *gpsData) setStatus(s string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.s = s
}

func (g *gpsData) getTime() time.Time {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.tm
}

func (g *gpsData) setTime(t time.Time) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.tm = t
}

func (g *gpsData) getLocation() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.loc
}

func (g *gpsData) setLocation(l string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.loc = l
}

func (g *gpsData) getFixQuality() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.q
}

func (g *gpsData) setFixQuality(q string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.q = q
}

func (g *gpsData) getNumSatellites() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.n
}

func (g *gpsData) setNumSatellites(n int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.n = n
}

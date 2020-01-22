package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

// parseRMCTime extracts the time from an **RMC line
func parseRMCTime(fields []string) (time.Time, error) {
	d := fields[9]

	year, err := strconv.Atoi(d[4:6])
	if err != nil {
		log.Printf("%+v", err)
		return time.Time{}, err
	}
	year += 2000

	mon, err := strconv.Atoi(d[2:4])
	if err != nil {
		log.Printf("%+v", err)
		return time.Time{}, err
	}

	day, err := strconv.Atoi(d[0:2])
	if err != nil {
		log.Printf("%+v", err)
		return time.Time{}, err
	}

	t := strings.Split(fields[1], ".")

	hour, err := strconv.Atoi(t[0][0:2])
	if err != nil {
		log.Printf("%+v", err)
		return time.Time{}, err
	}

	min, err := strconv.Atoi(t[0][2:4])
	if err != nil {
		log.Printf("%+v", err)
		return time.Time{}, err
	}

	sec, err := strconv.Atoi(t[0][4:6])
	if err != nil {
		log.Printf("%+v", err)
		return time.Time{}, err
	}

	return time.Date(year, time.Month(mon), day, hour, min, sec, 0, time.UTC), nil
}

// parseDegMinToFloat parses NMEA format of (d)ddmm.mmmm to decimal degrees
func parseDegMinToFloat(dm string) (float64, error) {
	// handle latitude (0-90) or longitude (0-180) degrees
	t := strings.Split(dm, ".")
	p := 2
	if len(t[0]) == 5 {
		p = 3
	}

	degrees, err := strconv.ParseFloat(dm[0:p], 64)
	if err != nil {
		log.Printf("%+v", err)
		return 0, err
	}

	minutes, err := strconv.ParseFloat(dm[p:], 64)
	if err != nil {
		log.Printf("%+v", err)
		return 0, err
	}

	// make sure degrees are within bounds
	if (p == 2 && math.Abs(degrees) > 90) || (p == 3 && math.Abs(degrees) > 180) {
		err := fmt.Errorf("invalid location")
		log.Printf("%+v", err)
		return 0.0, err
	}

	return degrees + minutes/60, nil
}

// latLonToGridsquare converts decimal latitude & longitude to a maidenhead gridsquare
func latLonToGridsquare(lat, lon float64) (string, error) {
	upper := []rune("ABCDEFGHIJKLMNOPQRSTUVWX")
	lower := []rune("abcdefghijklmnopqrstuvwx")

	if (math.Abs(lat) > 90) || (math.Abs(lon) > 180) {
		err := fmt.Errorf("invalid location")
		log.Printf("%+v", err)
		return "", err
	}

	adjLat := lat + 90
	adjLon := lon + 180
	return fmt.Sprintf("%c%c%d%d%c%c",
		upper[int(adjLon/20)],
		upper[int(adjLat/10)],
		int(math.Mod((adjLon/2), 10)),
		int(math.Mod(adjLat, 10)),
		lower[int((adjLon-float64(2*int(adjLon/2)))*60/5)],
		lower[int((adjLat-float64(int(adjLat)))*60/2.5)],
	), nil
}

// parseRMCLocation extracts the location (as maidenhead gridsquare) from an **RMC line
func parseRMCLocation(fields []string) (string, error) {
	lat, err := parseDegMinToFloat(fields[3])
	if err != nil {
		log.Printf("%+v", err)
		return "", err
	}
	if fields[4] == "S" {
		lat = -lat
	}

	lon, err := parseDegMinToFloat(fields[5])
	if err != nil {
		log.Printf("%+v", err)
		return "", err
	}
	if fields[6] == "W" {
		lon = -lon
	}

	gridsquare, err := latLonToGridsquare(lat, lon)
	if err != nil {
		log.Printf("%+v", err)
		return "", err
	}

	return gridsquare, nil
}

// parseRMC extracts the time and location (as maidenhead gridsquare) from an **RMC line
func parseRMC(s string) (time.Time, string, error) {
	// parse comma delimted records to fields
	r := csv.NewReader(strings.NewReader(s))
	fields, err := r.Read()
	if err != nil {
		log.Printf("%+v", err)
		return time.Time{}, "", err
	}

	// need at least 11 fields to get time, location, and checksum
	if len(fields) < 11 {
		err := fmt.Errorf("invalid RMC line")
		log.Printf("%+v", err)
		return time.Time{}, "", err
	}

	// validate checksum
	checksum := 0
	strchk := strings.Split(s, "*")
	if len(strchk) < 2 {
		err := fmt.Errorf("missing checksum")
		log.Printf("%+v", err)
		return time.Time{}, "", err
	}

	for _, c := range strchk[0] {
		checksum ^= int(c)
	}
	if fmt.Sprintf("%X", checksum) != strings.TrimSpace(strchk[1]) {
		err := fmt.Errorf("RMC line bad checksum")
		log.Printf("%+v", err)
		return time.Time{}, "", err
	}

	// status is valid?
	if fields[2] != "A" {
		err := fmt.Errorf("receiver not in valid state")
		log.Printf("%+v", err)
		return time.Time{}, "", err
	}

	// get time
	t, err := parseRMCTime(fields)
	if err != nil {
		log.Printf("%+v", err)
		return time.Time{}, "", err
	}

	// get maidenhead gridsquare
	l, err := parseRMCLocation(fields)
	if err != nil {
		log.Printf("%+v", err)
		return time.Time{}, "", err
	}

	return t, l, nil
}

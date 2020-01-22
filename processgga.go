package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// parseGGA extracts the fix quality and number of satellites being tracked from a **GGA line
func parseGGA(s string) (string, int, error) {
	// parse comma delimted records to fields
	r := csv.NewReader(strings.NewReader(s))
	fields, err := r.Read()
	if err != nil {
		log.Printf("%+v", err)
		return "", 0, err
	}

	// need at least 8 fields to get fix quality and number of satellites being tracked
	if len(fields) < 8 {
		err := fmt.Errorf("invalid GGA line")
		log.Printf("%+v", err)
		return "", 0, err
	}

	// validate checksum
	checksum := 0
	strchk := strings.Split(s, "*")
	if len(strchk) < 2 {
		err := fmt.Errorf("missing checksum")
		log.Printf("%+v", err)
		return "", 0, err
	}

	for _, c := range strchk[0] {
		checksum ^= int(c)
	}
	if fmt.Sprintf("%X", checksum) != strings.TrimSpace(strchk[1]) {
		err := fmt.Errorf("GGA line bad checksum")
		log.Printf("%+v", err)
		return "", 0, err
	}

	// get fix quality
	q, err := strconv.Atoi(fields[6])
	if err != nil {
		log.Printf("%+v", err)
		return "", 0, err
	}
	qs := "invalid"
	switch q {
	case 1:
		qs = "GPS fix (SPS)"
	case 2:
		qs = "DGPS fix"
	case 3:
		qs = "PPS fix"
	case 4:
		qs = "Real Time Kinematic"
	case 5:
		qs = "Float RTK"
	case 6:
		qs = "estimated (dead reckoning)"
	case 7:
		qs = "Manual input mode"
	case 8:
		qs = "Simulation mode"
	}

	// get number of satellites being tracked
	n, err := strconv.Atoi(fields[7])
	if err != nil {
		log.Printf("%+v", err)
		return "", 0, err
	}

	return qs, n, nil
}

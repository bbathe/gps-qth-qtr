package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
	"gopkg.in/yaml.v2"
)

// configuration holds the application configuration
type configuration struct {
	GPSDevice struct {
		Port string
		Baud int
	}
}

var (
	// application configuration
	config configuration

	// last data from gps device
	gpsdata gpsData
)

// readLineFromPort reads bytes from port and accumulates string until delim is met
// does not return delim in string
func readLineFromPort(p *serial.Port, delim byte) (string, error) {
	var s string
	buf := []byte{0}

	for {
		n, err := p.Read(buf)
		if err != nil {
			log.Printf("%+v", err)
			return "", err
		}

		if n > 0 {
			if buf[0] == delim {
				return s, nil
			}
			s += string(buf[:n])
		}
	}
}

// gatherGpsData reads from gps port until a **RMC line is successfully processed
// system time is updated
func gatherGpsData() {
	config := &serial.Config{
		Name: config.GPSDevice.Port,
		Baud: config.GPSDevice.Baud,
	}

	p, err := serial.OpenPort(config)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	defer p.Close()

	var gotrmc bool
	var gotgga bool

	for {
		s, err := readLineFromPort(p, '$')
		if err != nil {
			log.Printf("%+v", err)
			continue
		}

		if len(s) > 5 {
			sentence := s[2:5]

			switch sentence {
			case "RMC":
				t, l, err := parseRMC(s)
				if err != nil {
					log.Printf("%+v", err)
					continue
				}

				// keep last known values
				gpsdata.setTime(t)
				gpsdata.setLocation(l)

				// update system time
				err = setSystemTime(t)
				if err != nil {
					gpsdata.setStatus(fmt.Sprintf("error setting time: %+v", err))
				}

				gotrmc = true
			case "GGA":
				q, n, err := parseGGA(s)
				if err != nil {
					log.Printf("%+v", err)
					continue
				}

				// keep last known values
				gpsdata.setFixQuality(q)
				gpsdata.setNumSatellites(n)

				gotgga = true
			}
		}

		if gotrmc && gotgga {
			// done for now
			return
		}
	}
}

func main() {
	// show file & location, date & time
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// log to file
	f, err := os.OpenFile("gps-qth-qtr.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// command line app
	flag.Usage = func() {
		log.Printf("invalid command line")
	}

	// get config file from command line
	configFile := flag.String("config", "", "Application configuration file")
	flag.Parse()

	if len(*configFile) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// read config
	// #nosec G304
	bytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// prime data
	q1 := scheduleOnce(func() {
		gatherGpsData()
	}, 0)
	defer close(q1)

	// and create recurring task
	q2 := scheduleRecurring(func() {
		gatherGpsData()
	}, 300*time.Second)
	defer close(q2)

	// returns on exit
	err = systemTray()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

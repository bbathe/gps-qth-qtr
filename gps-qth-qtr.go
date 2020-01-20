package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"
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

// gpsData is a structure to control concurrent access to the data from the gps device
type gpsData struct {
	tm  time.Time
	loc string
	mu  sync.RWMutex
}

var (
	// application configuration
	config configuration

	// last data from gps device
	gpsdata gpsData
)

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

	for {
		s, err := readLineFromPort(p, '$')
		if err != nil {
			log.Printf("%+v", err)
			continue
		}

		if len(s) > 5 {
			// **RMC (GNRNC, GPRNC)
			if s[2:5] == "RMC" {
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
					log.Printf("%+v", err)
					continue
				}
			}
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

	// concurrent
	go gatherGpsData()

	// returns on exit
	systemTray()
}

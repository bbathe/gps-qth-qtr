package main

import (
	"flag"
	"fmt"
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
	}

	for {
		s, err := readLineFromPort(p, '$')
		if err != nil {
			log.Printf("%+v", err)
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

func testOutput() {
	for {
		time.Sleep(time.Second * 3)

		t := gpsdata.getTime()
		l := gpsdata.getLocation()
		fmt.Printf("%+v | %+v\n", t, l)
	}
}

func main() {
	// show file & location, date & time
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// command line app
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage of %s\n", os.Args[0])
		flag.PrintDefaults()
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

	//	concurrent processes
	var wg sync.WaitGroup

	funcs := []func(){
		gatherGpsData,
		testOutput,
	}

	for _, f := range funcs {
		wg.Add(1)

		go func(fn func()) {
			defer wg.Done()
			fn()
		}(f)
	}
	wg.Wait()
}

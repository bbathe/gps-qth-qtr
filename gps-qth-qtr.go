package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/tarm/serial"
	"gopkg.in/yaml.v2"
)

// configuration holds the application configuration
type configuration struct {
	GPSDevice struct {
		Port     string
		Baud     int
		PollRate time.Duration
	}
}

var (
	// application configuration
	config configuration

	// last data from gps device
	gpsdata = newGPSData()

	// prevent concurrent processing of gps data
	nbmGatherGpsData = NewNonBlockingMutex()
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

// gatherGpsData reads from gps port until **RMC & **GGA lines are successfully processed
// system time is updated as long as the quality of the gps signal is good enough (HDOP < 5)
func gatherGpsData() {
	if nbmGatherGpsData.Lock() {
		defer nbmGatherGpsData.Unlock()

		// we want the update to the global gpsdata to be atomic and only if there was no errors in gathering data
		var err error
		newgpsdata := newGPSData()
		defer func() {
			if err != nil {
				// reset newgpsdata
				newgpsdata = newGPSData()

				// set message to error string
				newgpsdata.setStatus(err.Error())
			}
			// copy over new values
			gpsdata.copy(newgpsdata)
		}()

		config := &serial.Config{
			Name: config.GPSDevice.Port,
			Baud: config.GPSDevice.Baud,
		}

		var p *serial.Port
		p, err = serial.OpenPort(config)
		if err != nil {
			log.Printf("%+v", err)
			return
		}
		defer p.Close()

		var gotrmc bool
		var gotgga bool

		for {
			var s string
			s, err = readLineFromPort(p, '$')
			if err != nil {
				log.Printf("%+v", err)
				return
			}

			if len(s) > 5 {
				sentence := s[2:5]

				switch sentence {
				case "RMC":
					var t time.Time
					var l string
					var lat, lon float64
					t, l, lat, lon, err = parseRMC(s)
					if err != nil {
						log.Printf("%+v|%+s", err, s)
						return
					}

					// keep values
					newgpsdata.setTime(t)
					newgpsdata.setGridsquare(l)
					newgpsdata.setLatitude(lat)
					newgpsdata.setLongitude(lon)

					gotrmc = true
				case "GGA":
					var q string
					var n int
					var h float64
					q, n, h, err = parseGGA(s)
					if err != nil {
						log.Printf("%+v|%+s", err, s)
						return
					}

					// keep values
					newgpsdata.setFixQuality(q)
					newgpsdata.setNumSatellites(n)
					newgpsdata.setHDOP(h)

					gotgga = true
				}
			}

			// if we were able to capture all the data we need
			if gotrmc && gotgga {
				// and gps signal good enough
				if newgpsdata.getHDOP() < 5 {
					// update system time
					err = setSystemTime(newgpsdata.getTime())
					if err != nil {
						log.Printf("%+v", err)
					}
				}
				return
			}
		}
	}
}

func main() {
	// show file & location, date & time
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// log & config files are in the same directory as the executable with the same base name
	fn, err := os.Executable()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	basefn := strings.TrimSuffix(fn, path.Ext(fn))

	// log to file
	f, err := os.OpenFile(basefn+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// read config
	// #nosec G304
	bytes, err := ioutil.ReadFile(basefn + ".yaml")
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
	}, config.GPSDevice.PollRate*time.Second)
	defer close(q2)

	// returns on exit
	err = systemTray()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

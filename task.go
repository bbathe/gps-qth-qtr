package main

import "time"

// scheduleRecurring creates a recurring task, returns quit channel
func scheduleRecurring(what func(), delay time.Duration) chan bool {
	ticker := time.NewTicker(delay)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				what()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

// scheduleOnce creates a task that fires once, returns quit channel
func scheduleOnce(what func(), delay time.Duration) chan bool {
	if delay == 0 {
		delay = 1
	}
	ticker := time.NewTicker(delay)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				ticker.Stop()
				what()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

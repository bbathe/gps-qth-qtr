package main

import (
	"log"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procSetSystemTime *windows.Proc
)

// setSystemTime calls the windows SetSystemTime API
func setSystemTime(t time.Time) error {
	// convert time types
	systime := windows.Systemtime{
		Year:   uint16(t.Year()),
		Month:  uint16(t.Month()),
		Day:    uint16(t.Day()),
		Hour:   uint16(t.Hour()),
		Minute: uint16(t.Minute()),
		Second: uint16(t.Second()),
	}

	// make call to windows api
	r1, _, err := procSetSystemTime.Call(uintptr(unsafe.Pointer(&systime)))
	if r1 == 0 {
		log.Printf("%+v", err)
		return err
	}
	return nil
}

func init() {
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		log.Fatal(err)
	}

	procSetSystemTime, err = dll.FindProc("SetSystemTime")
	if err != nil {
		log.Fatal(err)
	}
}

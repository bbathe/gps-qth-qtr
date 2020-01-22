// +build windows

package main

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/lxn/walk"
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

func systemTray() error {
	// our window
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// create the notify icon and make sure we clean it up on exit
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	defer func() {
		err := ni.Dispose()
		if err != nil {
			log.Printf("%+v", err)
		}
	}()

	// set the icon and a tool tip text
	icon, err := walk.Resources.Icon("3")
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	err = ni.SetIcon(icon)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	err = ni.SetToolTip("gps-qth-qtr")
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// gridsquare action in context menu
	gridsquareAction := walk.NewAction()
	err = gridsquareAction.SetText("Gridsquare")
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	gridsquareAction.Triggered().Attach(func() {
		err := walk.Clipboard().SetText(gpsdata.getLocation())
		if err != nil {
			log.Printf("%+v", err)
		}
	})
	err = ni.ContextMenu().Actions().Add(gridsquareAction)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// status action in context menu
	statusAction := walk.NewAction()
	err = statusAction.SetText("Status")
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	statusAction.Triggered().Attach(func() {
		status := fmt.Sprintf(
			"Message: %s\nGridsquare: %s\nTime: %v\nFix Quality: %s\nSatellites: %d",
			gpsdata.getStatus(),
			gpsdata.getLocation(),
			gpsdata.getTime(),
			gpsdata.getFixQuality(),
			gpsdata.getNumSatellites(),
		)

		walk.MsgBox(mw, "Status", status, walk.MsgBoxIconInformation)
	})
	err = ni.ContextMenu().Actions().Add(statusAction)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// exit action in context menu
	exitAction := walk.NewAction()
	err = exitAction.SetText("Exit")
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})
	err = ni.ContextMenu().Actions().Add(exitAction)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// make notify icon visible
	err = ni.SetVisible(true)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// message loop
	mw.Run()

	return nil
}

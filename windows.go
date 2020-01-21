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

func systemTray() {
	// our window
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	// create the notify icon and make sure we clean it up on exit
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := ni.Dispose()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// set the icon and a tool tip text
	icon, err := walk.Resources.Icon("3")
	if err != nil {
		log.Fatal(err)
	}
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("gps-qth-qtr"); err != nil {
		log.Fatal(err)
	}

	// gridsquare action in context menu
	gridsquareAction := walk.NewAction()
	if err := gridsquareAction.SetText("Gridsquare"); err != nil {
		log.Fatal(err)
	}
	gridsquareAction.Triggered().Attach(func() {
		err := walk.Clipboard().SetText(gpsdata.getLocation())
		if err != nil {
			log.Printf("%+v", err)
		}
	})
	if err := ni.ContextMenu().Actions().Add(gridsquareAction); err != nil {
		log.Fatal(err)
	}

	// status action in context menu
	statusAction := walk.NewAction()
	if err := statusAction.SetText("Status"); err != nil {
		log.Fatal(err)
	}
	statusAction.Triggered().Attach(func() {
		status := fmt.Sprintf("Gridsquare: %s\nTime: %v", gpsdata.getLocation(), gpsdata.getTime())
		walk.MsgBox(mw, "Status", status, walk.MsgBoxIconInformation)
	})
	if err := ni.ContextMenu().Actions().Add(statusAction); err != nil {
		log.Fatal(err)
	}

	// exit action in context menu
	exitAction := walk.NewAction()
	if err := exitAction.SetText("Exit"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	// make notify icon visible
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// message loop
	mw.Run()
}

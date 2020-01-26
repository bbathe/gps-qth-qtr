// +build windows

package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"unsafe"

	"github.com/lxn/walk"
	declarative "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var (
	procSetSystemTime *windows.Proc
	appIcon           *walk.Icon
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

func RunStatusWindow() error {
	var statusWindow *walk.MainWindow

	mw := declarative.MainWindow{
		AssignTo: &statusWindow,
		Name:     "statusmw",
		Title:    "Status Data",
		Icon:     appIcon,
		Size:     declarative.Size{Width: 300, Height: 200},
		MinSize:  declarative.Size{Width: 300, Height: 200},
		MaxSize:  declarative.Size{Width: 300, Height: 200},
		Layout:   declarative.VBox{MarginsZero: true},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout:        declarative.Grid{Rows: 2},
				StretchFactor: 4,
				Children: []declarative.Widget{
					declarative.TableView{
						Name:                "statustv",
						ColumnsOrderable:    true,
						AlternatingRowBG:    true,
						HeaderHidden:        true,
						LastColumnStretched: true,
						Columns: []declarative.TableViewColumn{
							{Name: "Index", Hidden: true},
							{Name: "Name"},
							{Name: "Value"},
						},
						Model: NewTableDataModel(),
					},
					declarative.PushButton{
						Text: "OK",
						OnClicked: func() {
							statusWindow.Close()
						},
					},
				},
			},
		},
	}

	err := mw.Create()
	if err != nil {
		return err
	}

	// disable maximize, minimize, and resizing
	hwnd := statusWindow.Handle()
	win.SetWindowLong(hwnd, win.GWL_STYLE, win.GetWindowLong(hwnd, win.GWL_STYLE) & ^(win.WS_MAXIMIZEBOX|win.WS_MINIMIZEBOX|win.WS_SIZEBOX))

	// message loop
	statusWindow.Run()

	return nil
}

func NewTableDataModel() *TableDataModel {
	m := &TableDataModel{items: make([]*TableData, 0, 32)}

	// conditional message
	if gpsdata.getStatus() != "" {
		m.items = append(m.items, &TableData{
			Index: 0,
			Name:  "Message",
			Value: gpsdata.getStatus(),
		})
	}

	m.items = append(m.items, &TableData{
		Index: 1,
		Name:  "Gridsquare",
		Value: gpsdata.getLocation(),
	})

	m.items = append(m.items, &TableData{
		Index: 2,
		Name:  "Time",
		Value: gpsdata.getTime().Format("2006-01-02 15:04:05 UTC"),
	})

	m.items = append(m.items, &TableData{
		Index: 3,
		Name:  "Satellites",
		Value: fmt.Sprintf("%d", gpsdata.getNumSatellites()),
	})

	m.items = append(m.items, &TableData{
		Index: 4,
		Name:  "Fix Quality",
		Value: gpsdata.getFixQuality(),
	})

	m.items = append(m.items, &TableData{
		Index: 5,
		Name:  "HDOP",
		Value: strconv.FormatFloat(gpsdata.getHDOP(), 'f', -1, 64),
	})

	return m
}

type TableDataModel struct {
	walk.SortedReflectTableModelBase
	items []*TableData
}

func (m *TableDataModel) Items() interface{} {
	return m.items
}

type TableData struct {
	Index int
	Name  string
	Value string
}

func systemTray() error {
	var err error

	// load appIcon
	appIcon, err = walk.Resources.Icon("3")
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// our window
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// create our systray notify icon
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
	err = ni.SetIcon(appIcon)
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
	err = gridsquareAction.SetText("Copy Gridsquare")
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
	err = statusAction.SetText("Status...")
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	statusAction.Triggered().Attach(func() {
		err = RunStatusWindow()
		if err != nil {
			log.Printf("%+v", err)
		}
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

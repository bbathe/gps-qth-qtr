// +build windows

package main

import (
	"log"
	"time"
	"unsafe"

	"github.com/lxn/walk"
	declarative "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var (
	// pointer to SetSystemTime proc
	procSetSystemTime *windows.Proc

	// our icon
	appIcon *walk.Icon

	// prevent multiple status windows
	nbmRunStatusWindow = NewNonBlockingMutex()

	// reference to status window
	statusWindow *walk.MainWindow
)

// setSystemTime calls the windows SetSystemTime API
func setSystemTime(t time.Time) error {
	// convert time types
	systime := windows.Systemtime{
		Year:         uint16(t.Year()),
		Month:        uint16(t.Month()),
		Day:          uint16(t.Day()),
		Hour:         uint16(t.Hour()),
		Minute:       uint16(t.Minute()),
		Second:       uint16(t.Second()),
		Milliseconds: uint16(t.Nanosecond() / 1000000),
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

// runStatusWindow presents the user with a window containing the GPS data we have collected
func runStatusWindow() error {
	if nbmRunStatusWindow.Lock() {
		defer nbmRunStatusWindow.Unlock()

		mw := declarative.MainWindow{
			AssignTo: &statusWindow,
			Name:     "statusmw",
			Title:    "Status Data",
			Icon:     appIcon,
			Size:     declarative.Size{Width: 300, Height: 250},
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
							Model: newStatusTableDataModel(),
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

		// create window
		err := mw.Create()
		if err != nil {
			return err
		}

		// disable maximize, minimize, and resizing
		hwnd := statusWindow.Handle()
		win.SetWindowLong(hwnd, win.GWL_STYLE, win.GetWindowLong(hwnd, win.GWL_STYLE) & ^(win.WS_MAXIMIZEBOX|win.WS_MINIMIZEBOX|win.WS_SIZEBOX))

		// start message loop
		statusWindow.Run()
	} else {
		// bring already running status window to top
		statusWindow.Show()
	}

	return nil
}

// newStatusTableDataModel returns data model used to populate status tableview
func newStatusTableDataModel() *statusTableDataModel {
	m := &statusTableDataModel{items: make([]*statusTableData, 0, 8)}

	m.items = append(m.items, &statusTableData{
		Index: 0,
		Name:  "Message",
		Value: gpsdata.formatStatus(),
	})

	m.items = append(m.items, &statusTableData{
		Index: 1,
		Name:  "Gridsquare",
		Value: gpsdata.formatGridsquare(),
	})

	m.items = append(m.items, &statusTableData{
		Index: 2,
		Name:  "Latitude",
		Value: gpsdata.formatLatitude(),
	})

	m.items = append(m.items, &statusTableData{
		Index: 3,
		Name:  "Longitude",
		Value: gpsdata.formatLongitude(),
	})

	m.items = append(m.items, &statusTableData{
		Index: 4,
		Name:  "Time",
		Value: gpsdata.formatTime(),
	})

	m.items = append(m.items, &statusTableData{
		Index: 5,
		Name:  "Satellites",
		Value: gpsdata.formatNumSatellites(),
	})

	m.items = append(m.items, &statusTableData{
		Index: 6,
		Name:  "Fix Quality",
		Value: gpsdata.formatFixQuality(),
	})

	m.items = append(m.items, &statusTableData{
		Index: 7,
		Name:  "HDOP",
		Value: gpsdata.formatHDOP(),
	})

	return m
}

// statusTableDataModel is our datamodel type for status data
type statusTableDataModel struct {
	walk.SortedReflectTableModelBase
	items []*statusTableData
}

// Items is needed by status tableview
func (m *statusTableDataModel) Items() interface{} {
	return m.items
}

// statusTableData is our data type for status data
type statusTableData struct {
	Index int
	Name  string
	Value string
}

// systemTray create the UI element in the system tray for the user to interact with
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
		err := walk.Clipboard().SetText(gpsdata.formatGridsquare())
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
		err = runStatusWindow()
		if err != nil {
			log.Printf("%+v", err)
		}
	})
	err = ni.ContextMenu().Actions().Add(statusAction)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// update now action in context menu
	updateAction := walk.NewAction()
	err = updateAction.SetText("Update now")
	if err != nil {
		log.Printf("%+v", err)
		return err
	}
	updateAction.Triggered().Attach(func() {
		updated := gatherGpsData()
		if !updated {
			log.Printf("updateAction: failed to update")
		}
	})
	err = ni.ContextMenu().Actions().Add(updateAction)
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

	// start message loop
	mw.Run()

	return nil
}

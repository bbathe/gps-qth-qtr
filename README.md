# gps-qth-qtr

[![Tests](https://github.com/bbathe/gps-qth-qtr/workflows/Tests/badge.svg)](https://github.com/bbathe/gps-qth-qtr/actions) [![Release](https://github.com/bbathe/gps-qth-qtr/workflows/Release/badge.svg)](https://github.com/bbathe/gps-qth-qtr/actions)

gps-qth-qtr is a Windows system tray application that interfaces with a connected GPS receiver and keeps the system time accurate based on GPS time.  It also determines your [Maidenhead Gridsquare](https://en.wikipedia.org/wiki/Maidenhead_Locator_System) from the geolocation information.

## Description

I created this application to keep my laptop's time correct when I'm "off the grid" during mobile [Amateur Radio](http://www.arrl.org) activities using protocols where having the correct time is important like [FT4 and FT8](https://www.physics.princeton.edu/pulsar/k1jt/wsjtx.html).  While time is important with these protocols, the accuracy required is only within a couple seconds.  I'm not doing anything fancy to keep the system time more accurate than required by these protocols.

I did all the inital development using a [u-blox 8](https://www.u-blox.com) reciever from [Amazon](https://smile.amazon.com/gp/product/B071XY4R26).  This application only uses the 2 [NMEA 0183](https://en.wikipedia.org/wiki/NMEA_0183) sentences GGA and RMC and does not restrict the NMEA 0183 talker, so it should work with any navigation satellite system reciever as long as you can get the correct drivers installed so the data can be read from a COM port.

## Installation

To install this application:

1. Create the folder `C:\Program Files\gps-qth-qtr`
2. Download the ```gps-qth-qtr.exe.zip``` file from the [latest release](https://github.com/bbathe/gps-qth-qtr/releases) and unzip it into that folder
3. Create the ```gps-qth-qtr.yaml``` file (a [YAML](https://en.wikipedia.org/wiki/YAML) file is a plain text file with the ```.yaml``` extension) in that folder, with these attributes:
    ```
    gpsdevice:
      port: COM3
      baud: 9600
      pollrate: 900
    ```
    - ```port``` is the name of the Windows COM port to read from the connected GPS device, this is setup when you install the device driver for your GPS device.  You should be able to find this in Device Manager.
    - ```baud``` is the rate at which information is transferred from the COM port, this is a setting on the port that is setup when you install the device driver for your GPS device.  You should be able to find this in Device Manager, check the "Port Settings" tab for the device.
    - ```pollrate``` defines how often (in seconds) you want the gps-qth-qtr application to poll the connected GPS device and set the system time.
4. You can now double-click on the ```gps-qth-qtr.exe``` file to start the application.

There will be a log file created in the same directory as the executable and all errors are logged there.

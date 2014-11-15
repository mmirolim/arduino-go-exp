package main

import (
    "log"
    "time"

    "github.com/hybridgroup/gobot"
    "github.com/hybridgroup/gobot/platforms/firmata"
    "github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
    gbot := gobot.NewGobot()

    firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "COM6")
    relay := gpio.NewDirectPinDriver(firmataAdaptor, "pin", "7")

    work := func() {
        level := byte(1)

        gobot.Every(2*time.Second, func() {
            log.Println(level)
            relay.DigitalWrite(level)
            if level == 1 {
                level = 0
            } else {
                level = 1
            }
        })
    }

    robot := gobot.NewRobot("bot",
        []gobot.Connection{firmataAdaptor},
        []gobot.Device{relay},
        work,
    )

    gbot.AddRobot(robot)

    gbot.Start()
}
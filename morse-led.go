package main 

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
	"log"
	"strings"
	"strconv"
)

func main() {
	// morse flash code 1 is dot and 2 is flash
	morse := map[string]string{
		"A" : "12",
		"B" : "2111",
		"C" : "2121",
		"D" : "211",
		"E" : "1",
		"F" : "1121",
		"G" : "221",
		"H" : "1111",
		"I" : "11",
		"J" : "1222",
		"K" : "212",
		"L" : "1211",
		"M" : "22",
		"N" : "21",
		"O" : "222",
		"P" : "1221",
		"Q" : "2212",
		"R" : "121",
		"S" : "111",
		"T" : "2",
		"U" : "112",
		"V" : "1112",
		"W" : "122",
		"X" : "2112",
		"Y" : "2122",
		"Z" : "2211",
	}
	// timings/delay
	d := map[int]int{
		0 : 200,  // dot
		1 : 200,  // between parts of same letter
		2 : 600,  // between two letters
		3 : 1400, // between words
	}
	// string to flash
	msg := "HELLO WORLD"
	
	gbot := gobot.NewGobot()
	server := api.NewAPI(gbot)
	server.Port = "4000"
	server.Start()

    firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "COM6")
    led := gpio.NewLedDriver(firmataAdaptor, "led", "13")
    work := func() {
        gobot.Every(15*time.Second, func() {
        	log.Println(msg)
        	for _, v := range strings.Split(msg, "") {
        		// check if space
        		if v == " " {
        			time.Sleep(time.Duration(d[3])*time.Millisecond)
        		} else {
        			time.Sleep(time.Duration(d[2])*time.Millisecond)
        		}
        		for _, chr := range strings.Split(morse[v], "") {
        			time.Sleep(time.Duration(d[1])*time.Millisecond)
        			led.On()
        			i, _ := strconv.Atoi(chr)
        			time.Sleep(time.Duration(d[i])*time.Millisecond)
        			led.Off()
        		}
        	}
        })
    }

    robot := gobot.NewRobot(
    	"bot",
        []gobot.Connection{firmataAdaptor},
        []gobot.Device{led},
        work,
    )

    gbot.AddRobot(robot)

    gbot.Start()
}
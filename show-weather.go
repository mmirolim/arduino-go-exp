package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tarm/goserial"
)

func main() {
	cities := []string{
		"Bangkok",
		"Beijing",
		"Bogotá",
		"Buenos Aires",
		"Cairo", "Delhi",
		"Dhaka",
		"Guangzhou",
		"Istanbul",
		"Jakarta",
		"Karachi",
		"Kinshasa",
		"Kolkata",
		"Lagos",
		"Lahore",
		"London",
		"Los Angeles",
		"Manila",
		"Mexico",
		"Moscow",
		"Mumbai",
		"New York",
		"Osaka",
		"Rio de Janeiro",
		"São Paulo",
		"Seoul",
		"Shenzhen",
		"Shanghai",
		"Tianjin",
		"Tokyo",
		"Tashkent",
		"Quebec"}

	var hotCities int
	// send to arduino via serial port
	c := &serial.Config{Name: "COM6", Baud: 9600}
	s, err := serial.OpenPort(c)
	time.Sleep(2 * time.Second)
	failOnError(err, "Err on serial port open")

	for _, v := range cities {
		uri := "http://api.openweathermap.org/data/2.5/weather?q=" + v + "&units=metric"
		res, err := http.Get(uri)
		failOnError(err, "Err during making http request")

		data, err := ioutil.ReadAll(res.Body)
		failOnError(err, "Err in ReadAll")

		var dat map[string]interface{}

		if err := json.Unmarshal(data, &dat); err != nil {
			log.Println(err)
			continue
		}
		main := dat["main"].(map[string]interface{})
		temp := main["temp"].(float64)
		fmt.Println(v, "temp", temp)

		tempStr := strconv.FormatFloat(temp, 'f', -1, 64)
		// write to serial port
		n, err := s.Write([]byte(" " + v + "/ Temp " + tempStr + " C"))
		failOnError(err, "Err on Serial write")
		log.Println(n)

		buf := make([]byte, 128)
		n, err = s.Read(buf)
		failOnError(err, "Err during read from serial port")
		if n > 3 {
			for _, chr := range buf[n-3 : n] {
				hotCities = int(chr) - 48
			}
			log.Print(hotCities)
		}
		time.Sleep(1 * time.Second)
	}
	// write to serial port
	n, err := s.Write([]byte(" " + "Hot Cities" + "/ total " + strconv.Itoa(hotCities)))
	failOnError(err, "Err on Serial write")
	log.Println(n)
	time.Sleep(5 * time.Second)

}

func failOnError(e error, msg string) {
	if e != nil {
		log.Fatal("%s: %s", msg, e)
		panic(fmt.Sprintf("%s: %s", msg, e))
	}
}

package main

import (
      "github.com/tarm/goserial"
      "log"
      "time"
)

func main() {
        c := &serial.Config{Name: "COM6", Baud: 9600}
        s, err := serial.OpenPort(c)
        time.Sleep(2 * time.Second)
        if err != nil {
                log.Fatal(err)
        }

        n, err := s.Write([]byte("I am Mirolim /second row"))
        if err != nil {
                log.Fatal(err)
        }
        log.Println(n)

        time.Sleep(2 * time.Second)
}
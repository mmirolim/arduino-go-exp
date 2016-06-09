package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	port = flag.String("p", "/dev/ttyACM0", "port name")
)

func main() {
	flag.Parse()
	// connect to serial port
	f, err := connectSerial(*port)
	exitOnErr("serial port connection failed", err)
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	// start listening std in and out data from serial port
	listen := bufio.NewReader(os.Stdin)
	// read data from std in to write into serial port
	go func() {
		buf := make([]byte, 128)
		var readCount int
		for {
			n, err := f.Read(buf)
			if err != nil {
				exitOnErr("port read err", err)
			}

			readCount++
			fmt.Printf("Read %v %v bytes: % 02x %s\n", readCount, n, buf[:n], buf[:n])

		}
	}()
	for {
		s, err := listen.ReadString('\n')
		exitOnErr("readstring", err)
		switch s {
		case "exit\n":
			log.Println("exiting")
			break
		default:
			// write read string to serial port
			s = s[:len(s)-1]
			fmt.Println("write to serial", s, []byte(s))
			n, err := f.Write([]byte(s))
			exitOnErr("port write err", err)
			fmt.Println("bytes written ", n)
		}
	}
}

func connectSerial(port string) (*os.File, error) {
	const (
		baud       uint32 = 9600
		size              = 8
		stopBits          = 1
		parityNone        = 'N'
		timeout           = time.Millisecond * 500
	)

	var (
		f   *os.File
		err error
	)

	f, err = os.OpenFile(port, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		return nil, err
	}

	cflagToUse := syscall.CREAD | syscall.CLOCAL | baud | syscall.CS8

	fd := f.Fd()
	vmin, vtime := posixTimeoutValues(timeout)

	t := syscall.Termios{
		Iflag:  syscall.IGNPAR,
		Cflag:  cflagToUse,
		Cc:     [32]uint8{syscall.VMIN: vmin, syscall.VTIME: vtime},
		Ispeed: baud,
		Ospeed: baud,
	}

	if _, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&t)),
		0,
		0,
		0,
	); errno != 0 {
		return nil, errno
	}

	if err = syscall.SetNonblock(int(fd), false); err != nil {
		return nil, err
	}

	return f, err
}

// code from tarm serial github
func posixTimeoutValues(t time.Duration) (vmin uint8, vtime uint8) {
	const MAXUINT8 = 1<<8 - 1 // 255
	// set blocking / non blocking read
	var minBytesToRead uint8 = 1
	var readTimeoutInDeci int64
	if t > 0 {
		// EOF on zero read
		minBytesToRead = 1
		// convert timeout to deciseconds as expected by VTIME
		readTimeoutInDeci = (t.Nanoseconds() / 1e8)
		// capping the timeout
		if readTimeoutInDeci < 1 {
			// min possible timeout 0.1s
			readTimeoutInDeci = 1
		} else if readTimeoutInDeci > MAXUINT8 {
			// max possible timeout 25.5s
			readTimeoutInDeci = MAXUINT8
		}
	}

	return minBytesToRead, uint8(readTimeoutInDeci)
}

// Discards data written to the port but not transmitted,
// or data received but not read
func flush(f *os.File) error {
	const TCFLSH = 0x540B
	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(f.Fd()),
		uintptr(TCFLSH),
		uintptr(syscall.TCIOFLUSH),
	)
	return err
}

func exitOnErr(msg string, err error) {
	if err != nil {
		log.Println(msg, err)
		os.Exit(1)
	}
}

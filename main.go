package main

import (
	"fmt"
	"github.com/tdaira/food-manager/device"
	"github.com/tdaira/food-manager/gcloud"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var stepperSleep = 100 * time.Millisecond // msec
var cameraSleep = 10 * time.Second        // sec

func main() {
	// Init RaspberryPi adaptor.
	adaptor := raspi.NewAdaptor()

	// Turn on LED.
	led := device.NewLED(adaptor, "29")
	err := led.ON()
	if err != nil {
		log.Fatal(err)
	}

	// Run Motor.
	motor := device.NewMotor(adaptor, 50*time.Millisecond)
	motor.Run()

	// Init cloud Storage.
	storage, err := gcloud.NewStorage("food_watcher")
	if err != nil {
		log.Fatal(err)
	}

	// Set signal handler for device termination.
	setSignalHandler(led, motor)

	// Read switch value and change motor direction.
	for {
		for i := 0; i < int(cameraSleep/stepperSleep); i++ {
			val1, err := adaptor.DigitalRead("19")
			if err != nil {
				panic(err)
			}
			val2, err := adaptor.DigitalRead("21")
			if err != nil {
				panic(err)
			}
			if val1 == 1 {
				motor.SetDirection(true)
			}
			if val2 == 1 {
				motor.SetDirection(false)
			}
			time.Sleep(100 * time.Millisecond)
		}
		motor.Stop()
		path := createPath()
		err := takePhoto(path)
		if err != nil {
			log.Fatal(err)
		}
		storage.Upload(path)
		if err != nil {
			log.Fatal(err)
		}
		motor.Run()
	}
}

func setSignalHandler(led *device.LED, motor *device.Motor) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				err := led.OFF()
				if err != nil {
					log.Fatal(err)
				}
				err = motor.Stop()
				if err != nil {
					log.Fatal(err)
				}
				os.Exit(0)

			default:
				fmt.Println("Unknown signal.")
				os.Exit(1)
			}
		}
	}()
}

func createPath() string {
	now := time.Now()
	secs := now.Unix()
	return "/tmp/" + strconv.Itoa(int(secs)) + ".jpg"
}

func takePhoto(path string) error {
	log.Print("Take photo: " + path)
	return exec.Command("raspistill", "-o", path).Run()
}

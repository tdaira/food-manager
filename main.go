package main

import (
	"fmt"
	"github.com/tdaira/food-manager/device"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
	motor := device.NewMotor(adaptor, 20*time.Millisecond)
	motor.Run()

	// Set signal handler for device termination.
	setSignalHandler(led, motor)

	// Read switch value and change motor direction.
	for {
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

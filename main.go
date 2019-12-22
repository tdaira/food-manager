package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

type Direction int

const (
	Forward Direction = iota
	Backward
)

func (d Direction) String() string {
	switch d {
	case Forward:
		return "forward"
	case Backward:
		return "backward"
	default:
		return "forward"
	}
}

func main() {
	adaptor := raspi.NewAdaptor()
	stepper := gpio.NewStepperDriver(adaptor, [4]string{"31", "33", "35", "37"}, gpio.StepperModes.DualPhaseStepping, 2048)
	direction := Forward
	err := stepper.SetSpeed(5)
	if err != nil {
		fmt.Printf("set speed error: %+v", err)
		return
	}
	err = stepper.SetDirection(Forward.String())
	if err != nil {
		fmt.Printf("set direction error: %+v", err)
		return
	}
	err = stepper.Run()
	if err != nil {
		fmt.Printf("stepper execution error: %+v", err)
		return
	}

	work := func() {
		gobot.Every(5*time.Second, func() {
			if direction == Forward {
				direction = Backward
			} else {
				direction = Forward
			}
			err := stepper.SetDirection(direction.String())
			if err != nil {
				fmt.Printf("stepper execution error: %+v", err)
			}
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{stepper},
		work,
	)

	err = robot.Start()
	if err != nil {
		fmt.Printf("bot execution error: %+v", err)
		return
	}
}

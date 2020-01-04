package device

import (
	"gobot.io/x/gobot/platforms/raspi"
)

type LED struct {
	adaptor *raspi.Adaptor
	pin     string
}

func NewLED(adaptor *raspi.Adaptor, pin string) *LED {
	return &LED{adaptor: adaptor, pin: pin}
}

func (l *LED) ON() error {
	return l.adaptor.DigitalWrite(l.pin, 1)
}

func (l *LED) OFF() error {
	return l.adaptor.DigitalWrite(l.pin, 0)
}

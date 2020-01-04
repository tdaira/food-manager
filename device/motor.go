package device

import (
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"time"
)

type Motor struct {
	adaptor   *raspi.Adaptor
	delay     time.Duration
	steps     [][4]byte
	pins      [4]string
	running   bool
	direction bool
}

func NewMotor(adaptor *raspi.Adaptor, delay time.Duration) *Motor {
	dualPhaseStepping := [][4]byte{
		{1, 0, 0, 1},
		{1, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 1},
	}
	return &Motor{
		adaptor,
		delay,
		dualPhaseStepping,
		[4]string{"31", "33", "35", "37"},
		false,
		false,
	}
}

func (m *Motor) Run() {
	m.running = true

	go func() {
		for {
			if !m.running {
				break
			}
			for _, step := range m.steps {
				for i := 0; i < len(step); i++ {
					if m.direction {
						err := m.adaptor.DigitalWrite(m.pins[i], step[len(step)-i-1])
						if err != nil {
							log.Fatal(err)
						}
					} else {
						err := m.adaptor.DigitalWrite(m.pins[i], step[i])
						if err != nil {
							log.Fatal(err)
						}
					}
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}

func (m *Motor) SetDirection(direction bool) {
	m.direction = direction
}

func (m *Motor) Stop() error {
	m.running = false
	powerOffSignal := [4]byte{0, 0, 0, 0}
	for i, v := range powerOffSignal {
		return m.adaptor.DigitalWrite(m.pins[i], v)
	}
	return nil
}

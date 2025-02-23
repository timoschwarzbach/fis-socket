package content

import (
	"time"
)

type GenericSequence struct {
	controller *Controller
	screen     string
	duration   int
}

func (c *Controller) Generic(screen string, duration ...int) *GenericSequence {
	dur := 5
	if len(duration) > 0 {
		dur = duration[0]
	}
	return &GenericSequence{
		controller: c,
		screen:     screen,
		duration:   dur,
	}
}

func (i *GenericSequence) Display() {
	i.controller.send(i.screen, "")

	time.Sleep(time.Duration(i.duration) * time.Second)

	i.controller.next()
}

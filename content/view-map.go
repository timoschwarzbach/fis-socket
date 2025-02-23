package content

import (
	"time"
)

type MapSequence struct {
	controller *Controller
	source     string
}

func (c *Controller) MapView(file string) *MapSequence {
	if file == "" {
		file = "test1"
	}
	return &MapSequence{
		controller: c,
		source:     file + ".json",
	}
}

func (i *MapSequence) Display() {

	i.controller.send("mapview", "http://localhost:8080/mapscreens/"+i.source)

	// wait 15 seconds
	time.Sleep(15 * time.Second)

	i.controller.next()
}

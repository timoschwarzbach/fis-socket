package sequences

import (
	"fis/socket/socket"
)

type Sequence interface {
	Display()
}

type Controller struct {
	socketServer    *socket.Server
	sequenceService *SequenceService
	sequence        []Sequence
	index           int
}

func CreateController(server *socket.Server, dbSync chan bool) *Controller {
	c := &Controller{
		index:           0,
		socketServer:    server,
		sequenceService: CreateSequenceService(dbSync),
	}

	c.sequence = c.demoSequence()

	c.next()

	return c
}

func (c *Controller) demoSequence() []Sequence {
	// img := c.Image()
	img := c.Sequence(c.sequenceService)
	return []Sequence{
		c.Generic("stations"),
		// c.Generic("map", 2),
		// c.BocholtBorkenerVolksblatt(),
		// c.MapView("test1"),
		// c.MapView("testbustreff"),
		// c.MapView("testbahnhof"),
		img,
		img,
		// c.Tagesschau(),
		// c.Video(),
	}

}

type Packet struct {
	Screen string `json:"screen"`
	Data   string `json:"data"`
}

// allows sequence items to send messages to the fis displays
func (c *Controller) send(screen string, data any) {
	// fmt.Printf("Sending %s\n", screen)
	c.socketServer.Emit("screen", screen, data)
}

// continue the sequence
func (c *Controller) next() {
	c.index = (c.index + 1) % len(c.sequence)
	c.sequence[c.index].Display()
}

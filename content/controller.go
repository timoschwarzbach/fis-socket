package content

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
	// item := c.Sequence(c.sequenceService)
	return []Sequence{
		// c.Generic("stations"),
		c.Generic("map", 2),
		// c.Generic("map-fullRoute", 5),
		// c.MapView("test1"),
		// c.MapView("testbustreff"),
		// c.MapView("testbahnhof"),
		// item,
		// item,
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

// allows a sequence to tell the client to load a resource into cache for display in the future
func (c *Controller) prefetch(resource string) {
	c.socketServer.Emit("prefetch", resource)
}

// continue the sequence
func (c *Controller) next() {
	c.index = (c.index + 1) % len(c.sequence)
	c.sequence[c.index].Display()
}

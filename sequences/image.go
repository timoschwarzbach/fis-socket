package sequences

import (
	"time"
)

type ImageSequence struct {
	controller *Controller
	imagePool  []string
	current    int
}

func (c *Controller) Image() *ImageSequence {
	return &ImageSequence{
		controller: c,
		imagePool: []string{
			"bsag-aktuelles-whatsapp-kanal-abonnieren.jpg",
			"csm_bsag-header-mit-bus-und-bahn-zum-werderspiel_e43d540dfc.jpg",
		},
		current: 0,
	}
}

func (i *ImageSequence) Display() {
	i.controller.send("image", "http://localhost:8080/"+i.imagePool[i.current])

	// wait 5 seconds
	time.Sleep(5 * time.Second)

	i.prepareNext()
	i.controller.next()
}

func (i *ImageSequence) prepareNext() {
	i.current = (i.current + 1) % len(i.imagePool)
}

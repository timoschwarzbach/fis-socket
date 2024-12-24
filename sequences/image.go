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
			"image1.jpg",
			"image2.jpg",
			"image3.jpg",
			"20231004_115158617_iOS.jpg",
			"202308xx_Gleisanschluss-Neues-Werk-Cottbus.jpg",
			"f8b8fb494f0b00379e402c5162fb00d9.jpeg",
			"Kampagnenmotiv_Nur-fuer-alle_S-Bahn-Berlin_Gemeinsam-data.jpg",
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

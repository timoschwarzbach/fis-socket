package sequences

import (
	"fmt"
	"time"
)

type BBVSequence struct {
	controller *Controller
	newsPool   []News
	current    int
}

func (c *Controller) BocholtBorkenerVolksblatt() *BBVSequence {
	return &BBVSequence{
		controller: c,
		newsPool: []News{
			{topline: "Gigaset in neuen Händen", headline: "Neustart in Bocholt unter Hongkonger Führung", image: "630_0900_4307391_gigaset-1648x824.jpg"},
			{topline: "Rentner macht einen von 350 in Bocholt reich", headline: "Christoph Prüm startet ungewöhnliches Experiment", image: "630_0900_4306357_Grunderbe_003-1648x824.jpg"},
		},
		current: 0,
	}
}

func (i *BBVSequence) Display() {
	if i.current > len(i.newsPool) {
		fmt.Println("No news available")
		i.prepareNext()
		i.controller.next()
		return
	}

	// todo bbv intro

	news := i.newsPool[i.current]
	packet := NewsPacket{Topline: news.topline, Headline: news.headline, Image: "http://localhost:8080/bbv/" + news.image, Badge: "http://localhost:8080/bbv/logo.svg"}
	i.controller.send("news", packet)

	// wait 5 seconds
	time.Sleep(5 * time.Second)

	i.prepareNext()
	i.controller.next()
}

func (i *BBVSequence) prepareNext() {
	i.current = (i.current + 1) % len(i.newsPool)
}

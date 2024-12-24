package sequences

import (
	"fmt"
	"path/filepath"
	"time"
)

type News struct {
	topline  string
	headline string
	image    string
}

type NewsPacket struct {
	Topline  string `json:"topline"`
	Headline string `json:"headline"`
	Image    string `json:"image"`
	Badge    string `json:"badge"`
}

type TagesschauSequence struct {
	controller *Controller
	newsPool   []News
	current    int
}

func (c *Controller) Tagesschau() *TagesschauSequence {
	return &TagesschauSequence{
		controller: c,
		newsPool: []News{
			{topline: "Misstrauensvotum", headline: "Opposition stürzt Frankreichs Regierung", image: "barnier-204.webp"},
			{topline: "Minderheitsregierung in Sachsen", headline: "CDU und SPD stellen Koalitionsvertrag vor", image: "mdr-petra-koepping-116.webp"},
			{topline: "Einsatz von 500 Bundespolizisten", headline: "Großrazzia gegen Netzwerk von Schleusern", image: "polizei-razzia-106.webp"},
		},
		current: 0,
	}
}

func (i *TagesschauSequence) Display() {
	if i.current > len(i.newsPool) {
		fmt.Println("No news available")
		i.prepareNext()
		i.controller.next()
		return
	}

	i.showTagesschauIntro()

	news := i.newsPool[i.current]
	packet := NewsPacket{Topline: news.topline, Headline: news.headline, Image: "http://localhost:8080/tagesschau/" + news.image, Badge: "http://localhost:8080/tagesschau/badge.png"}
	i.controller.send("news", packet)

	// wait 5 seconds
	time.Sleep(5 * time.Second)

	i.prepareNext()
	i.controller.next()
}

func (i *TagesschauSequence) prepareNext() {
	i.current = (i.current + 1) % len(i.newsPool)
}

func (i *TagesschauSequence) showTagesschauIntro() {
	i.controller.send("video", "http://localhost:8080/tagesschau/tagesschau_intro.mp4")
	absPath, _ := filepath.Abs("./static/tagesschau/tagesschau_intro.mp4")
	time.Sleep(getVideoDuration(absPath))
}

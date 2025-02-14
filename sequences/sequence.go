package sequences

import (
	"encoding/json"
	"time"
)

type RemoteSequence struct {
	controller *Controller
	service    *SequenceService
}

func (c *Controller) Sequence(s *SequenceService) *RemoteSequence {
	return &RemoteSequence{
		controller: c,
		service:    s,
	}
}

// {"aspects":["aspect-16-9"],"slides":[{"background":"cm72mbowb0000cg0m0jjf1p5z","bottom":{}}]}
type sequenceData struct {
	Slides []sequenceSlide `json:"slides"`
}

type sequenceSlide struct {
	Background string                 `json:"background"`
	Bottom     map[string]interface{} `json:"bottom"`
}

func (rs *RemoteSequence) Display() {
	var displayData sequenceData
	err := json.Unmarshal([]byte(rs.service.current.displayJSON), &displayData)
	if err != nil {
		panic(err)
	}
	for index := range displayData.Slides {
		slide := displayData.Slides[index]
		rs.controller.send("image", "http://localhost:8080/"+slide.Background)
		time.Sleep(5 * time.Second)
	}

	// wait 5 seconds
	time.Sleep(5 * time.Second)

	rs.controller.next()
}

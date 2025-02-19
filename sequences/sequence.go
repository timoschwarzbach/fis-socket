package sequences

import (
	"log"
	"path/filepath"
	"time"

	"github.com/goccy/go-json"
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

type slide struct {
	Background string                 `json:"background"`
	Bottom     map[string]interface{} `json:"bottom"`
	Duration   StringInt              `json:"duration"`
}

func (rs *RemoteSequence) Display() {
	log.Println("Sequence:\tDisplaying a remote sequence")
	var slides []slide
	err := json.Unmarshal([]byte(rs.service.current.Slides), &slides)
	if err != nil {
		log.Panicln("Sequence: Error unmarshalling slides")
	}
	for index := range slides {
		slide := slides[index]
		log.Printf("Sequence:\tSending slide %d of %d\n", index+1, len(slides))
		fileCategory, backgroundFile := rs.service.getLocalFileReferenceFromId(slide.Background)
		rs.controller.send(fileCategory, "http://localhost:8080/"+backgroundFile)

		if slide.Duration > 0 {
			time.Sleep(time.Duration(slide.Duration) * time.Millisecond)
		} else if fileCategory == "video" {
			absPath, _ := filepath.Abs("./static/" + backgroundFile)
			time.Sleep(getVideoDuration(absPath))
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Sequence:\tProgressing to next sequence")
	rs.service.Step()
	rs.controller.next()
}

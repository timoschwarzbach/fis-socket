package content

import (
	"log"
	"os/exec"
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

type mediaData struct {
	Type       string                 `json:"type"`
	Background string                 `json:"background"`
	Bottom     map[string]interface{} `json:"bottom"`
}

func (rs *RemoteSequence) Display() {
	sequence := rs.service.current
	log.Println("Sequence:\tDisplaying a remote sequence", sequence.Category)

	// tell client to load all images
	// in the best case: this has already been done trough the previous sequence
	sendFetchInstruction(rs.service.current, rs.controller)
	sendFetchInstruction(rs.service.next, rs.controller)

	// display slides
	for index := range sequence.Slides {
		slide := sequence.Slides[index]
		log.Printf("Sequence:\tSending slide %d of %d\n", index+1, len(sequence.Slides))

		fileCategory, backgroundFile := "fallback", "fallback"

		if sequence.Category == "default" {
			fileCategory, backgroundFile = rs.service.fileFromId(slide.Background)
		}

		if sequence.Category == "tagesschau" {
			fileCategory, backgroundFile = "image", slide.Background
		}

		// send item
		rs.controller.send("media", &mediaData{
			Type:       fileCategory,
			Background: backgroundUrl(rs.service.current.Category, backgroundFile),
			Bottom:     slide.Bottom,
		})

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

func sendFetchInstruction(s *Sequences, c *Controller) {
	for index := range s.Slides {
		slide := s.Slides[index]
		c.prefetch(slide.Background)
	}
}

/* Video duration */

type fFProbeOutput struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func getVideoDuration(filename string) time.Duration {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "json", filename)
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error getting video duration:", err)
		return 0
	}

	var ffprobeOutput fFProbeOutput
	if err := json.Unmarshal(output, &ffprobeOutput); err != nil {
		log.Println("Error getting video duration:", err)
		return 0
	}

	duration, err := time.ParseDuration(ffprobeOutput.Format.Duration + "s")
	if err != nil {
		log.Println("Error parsing video duration:", err)
		return 0
	}
	return duration
}

/* Background url from sequence category + media id */
func backgroundUrl(category string, mediaId string) string {
	if category == "tagesschau" {
		return "http://localhost:8080/tagesschau/" + mediaId
	}
	return "http://localhost:8080/" + mediaId
}

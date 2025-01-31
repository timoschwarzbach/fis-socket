package sequences

import (
	"encoding/json"
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

type VideoSequence struct {
	controller *Controller
	videoPool  []string
	current    int
}

func (c *Controller) Video() *VideoSequence {
	return &VideoSequence{
		controller: c,
		videoPool:  []string{"vimeo_558475288_1920x1080_mute.mp4"},
		current:    0,
	}
}

func (i *VideoSequence) Display() {
	i.controller.send("video", "http://localhost:8080/"+i.videoPool[i.current])

	// wait for the duration of the video

	absPath, _ := filepath.Abs("./static/" + i.videoPool[i.current])
	time.Sleep(getVideoDuration(absPath))

	i.prepareNext()
	i.controller.next()
}

func (i *VideoSequence) prepareNext() {
	i.current = (i.current + 1) % len(i.videoPool)
}

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

package wallpaper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// Image is the path to a image file that can be used to as a background.
type Image string

func (i Image) String() string {
	return string(i)
}

// Mode defines how to position or transform the images for the background.
type Mode string

// Supported modes for image scaling.
// See `man feh` for more details.
const (
	Center Mode = "center"
	Fill   Mode = "fill"
	Max    Mode = "max"
)

// Background is used to manage the Image that is displayed as the wallpaper.
type Background struct {
	Mode Mode
}

// Set will set the Background to the provided Image.
func (w *Background) Set(image Image) error {
	mode := fmt.Sprintf("--bg-%s", w.Mode)
	cmd := exec.Command("feh", mode, image.String())
	return cmd.Run()
}

// IsImage determines if a given file is an Image that can be used for a background.
func IsImage(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		// if we can't open the file, then it definitely is not an image
		log.Printf("[INFO] unable to open file: %s", path)
		return false
	}
	defer f.Close()

	// Only need the first 512 bytes to determine content type
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		// if we can't get those bytes, it is not an image
		log.Printf("[INFO] unable to read bytes to determine type: %s", path)
		return false
	}
	contentType := http.DetectContentType(buffer)
	log.Printf("[DEBUG] file %s is %s", path, contentType)
	switch contentType {
	case "image/jpeg":
		return true
	case "image/png":
		return true
	default:
		return false
	}
}

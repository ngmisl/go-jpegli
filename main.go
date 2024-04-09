package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gen2brain/jpegli"
)

// ToJpeg converts a PNG image to Jpegli format
func ToJpeg(imageBytes []byte) ([]byte, error) {
	// Detect the content type
	contentType := http.DetectContentType(imageBytes)
	if contentType != "image/png" {
		return nil, fmt.Errorf("unable to convert %#v to Jpegli", contentType)
	}

	// Decode the PNG image bytes
	img, err := png.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	// Encode the image as a Jpegli file
	buf := new(bytes.Buffer)
	if err := jpegli.Encode(buf, img, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	// Check if any command-line arguments are provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [PNG_FILE_1] [PNG_FILE_2] ... [PNG_FILE_N]")
	}

	for _, arg := range os.Args[1:] {
		// Read the input PNG image from a file
		inputImageBytes, err := os.ReadFile(arg)
		if err != nil {
			log.Printf("Error reading PNG file %s: %v", arg, err)
			continue
		}

		// Convert PNG to Jpegli
		convertedImageBytes, err := ToJpeg(inputImageBytes)
		if err != nil {
			log.Printf("Error converting PNG %s to Jpegli: %v", arg, err)
			continue
		}

		// Construct the output file name
		outputFilename := filepath.Base(arg[:len(arg)-len(filepath.Ext(arg))]) + ".jpeg"
		if err := os.WriteFile(outputFilename, convertedImageBytes, 0644); err != nil {
			log.Printf("Error saving Jpegli file %s: %v", outputFilename, err)
			continue
		}

		log.Printf("Conversion of %s to Jpegli completed successfully.", arg)
	}
}

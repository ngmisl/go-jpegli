package main

import (
	"bytes"
	"fmt"
	"image/png" // Import PNG support
	"log"
	"net/http"
	"os"

	"github.com/gen2brain/jpegli" // Import Jpegli package
)

// ToJpeg converts a PNG image to Jpegli format
func ToJpeg(imageBytes []byte) ([]byte, error) {
	// Detect the content type
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/png":
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

	return nil, fmt.Errorf("unable to convert %#v to Jpegli", contentType)
}

func main() {
	// Read the input PNG image from a file
	inputImageBytes, err := os.ReadFile("1.png")
	if err != nil {
		log.Fatal("Error reading PNG file:", err)
	}

	// Convert PNG to Jpegli
	convertedImageBytes, err := ToJpeg(inputImageBytes)
	if err != nil {
		log.Fatal("Error converting PNG to Jpegli:", err)
	}

	// Write the Jpegli output to a file
	if err := os.WriteFile("output.jpeg", convertedImageBytes, 0644); err != nil {
		log.Fatal("Error saving Jpegli file:", err)
	}

	log.Println("Conversion to Jpegli completed successfully.")
}

package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
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
		inputImageBytes, err := os.ReadFile(filepath.Clean(arg))
		if err != nil {
			log.Printf("Error reading PNG file %s: %v", arg, err)
			continue
		}

		// Convert PNG to Jpegli
		jpegliBytes, err := ToJpeg(inputImageBytes)
		if err != nil {
			log.Printf("Error converting PNG %s to Jpegli: %v", arg, err)
			continue
		}

		// Convert PNG to JPEG using standard compression
		jpegBuf := new(bytes.Buffer)
		img, err := png.Decode(bytes.NewReader(inputImageBytes))
		if err != nil {
			log.Printf("Error decoding PNG %s: %v", arg, err)
			continue
		}
		if err := jpeg.Encode(jpegBuf, img, &jpeg.Options{Quality: 100}); err != nil {
			log.Printf("Error encoding JPEG %s: %v", arg, err)
			continue
		}

		// Calculate the compression difference
		pngSize := len(inputImageBytes)
		jpegliSize := len(jpegliBytes)
		jpegSize := len(jpegBuf.Bytes())
		jpegliCompressionRatio := float64(jpegliSize) / float64(pngSize) * 100
		jpegCompressionRatio := float64(jpegSize) / float64(pngSize) * 100

		fmt.Printf("File: %s\n", arg)
		fmt.Printf("PNG size: %d bytes\n", pngSize)
		fmt.Printf("Jpegli size: %d bytes (%.2f%% of PNG)\n", jpegliSize, jpegliCompressionRatio)
		fmt.Printf("JPEG size: %d bytes (%.2f%% of PNG)\n", jpegSize, jpegCompressionRatio)
		fmt.Printf("Jpegli compression is %.2f%% smaller than JPEG compression\n", jpegCompressionRatio-jpegliCompressionRatio)

		// Construct the output file name
		outputDir := filepath.Dir(arg)
		outputFilename := filepath.Base(arg[:len(arg)-len(filepath.Ext(arg))]) + ".jpeg"
		outputPath := filepath.Join(outputDir, outputFilename)

		if err := os.WriteFile(filepath.Clean(outputPath), jpegliBytes, 0644); err != nil {
			log.Printf("Error saving Jpegli file %s: %v", outputPath, err)
			continue
		}

		log.Printf("Conversion of %s to Jpegli completed successfully.", arg)
	}
}

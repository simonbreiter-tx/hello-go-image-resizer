// main.go
package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"log"
	"os"

	"github.com/disintegration/gift"
)

func main() {
	filePath := flag.String("file", "", "Path to the GIF file")
	cropX := flag.Int("x", 0, "X coordinate for cropping")
	cropY := flag.Int("y", 0, "Y coordinate for cropping")
	cropWidth := flag.Int("width", 100, "Width of the cropped region")
	cropHeight := flag.Int("height", 100, "Height of the cropped region")
	outputFile := flag.String("output", "cropped.gif", "Output file name")

	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide a GIF file path using the -file flag.")
		os.Exit(1)
	}

	err := processGIF(*filePath, *cropX, *cropY, *cropWidth, *cropHeight, *outputFile)
	if err != nil {
		log.Fatalf("Error processing GIF: %v", err)
	}

	fmt.Println("GIF processed and saved to", *outputFile)
}

func processGIF(filePath string, cropX, cropY, cropWidth, cropHeight int, outputFile string) error {
	// Open the GIF file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening GIF file: %w", err)
	}
	defer file.Close()

	// Decode the GIF
	srcGif, err := gif.DecodeAll(file)
	if err != nil {
		return fmt.Errorf("error decoding GIF: %w", err)
	}

	// Crop the GIF
	croppedGif := cropGIF(srcGif, cropX, cropY, cropWidth, cropHeight)

	// Save the cropped GIF
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	err = gif.EncodeAll(outFile, croppedGif)
	if err != nil {
		return fmt.Errorf("error encoding cropped GIF: %w", err)
	}

	return nil
}

func cropGIF(srcGif *gif.GIF, cropX, cropY, cropWidth, cropHeight int) *gif.GIF {
	croppedGif := &gif.GIF{
		Image:     make([]*image.Paletted, len(srcGif.Image)),
		Delay:     srcGif.Delay,
		LoopCount: srcGif.LoopCount,
		Disposal:  srcGif.Disposal,
		Config:    image.Config{ColorModel: srcGif.Config.ColorModel, Width: cropWidth, Height: cropHeight},
		BackgroundIndex: srcGif.BackgroundIndex,
	}

	for i, frame := range srcGif.Image {
		// Create a new RGBA image for the cropped frame
		rgba := image.NewRGBA(image.Rect(0, 0, cropWidth, cropHeight))

		// Create a gift filter to crop the image
		g := gift.New(gift.Crop(image.Rect(cropX, cropY, cropX+cropWidth, cropY+cropHeight)))

		// Apply the filter to crop the frame
		g.Draw(rgba, frame)

		// Convert the RGBA image back to a Paletted image
		palettedImage := image.NewPaletted(image.Rect(0, 0, cropWidth, cropHeight), frame.Palette)
		for x := 0; x < cropWidth; x++ {
			for y := 0; y < cropHeight; y++ {
				idx := frame.Palette.Index(rgba.At(x, y))
				palettedImage.SetColorIndex(x, y, uint8(idx))
			}
		}

		croppedGif.Image[i] = palettedImage
	}

	return croppedGif
}
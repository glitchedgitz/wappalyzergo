package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
)

var mainFile wappalyzer.Fingerprints
var total = 0
var totalExtracted = 0

func main() {
	readMainFile()

	// list file from directory
	listFile("technologies")
	writeMainFile()

	fmt.Println("Total:", total)
	fmt.Println("Total Extracted:", totalExtracted)
}

func readMainFile() {
	// Open the file
	file, err := os.Open(`D:\sdks\go\src\github.com\glitchedgitz\wappalyzergo\fingerprints_data.json`)
	if err != nil {
		log.Fatal("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
		return
	}

	// Covnert to json
	err = json.Unmarshal(content, &mainFile)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
		return
	}
}

func writeMainFile() {
	// Open the file
	file, err := os.Create("fingerprints_data.json")
	if err != nil {
		log.Fatal("Error creating file:", err)
		return
	}
	defer file.Close()

	// Convert the JSON data to a byte slice
	data, err := json.MarshalIndent(mainFile, "", "    ")
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
		return
	}

	// Write the data to the file
	_, err = file.Write(data)
	if err != nil {
		log.Fatal("Error writing to file:", err)
		return
	}
}

func FindDominantColor(imagePath string) (color.RGBA, error) {
	total++
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("error opening image file: %s", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("error decoding image: %s", err)
	}

	// Create a context
	dc := gg.NewContextForImage(img)

	// Convert image to RGBA format
	rgbaImg := dc.Image().(*image.RGBA)

	// Count occurrences of each color
	colorCounts := make(map[color.RGBA]int)
	for y := 0; y < rgbaImg.Bounds().Max.Y; y++ {
		for x := 0; x < rgbaImg.Bounds().Max.X; x++ {
			color := rgbaImg.RGBAAt(x, y)
			colorCounts[color]++
		}
	}

	// Find the color with the highest occurrence
	var dominantColor color.RGBA
	maxCount := 0
	for col, count := range colorCounts {
		if count > maxCount {
			dominantColor = col
			maxCount = count
		}
	}
	totalExtracted++
	return dominantColor, nil
}

func listFile(directory string) {
	// Open the directory
	dir, err := os.Open(directory)
	if err != nil {
		log.Fatal("Error opening directory:", err)
		return
	}
	defer dir.Close()

	// Read all files in the directory
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal("Error reading directory contents:", err)
		return
	}

	// Loop through each file in the directory
	for _, fileInfo := range fileInfos {
		// Check if it's a regular file (not a directory)
		if fileInfo.Mode().IsRegular() {
			// Print the file name
			fmt.Println(fileInfo.Name())

			// Get the full path of the file
			fullPath := filepath.Join(directory, fileInfo.Name())

			// Check the file extension
			ext := filepath.Ext(fullPath)
			switch ext {
			case ".svg":
				// Convert SVG to PNG
				err := convertSVGToPNG(fullPath)
				if err != nil {
					log.Printf("Error converting SVG to PNG: %v", err)
					continue
				}
				fullPath = fullPath + ".png"
			case ".jpg", ".jpeg":
				// Convert JPEG to PNG
				err := convertJPEGToPNG(fullPath)
				if err != nil {
					log.Printf("Error converting JPEG to PNG: %v", err)
					continue
				}
				fullPath = fullPath + ".png"
			}

			// Read the json from file
			file, err := os.Open(fullPath)
			if err != nil {
				log.Fatal("Error opening file:", err)
				return
			}
			defer file.Close()

			// Decode the image
			img, _, err := image.Decode(file)
			if err != nil {
				log.Fatal("Error decoding image:", err)
				return
			}

			// Create a rasterizer
			rast := rasterizer.NewCanvas(1)

			// Rasterize the image
			_, err = rast.DrawImage(0, 0, img)
			if err != nil {
				log.Fatal("Error rasterizing image:", err)
				return
			}

			// Convert the rasterizer to a GG image
			ggImg := rast.Image()

			// Convert image to RGBA format
			rgbaImg := ggImg.(*image.RGBA)

			// Count occurrences of each color
			colorCounts := make(map[color.RGBA]int)
			for y := 0; y < rgbaImg.Bounds().Max.Y; y++ {
				for x := 0; x < rgbaImg.Bounds().Max.X; x++ {
					color := rgbaImg.RGBAAt(x, y)
					colorCounts[color]++
				}
			}

			// Find the color with the highest occurrence
			var dominantColor color.RGBA
			maxCount := 0
			for col, count := range colorCounts {
				if count > maxCount {
					dominantColor = col
					maxCount = count
				}
			}

			// Print the dominant color
			fmt.Printf("Dominant color of %s: %s\n", fileInfo.Name(), ColorToHex(dominantColor))
		}
	}
}

func convertSVGToPNG(filePath string) error {
	// Open the SVG file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the SVG file
	svgImg, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Create a canvas
	c := canvas.New(0, 0)

	// Rasterize the SVG
	err = rasterizer.Draw(c, svgImg)
	if err != nil {
		return err
	}

	// Save the rasterized image as PNG
	pngFile, err := os.Create(filePath + ".png")
	if err != nil {
		return err
	}
	defer pngFile.Close()

	err = canvas.Encode(pngFile, c, canvas.PNGWriter())
	if err != nil {
		return err
	}

	return nil
}

func convertJPEGToPNG(filePath string) error {
	// Open the JPEG file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the JPEG file
	jpegImg, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Create a new PNG file
	pngFile, err := os.Create(filePath + ".png")
	if err != nil {
		return err
	}
	defer pngFile.Close()

	// Encode the JPEG image as PNG
	err = png.Encode(pngFile, jpegImg)
	if err != nil {
		return err
	}

	return nil
}

func ColorToHex(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

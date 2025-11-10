// ibconv v0.1
// Author: Bugra Coskun
// License: GPLv3

package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	_ "image/gif" // Import for GIF decoding

	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

const VERSION string = "v0.2"
var SEPARATOR string = "\n" + strings.Repeat("-",50) + "\n"

// We'll pass this config struct around to make handling flags easier.
type Config struct {
	InputFolder  string
	OutputFolder string
	TargetFormat string
	TargetWidth  uint
	TargetHeight uint
	Verbose      bool
	Help         bool
}

func main() {
	//fmt.Printf(SEPARATOR)
	fmt.Printf("\n[ ibconv %s ]\n",VERSION)
	fmt.Println("Author: Bugra Coskun")
	fmt.Println("LICENSE: GPLv3")
	fmt.Println("")
	//fmt.Printf(SEPARATOR)

	inFolder := flag.String("i", "./source", "Path of input folder")
	outFolder := flag.String("o", "./sink", "Path of output folder")
	resolution := flag.String("r", "280,180", "Image size W,H (e.g., '280,180')")
	format := flag.String("f", "jpg", "Output format (jpg or png)")
	verbose := flag.Bool("v", false, "Enable verbose output")
	help := flag.Bool("h", false, "Prints help output")

	flag.Parse()


	width, height, err := parseResolution(*resolution)
	if err != nil {
		log.Fatalf("Invalid resolution format. Must be W,H. Error: %v", err)
	}

	// 4. Validate the format
	targetFormat := strings.ToLower(*format)
	if targetFormat != "jpg" && targetFormat != "png" {
		log.Fatalf("Invalid format. Must be 'jpg' or 'png'.")
	}

	// 5. Create the config
	config := Config{
		InputFolder:  *inFolder,
		OutputFolder: *outFolder,
		TargetFormat: targetFormat,
		TargetWidth:  width,
		TargetHeight: height,
		Verbose:      *verbose,
		Help:         *help,
	}

	if config.Verbose {
		fmt.Println("Starting ibconv...")
		fmt.Printf("Input Folder: %s\n", config.InputFolder)
		fmt.Printf("Output Folder: %s\n", config.OutputFolder)
		fmt.Printf("Target Size: %dx%d\n", config.TargetWidth, config.TargetHeight)
		fmt.Printf("Target Format: %s\n", config.TargetFormat)
	}

	if config.Help {
		printHelp()
	} else {
		runConversion(config)
	}
}

func printHelp() {
	fmt.Println(`Usage: ibconv [OPTIONS]...

	A simple bulk image converter.

	The program converts all images (jpg, png, gif) from an input
	folder, resizes them, and saves them to an output folder in
	the specified format.

	Default behavior (no arguments):
	  Converts images from ./source to ./sink at 280x180 resolution
	  in 'jpg' format.

	Options:
	  -i <folder>     Path to the input source folder.
					  (default: ./source)

	  -o <folder>     Path to the output sink folder.
					  (default: ./sink)

	  -r <W,H>        Target resolution (Width,Height) to resize images to.
					  (e.g., "800,600")
					  (default: 280,180)

	  -f <format>     Target output format. Can be 'jpg' or 'png'.
					  (default: jpg)

	  -v              Enable verbose output, showing processing details.

	  -h              Print this help and usage message.
	`)
}

func runConversion(config Config) {

	// Check if the input folder exists and is a directory
	info, err := os.Stat(config.InputFolder)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Error: Input folder '%s' does not exist.",
				config.InputFolder)

		} else {
			log.Fatalf("Error checking input folder: %v", err)
		}
	}

	if !info.IsDir() {
		// Path exists but is a file, not a directory
		log.Fatalf("Error: Input path '%s' is a file, not a folder.",
			config.InputFolder)
	}

	if err := os.MkdirAll(config.OutputFolder, 0755); err != nil {
		// Create the output directory if it doesn't exist
		log.Fatalf("Failed to create output directory: %v", err)
	}

	if config.Verbose {
		fmt.Printf("Scanning folder: %s\n", config.InputFolder)
	}

	// Walk through the input directory
	err = filepath.WalkDir(config.InputFolder,
		func(path string, d os.DirEntry, err error) error {

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // Skip directories
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {

			// Determine the new output path
			relPath, err := filepath.Rel(config.InputFolder, path)
			if err != nil {
				log.Printf("Warning: Could not get relative path for %s: %v", path, err)
				return nil
			}

			// Change the extension to the target format
			baseName := relPath[0 : len(relPath)-len(ext)]
			newExt := "." + config.TargetFormat
			newPath := filepath.Join(config.OutputFolder, baseName+newExt)
			
			// Ensure the new subdirectory (if any) exists
			if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
				log.Printf("Warning: Could not create sub-directory for %s: %v", newPath, err)
				return nil
			}

			// Process the image
			err = processImage(path, newPath, config)
			if err != nil {
				log.Printf("Warning: Failed to process %s: %v", path, err)
			} else if config.Verbose {
				fmt.Printf("Converted: %s -> %s\n", path, newPath)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}

	fmt.Println("Conversion complete.")
}

// processImage now takes the config to know what size/format to use
func processImage(inputPath, outputPath string, config Config) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("could not decode image (format %s): %w", format, err)
	}

	// Resize using the dimensions from config
	resizedImg := resize.Resize(config.TargetWidth, config.TargetHeight, img, resize.Lanczos3)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer outFile.Close()

	// Use a switch to encode to the correct target format
	switch config.TargetFormat {
	case "jpg":
		return jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(outFile, resizedImg)
	}
	
	return fmt.Errorf("unknown target format: %s", config.TargetFormat)
}

// parseResolution is a helper function to handle the "W,H" string
func parseResolution(res string) (width, height uint, err error) {
	parts := strings.Split(res, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format, expected 'W,H'")
	}

	w, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %w", err)
	}

	h, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %w", err)
	}

	if w <= 0 || h <= 0 {
		return 0, 0, fmt.Errorf("width and height must be positive")
	}

	return uint(w), uint(h), nil
}

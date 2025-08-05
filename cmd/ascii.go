package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"scrawl/utils"
	"strings"

	"github.com/spf13/cobra"
)

var (
	width   int 
	height  int
	invert  bool
	color   bool
	outFile string
)

var asciiChars = " .:-=+*#%@"

func Ascii(cmd *cobra.Command, args []string) error {
	path := args[0]

	file, err := os.Open(path)
	if err != nil {
		utils.ExitWithError("Error opening file", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		utils.ExitWithError("Error decoding image", err)
	}

	bounds := img.Bounds()
	if height == 0 {
		aspectRatio := float64(bounds.Dy()) / float64(bounds.Dx())
		height = int(float64(width) * aspectRatio * 0.5)
	}
	resized := utils.Resize(img, width, height)

	var output strings.Builder

	bounds = resized.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			r8, g8, b8 := float64(r>>8), float64(g>>8), float64(b>>8)

			luminance := 0.2126*r8 + 0.7152*g8 + 0.0722*b8
			index := int(luminance * float64(len(asciiChars)-1) / 255)
			if invert {
				index = len(asciiChars) - 1 - index
			}

			char := string(asciiChars[index])
			if color {
				output.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s", int(r8), int(g8), int(b8), char))
			} else {
				output.WriteString(char)
			}
		}
		output.WriteString("\n")
	}

	if color {
		output.WriteString("\x1b[0m")
	}

	if outFile != "" {
		err = os.WriteFile(outFile, []byte(output.String()), 0644)
		if err != nil {
			utils.ExitWithError("Failed writing to output file", err)
		}
	} else {
		fmt.Print(output.String())
	}

	return nil

}

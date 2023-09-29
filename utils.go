package main

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"os"
)

func LoadFont(path string) *truetype.Font {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	return f
}

func SizedFont(f *truetype.Font, points float64) font.Face {
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face
}

func loadImage(path string) image.Image {
	img, err := gg.LoadImage(path)
	if err != nil {
		panic(err)
	}
	return img
}

package main

import (
	"fmt"
	"path/filepath"

	"github.com/Ernyoke/Imger/imgio"
	"github.com/anthonynsimon/bild/transform"
	ocr "github.com/cam-per/digits-ocr"
)

func main() {
	//ocr.Fit(filepath.Join("..", "dataset"))
	ocr.Load(filepath.Join("..", "ocr.gob"))

	img, err := imgio.ImreadRGBA(filepath.Join("..", "dataset", "90.png"))
	if err != nil {
		return
	}

	img = transform.Resize(img, 30, 30, transform.NearestNeighbor)
	digit, err := ocr.Digit(img)
	fmt.Print(digit, err)
}

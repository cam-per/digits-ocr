package ocr

import (
	"image/color"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	"github.com/anthonynsimon/bild/adjust"

	"github.com/anthonynsimon/bild/transform"

	"github.com/Ernyoke/Imger/imgio"

	"github.com/therfoo/therfoo/tensor"
)

type data struct {
	in  tensor.Vector
	out tensor.Vector
}

type trainer struct {
	items []data
}

func (t *trainer) Get(index int) (*[]tensor.Vector, *[]tensor.Vector) {
	data := t.items[index]
	return &[]tensor.Vector{data.in}, &[]tensor.Vector{data.out}
}

func (t *trainer) Len() int {
	return len(t.items)
}

func (t *trainer) imageToVector(file string) (vector tensor.Vector, err error) {
	img, err := imgio.ImreadRGBA(file)
	if err != nil {
		return
	}

	img = transform.Resize(img, 30, 30, transform.NearestNeighbor)
	//imgio.Imwrite(img, xid.New().String()+".png")
	white := color.RGBA{255, 255, 255, 255}
	adjust.Apply(img, func(c color.RGBA) color.RGBA {
		if c == white {
			vector.Append(0.0)
		} else {
			vector.Append(1.0)
		}
		return c
	})
	return
}

func (t *trainer) fetchFiles(dir string) {
	t.items = t.items[:0]

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		digit, err := strconv.Atoi(f.Name()[:1])
		if err != nil {
			log.Fatal(err)
		}

		in, err := t.imageToVector(filepath.Join(dir, f.Name()))
		if err != nil {
			log.Fatal(err)
		}
		out := make(tensor.Vector, 10)
		out[digit] = 1.0

		t.items = append(t.items, data{in, out})
	}
}

var generator *trainer

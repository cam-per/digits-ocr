package ocr

import (
	"errors"
	"image"
	"image/color"
	"log"
	"path/filepath"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/cam-per/therfoo/optimizers/sgd"
	"github.com/therfoo/therfoo/layers/dense"
	"github.com/therfoo/therfoo/model"
	"github.com/therfoo/therfoo/tensor"
)

var nn *model.Model

func init() {
	generator = new(trainer)

	nn = model.New(
		model.WithCategoricalAccuracy(),
		model.WithCrossEntropyLoss(),
		model.WithEpochs(15000),
		model.WithInputShape(tensor.Shape{900}),
		model.WithOptimizer(
			sgd.New(sgd.WithBatchSize(10), sgd.WithLearningRate(0.005)),
		),
		model.WithTrainingGenerator(generator),
		model.WithValidatingGenerator(generator),
		model.WithVerbosity(true),
	)

	nn.Add(16, dense.New(dense.WithReLU()))
	nn.Add(10, dense.New(dense.WithSigmoid()))

	nn.Compile()
}

func Fit(datasetDir string) {
	generator.fetchFiles(datasetDir)
	nn.Fit()
	nn.Save(filepath.Join("..", "ocr.gob"))
	result := *nn.Predict(&[]tensor.Vector{generator.items[1].in})
	i, m := result[0].Max()
	log.Print(i, m)
}

func Load(confFile string) error {
	return nn.Load(confFile)
}

func imageToVector(img image.Image) (vector tensor.Vector) {
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

func Digit(img image.Image) (int, error) {
	if img.Bounds().Dx() != 30 && img.Bounds().Dy() != 30 {
		return 0, errors.New("image size must be 30x30")
	}
	result := *nn.Predict(&[]tensor.Vector{imageToVector(img)})
	i, _ := result[0].Max()
	return i, nil
}

package main

import (
	"image"
	"math"
	"math/rand"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Mouse",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	spritesheet, err := loadPicture("assets/trees.png")
	if err != nil {
		panic(err)
	}

	frames := []pixel.Rect{}
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
			frames = append(frames, pixel.R(x, y, x+32, y+32))
		}
	}

	var (
		camPos       = pixel.ZV
		camZoom      = 1.0
		camZoomSpeed = 1.2
		mouseOffset  = pixel.ZR
		trees        []*pixel.Sprite
		matrices     []pixel.Matrix
	)

	for !win.Closed() {
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		if win.JustPressed(pixelgl.MouseButtonRight) {
			mouseOffset.Min = win.MousePosition()
		}
		if win.Pressed(pixelgl.MouseButtonRight) {
			mouseOffset.Max = win.MousePosition()
		} else {
			camPos.X -= mouseOffset.W() / camZoom
			camPos.Y -= mouseOffset.H() / camZoom

			mouseOffset = pixel.ZR
		}

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			tree := pixel.NewSprite(spritesheet, frames[rand.Intn(len(frames))])
			trees = append(trees, tree)

			matrix := pixel.IM.Scaled(pixel.ZV, 2).Moved(cam.Unproject(win.MousePosition()))
			matrices = append(matrices, matrix)
		}

		win.Clear(colornames.Skyblue)

		for i, tree := range trees {
			tree.Draw(win, matrices[i])
		}

		win.Update()
	}
}

func loadPicture(filename string) (pixel.Picture, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

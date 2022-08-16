package main

import (
	"image"
	"math"
	"math/rand"
	"os"
	"time"

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
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		trees        []*pixel.Sprite
		matrices     []pixel.Matrix
	)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		relativeCamSpeed := camSpeed / camZoom

		if win.Pressed(pixelgl.KeyW) {
			camPos.Y += relativeCamSpeed * dt
		}
		if win.Pressed(pixelgl.KeyA) {
			camPos.X -= relativeCamSpeed * dt
		}
		if win.Pressed(pixelgl.KeyS) {
			camPos.Y -= relativeCamSpeed * dt
		}
		if win.Pressed(pixelgl.KeyD) {
			camPos.X += relativeCamSpeed * dt
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

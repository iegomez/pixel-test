package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
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

func run() {

	screenWidth := 1024
	screenHeight := 768

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, float64(screenWidth), float64(screenHeight)),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	spritesheet, err := loadPicture("assets/map_tiles.png")
	if err != nil {
		panic(err)
	}

	var matrices []pixel.Matrix
	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	var tiles []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 34 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 34 {
			tiles = append(tiles, pixel.R(x+1, y+1, x+32, y+32))
		}
	}

	tile := pixel.NewSprite(spritesheet, tiles[int(spritesheet.Bounds().Max.X/32)+1])

	xOffset := int(screenHeight / 32)
	yOffset := screenHeight * 2

	for i := -1 * xOffset; i < xOffset; i++ {
		for j := -1 * yOffset; j < yOffset; j++ {
			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(pixel.Vec{float64(i * 32), float64(j * 32)}))
		}
	}

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.Pressed(pixelgl.KeyLeft) {
			//camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			//camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		batch.Clear()
		for i := 0; i < 32; i++ {
			for j := 0; j < 64; j++ {
				tile.Draw(batch, matrices[i*j+j])
			}
		}
		batch.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}

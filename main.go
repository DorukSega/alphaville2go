package main

//lint:file-ignore U1000 Ignore all unused code, it's generated

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screen_width  int32 = 1000
	screen_height int32 = 1000
	N_ROWS        int32 = 256 // game width instead of screen
	N_COLS        int32 = 256
)

var fullscreen = false
var vsync = true

var memory = make([][]color.RGBA, int(N_ROWS))

// var old_memory = make([][]color.RGBA, int(N_ROWS))

var cameraX, cameraY float64 = 0, 0
var move_speed float64 = 2
var diagonal_speed = move_speed / 1

var isLeftArrowPressed = false
var isRightArrowPressed = false
var isUpArrowPressed = false
var isDownArrowPressed = false

func main() {
	init_mdimension(&memory, int(N_COLS))
	mt_sliceFill(&memory, color.RGBA{0, 0, 0, 255})
	// init_mdimension(&old_memory, int(N_COLS))
	// mt_sliceFill(&old_memory, color.RGBA{0, 0, 0, 255})

	tilemap1 := load_tilemap("./tilemap.bin")
	tileset1 := load_tileset("./tileset.png", 16)

	var flags uint32 = 0
	if fullscreen {
		flags |= rl.FlagFullscreenMode
	}
	if vsync {
		flags |= rl.FlagVsyncHint
	}

	rl.SetConfigFlags(flags)
	rl.InitWindow(screen_width, screen_height, "Alphaville")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		update_camera_position()

		draw_map(&tilemap1, &tileset1)
		rl.BeginDrawing()
		rl.ClearBackground(color.RGBA{17, 17, 17, 255})
		// for x := range memory {
		// 	for y := range memory[x] {
		// 		memory[x][y] = color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256))}
		// 	}
		// }

		draw_memory(pixel_ratio(screen_width))
		rl.EndDrawing()
	}
}

func pixel_ratio(width int32) int32 {
	return int32(math.Ceil(float64(width) / float64(N_ROWS)))
}

type tileset struct {
	image_set  [][][]color.RGBA
	tile_scale int32
}

func load_tileset(file_name string, tile_scale int32) tileset {
	// Open the file for reading
	file, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("image: %v\n", file_name)
		panic(err)
	}

	i_width, i_height := img.Bounds().Dx(), img.Bounds().Dy()
	t_x, t_y := i_width/int(tile_scale), i_height/int(tile_scale)
	tile_count := t_x * t_y
	//fmt.Printf("%v %v\n", t_x, t_y)

	var image_set = make([][][]color.RGBA, tile_count)
	for x := range image_set {
		image_set[x] = make([][]color.RGBA, tile_scale)
		for y := range image_set[x] {
			image_set[x][y] = make([]color.RGBA, tile_scale)
		}
	}

	for id := range image_set {
		for x := range image_set[id] {
			for y := range image_set[id][x] {
				itr_x, itr_y := id%t_x, id/t_x
				row := x + itr_x*int(tile_scale)
				col := y + itr_y*int(tile_scale)

				r, g, b, a := img.At(int(row), int(col)).RGBA()
				pixel := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
				//fmt.Printf("(%v %v %v) (%v %v) %v\n", id, x, y, row, col, pixel)
				image_set[id][x][y] = pixel
			}
		}
	}
	return tileset{image_set, tile_scale}
}

type tilemap struct {
	tmap   []byte
	width  int32
	height int32
}

func load_tilemap(file_path string) tilemap {
	file, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	width, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	height, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	var size = int(width) * int(height) // converting because this can be bigger than uint8
	var tmap = make([]byte, size)
	i := 0
	for {
		// Read a byte (uint8) from the file
		byteValue, err := reader.ReadByte()
		if err != nil {
			break // Reached the end of the file
		}
		tmap[i] = byteValue
		i++
	}
	return tilemap{tmap, int32(width), int32(height)}
}

func draw_map(tmap *tilemap, tset *tileset) {
	lenght := tmap.width * tset.tile_scale
	if cameraX == 0 && cameraY == 0 {
		cameraX = float64(lenght)/2 - float64(N_ROWS)/2
		cameraY = float64(lenght)/2 - float64(N_COLS)/2
	}

	for x := range memory {
		for y := range memory[x] {
			xI := x + int(math.Ceil(cameraX))
			yI := y + int(math.Ceil(cameraY))

			if xI < 0 || xI > int(lenght)-1 || yI < 0 || yI > int(lenght)-1 {
				memory[x][y] = color.RGBA{0, 0, 0, 255}
			} else {
				map_x, map_y := xI/int(tset.tile_scale), yI/int(tset.tile_scale)
				map_id := map_x + map_y*int(tmap.width)
				tileset_index := tmap.tmap[map_id]
				memory[x][y] = tset.image_set[tileset_index][xI%int(tset.tile_scale)][yI%int(tset.tile_scale)]
			}
		}
	}
}

func draw_memory(p_ratio int32) {
	for ix := range memory {
		for iy, pixel := range memory[ix] {
			//if old_memory[ix][iy] != pixel {
			rl.DrawRectangle(int32(ix)*p_ratio, int32(iy)*p_ratio, p_ratio, p_ratio, pixel)
			//}
		}
	}
}

func update_camera_position() {
	if rl.IsKeyDown(rl.KeyLeft) && rl.IsKeyDown(rl.KeyUp) {
		cameraX -= diagonal_speed
		cameraY -= diagonal_speed
	} else if rl.IsKeyDown(rl.KeyRight) && rl.IsKeyDown(rl.KeyUp) {
		cameraX += diagonal_speed
		cameraY -= diagonal_speed
	} else if rl.IsKeyDown(rl.KeyLeft) && rl.IsKeyDown(rl.KeyDown) {
		cameraX -= diagonal_speed
		cameraY += diagonal_speed
	} else if rl.IsKeyDown(rl.KeyRight) && rl.IsKeyDown(rl.KeyDown) {
		cameraX += diagonal_speed
		cameraY += diagonal_speed
	} else if rl.IsKeyDown(rl.KeyLeft) {
		cameraX -= move_speed
	} else if rl.IsKeyDown(rl.KeyRight) {
		cameraX += move_speed
	} else if rl.IsKeyDown(rl.KeyUp) {
		cameraY -= move_speed
	} else if rl.IsKeyDown(rl.KeyDown) {
		cameraY += move_speed
	}
}

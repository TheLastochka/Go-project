package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"time"
)

var (
	xCord int
	yCord int
	wWind int
	hWind int
)

type Img struct {
	w   int
	h   int
	rgb [][]color.Color // x,y
}
type Region struct {
	x0 int
	y0 int
	w  int
	h  int
}

func loadImg(path string) (img Img) {
	infile, err := os.Open(path)
	if err != nil {
		panic("can't open " + path)
	}
	defer infile.Close()

	src, err := png.Decode(infile)
	if err != nil {
		panic("can't decode " + path)
	}

	bounds := src.Bounds()
	img.w, img.h = bounds.Max.X, bounds.Max.Y

	img.rgb = make([][]color.Color, img.w)
	for x := 0; x < img.w; x++ {
		col := make([]color.Color, img.h)
		for y := 0; y < img.h; y++ {
			col[y] = src.At(x, y)
		}
		img.rgb[x] = col
	}
	return img
}
func saveImg(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
func getScreenshot(region ...Region) (screen Img) {
	reg := Region{xCord, yCord, wWind, hWind}
	if len(region) != 0 {
		reg = region[0]
	}
	src, _ := screenshot.Capture(reg.x0, reg.y0, reg.w, reg.h)
	//save(src, "all.png")
	bounds := src.Bounds()
	screen.w, screen.h = bounds.Max.X, bounds.Max.Y
	screen.rgb = make([][]color.Color, screen.w)
	for x := 0; x < screen.w; x++ {
		col := make([]color.Color, screen.h)
		for y := 0; y < screen.h; y++ {
			col[y] = src.At(x, y)
		}
		screen.rgb[x] = col
	}
	return screen
}
func findImage(path string, region ...Region) [2]int {
	reg := Region{xCord, yCord, wWind, hWind}
	if len(region) != 0 {
		reg = region[0]
	}
	img := loadImg(path)
	screen := getScreenshot(reg)
	var res [2]int
	for bigCol := 0; bigCol < screen.w-img.w; bigCol++ {
		//fmt.Println(bigRow)
		for bigRow := 0; bigRow < screen.h-img.h; bigRow++ {
			// calc
			err := false
			for x := 0; !err && x < img.w; x++ {
				for y := 0; !err && y < img.h; y++ {
					iR, iG, iB, _ := img.rgb[x][y].RGBA()
					sR, sG, sB, _ := screen.rgb[bigCol+x][bigRow+y].RGBA()
					if iR != sR || iG != sG || iB != sB {
						err = true
					}
				}
			}
			if !err {
				//fmt.Println("found!", bigRow, bigCol)
				res[0] = xCord + bigCol
				res[1] = yCord + bigRow
				return res
			}
		}
	}
	res[0] = -1
	res[1] = -1
	return res
}
func findPixel(r uint8, g uint8, b uint8, region ...Region) [2]int {
	reg := Region{xCord, yCord, wWind, hWind}
	if len(region) != 0 {
		reg = region[0]
	}
	screen := getScreenshot(reg)
	clr := color.RGBA{R: r, G: g, B: b, A: 255}
	for x := 0; x < screen.w; x++ {
		for y := 0; y < screen.h; y++ {
			if screen.rgb[x][y] == clr {
				return [2]int{x, y}
			}
		}
	}
	return [2]int{-1, -1}
}

func main() {
	sleepAfter := map[string]time.Duration{
		"restart":     60,
		"vkIcon":      12,
		"services":    14,
		"appIcon":     8,
		"watchBTN":    7,
		"siteOpenBtn": 7,
		"backArrow":   3,
		"closeAd":     1,
	}
	start := time.Now()
	//time.Sleep(2 * time.Second)
	//fpid, _ := robotgo.FindIds("HD-Player.exe")
	//
	//xWindow, yWindow, wWindow, hWindow := robotgo.GetBounds(fpid[0])
	//xCord = xWindow
	//yCord = yWindow-33
	//wWind = wWindow+33
	//hWind = hWindow+33
	xCord = 0
	yCord = 0
	wWind = 700
	hWind = 1080

	//cmd := exec.Command("calc")
	//cmd.Run()
	time.Sleep(2 * time.Second)
	com := "nircmd exec show \"calc\""
	fmt.Println(com)
	cmd := exec.Command(com)
	cmd.Run()

	//fpid, _ := robotgo.FindIds("HD-Player.exe")
	//Показ границ окна
	//x, y, w, h := robotgo.GetBounds(fpid[0])
	//fmt.Println("GetBounds is: ", x, y, w, h)
	//robotgo.Move(xCord,yCord)
	//time.Sleep(1* time.Second)
	//robotgo.Move(xCord+wWind,yCord)
	//time.Sleep(1* time.Second)
	//robotgo.Move(xCord+wWind,yCord+hWind)
	//time.Sleep(1* time.Second)
	//robotgo.Move(xCord,yCord+hWind)
	//time.Sleep(1* time.Second)
	//fmt.Println(fpid)

	if false {
		mas := findImage("./vkIcon.png")
		fmt.Println(mas)
		robotgo.Move(mas[0], mas[1])
	}

	fmt.Println(findPixel(243, 12, 67))

	end := time.Now()
	fmt.Println("time:", end.Sub(start))

	if false {
		cords := findImage("./vkIcon.png")
		robotgo.MoveClick(cords[0], cords[1])
		time.Sleep(sleepAfter["vkIcon"] * time.Second)

		cords = findImage("./services.png")
		robotgo.MoveClick(cords[0], cords[1])
		time.Sleep(sleepAfter["services"] * time.Second)

		cords = findImage("./appIcon.png")
		robotgo.MoveClick(cords[0], cords[1])
		time.Sleep(sleepAfter["appIcon"] * time.Second)

		//cords = findPixel(75, 179, 75)
		//robotgo.MoveClick(cords[0], cords[1])
		//time.Sleep(sleepAfter["watchBTN"] * time.Second)
		//
		//cords = findPixelInRegion(75, 179, 75, xCord, yCord+hWind/2, wWind, hWind/2)
		//robotgo.MoveClick(cords[0], cords[1])
		//time.Sleep(sleepAfter["siteOpenBtn"] * time.Second)

		cords = findImage("./backArrow.png", Region{xCord, yCord + hWind/2, wWind, hWind / 2})
		robotgo.MoveClick(cords[0], cords[1])
		time.Sleep(sleepAfter["backArrow"] * time.Second)

		cords = findImage("./closeAd.png", Region{xCord, yCord, wWind, hWind / 2})
		robotgo.MoveClick(cords[0], cords[1])
		time.Sleep(sleepAfter["closeAd"] * time.Second)
	}
}

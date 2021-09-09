package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"image"
	"image/color"
	"image/png"
	"os"
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
	rgb [][]color.Color
}

//fmt.Println(xWindow, yWindow, wWindow, hWindow)

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

	img.rgb = make([][]color.Color, img.h)
	for y := 0; y < img.h; y++ {
		row := make([]color.Color, img.w)
		for x := 0; x < img.w; x++ {
			row[x] = src.At(x, y)
		}
		img.rgb[y] = row
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
func getScreenshot() (screen Img) {
	start := time.Now()

	src, _ := screenshot.Capture(xCord, yCord, wWind, hWind)
	end := time.Now()

	fmt.Println("time:", end.Sub(start))
	//save(src, "all.png")
	bounds := src.Bounds()
	screen.w, screen.h = bounds.Max.X, bounds.Max.Y
	//fmt.Println(w, h)
	screen.rgb = make([][]color.Color, screen.h)
	for y := 0; y < screen.h; y++ {
		row := make([]color.Color, screen.w)
		for x := 0; x < screen.w; x++ {
			row[x] = src.At(x, y)
		}
		screen.rgb[y] = row
	}
	return screen
}
func findImage(path string) [2]int {
	img := loadImg(path)
	screen := getScreenshot()
	var res [2]int
	for bigRow := 0; bigRow < screen.h-img.h; bigRow++ {
		//fmt.Println(bigRow)
		for bigCol := 0; bigCol < screen.w-img.w; bigCol++ {
			// calc
			err := false
			for y := 0; !err && y < img.h; y++ {
				for x := 0; !err && x < img.w; x++ {
					iR, iG, iB, _ := img.rgb[y][x].RGBA()
					sR, sG, sB, _ := screen.rgb[bigRow+y][bigCol+x].RGBA()
					if iR != sR || iG != sG || iB != sB {
						err = true
					}
				}
			}
			if !err {
				//fmt.Println("found!", bigRow, bigCol)
				res[1] = bigRow + yCord
				res[0] = bigCol + xCord
				return res
			}
		}
	}
	res[0] = -1
	res[1] = -1
	return res
}
func main() {
	time.Sleep(2 * time.Second)
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

	start := time.Now()
	mas := findImage("./vkIcon.png")
	fmt.Println(mas)
	robotgo.Move(mas[0], mas[1])
	end := time.Now()

	fmt.Println("time:", end.Sub(start))
}

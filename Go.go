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

//fmt.Println(xWindow, yWindow, wWindow, hWindow)

func loadImg(path string) (w int, h int, rgb [][]color.Color) {
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
	w, h = bounds.Max.X, bounds.Max.Y

	rgb = make([][]color.Color, h)
	for y := 0; y < h; y++ {
		row := make([]color.Color, w)
		for x := 0; x < w; x++ {
			//color := src.At(x, y)
			//r, g, b, _ := color.RGBA()
			row[x] = src.At(x, y)
		}
		rgb[y] = row
	}
	return w, h, rgb
}

func saveImg(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
func getScreenshot() (int, int, [][]color.Color) {
	start := time.Now()

	src, _ := screenshot.Capture(xCord, yCord, wWind, hWind)
	end := time.Now()

	fmt.Println("time:", end.Sub(start))
	//save(src, "all.png")
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	//fmt.Println(w, h)
	rgb := make([][]color.Color, h)
	for y := 0; y < h; y++ {
		row := make([]color.Color, w)
		for x := 0; x < w; x++ {
			//color := src.At(x, y)
			//r, g, b, _ := color.RGBA()
			row[x] = src.At(x, y)
		}
		rgb[y] = row
	}
	return w, h, rgb
}
func findImage(path string) [2]int {
	wSmall, hSmall, small := loadImg(path)
	wBig, hBig, big := getScreenshot()
	var res [2]int
	for bigRow := 0; bigRow < hBig-hSmall; bigRow++ {
		//fmt.Println(bigRow)
		for bigCol := 0; bigCol < wBig-wSmall; bigCol++ {
			// calc
			err := false
			for y := 0; !err && y < hSmall; y++ {
				for x := 0; !err && x < wSmall; x++ {
					if small[y][x] != big[bigRow+y][bigCol+x] {
						err = true
					} else {
						fmt.Println("x:", x, "\ty:", y, "\tsmall:", small[y][x], "\tbig:", big[bigRow+y][bigCol+x])
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
	yCord = 500
	wWind = 100
	hWind = 300

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

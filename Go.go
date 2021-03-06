package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"time"
	"unsafe"
)

var (
	xCord        int
	yCord        int
	wWind        int
	hWind        int
	sitesWatched int
	errorsCount  int
	lastAd       time.Time
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
type Config struct {
	SleepAfter map[string]time.Duration `json:"sleep_after"`
	Icons      map[string]string        `json:"icons"`
	Colors     map[string]color.RGBA    `json:"colors"`
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
func findPixel(clr color.RGBA, region ...Region) [2]int {
	reg := Region{xCord, yCord, wWind, hWind}
	if len(region) != 0 {
		reg = region[0]
	}
	screen := getScreenshot(reg)
	for x := 0; x < screen.w; x++ {
		for y := 0; y < screen.h; y++ {
			if screen.rgb[x][y] == clr {
				return [2]int{x, y}
			}
		}
	}
	return [2]int{-1, -1}
}
func openConfig(path string) (Config, error) {
	file, _ := os.Open(path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Config{}
	err := decoder.Decode(&conf)
	return conf, err
}
func getBase64Shot() string {
	src, _ := screenshot.Capture(0, 0, 1920, 1080)
	emptyBuff := bytes.NewBuffer(nil)                         //Open up a new empty buff
	jpeg.Encode(emptyBuff, src, nil)                          //img is written to buff
	dist := make([]byte, int32(float32(emptyBuff.Len())*1.4)) //Open up storage space
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())        //buff is converted to base64
	index := bytes.IndexByte(dist, 0)                         //Note here that because of the fixed-length array applied, the part that has not been filled needs to be removed, and the output may be wrong.
	baseImage := dist[0:index]
	return *(*string)(unsafe.Pointer(&baseImage))
}
func addSiteView() {
	sitesWatched++
}
func addError() {
	errorsCount++
}
func log(s string, m string) {
	/* ?????????????? ???? ????, ?????? ?????????? ?????????????????? ?????????????????????? ???????? ????????
	?? ???????????????? ?? ???????? ???????????? ???????????? ????????????????
	 ?????????? ?????????????? ??????????-???? ?????????? ????????????
	 ?????????? ???????? ?????????????????? ???????? ????????
	 ?????????? ?????????? ??????????-???? ???????????????????? golang ?????? ??????????????
	 ?????????? ?????? ??????????-???? ???????? ?????? ????????????*/
}

func main() {
	var cords [2]int
	start := time.Now()
	lastAd = time.Now()
	conf, err := openConfig("./config.json")
	if err != nil {
		panic("can't open config")
	}

	if false {
		// ???????????????? hd-player ?? path ???????? ?????? ??????????????
		exec.Command("correct_wind.bat").Run()
		fpid, _ := robotgo.FindIds("HD-Player.exe")
		time.Sleep(2 * time.Second)
		xWindow, yWindow, wWindow, hWindow := robotgo.GetBounds(fpid[0])
		xCord = xWindow
		yCord = yWindow - 33
		wWind = wWindow + 33
		hWind = hWindow + 33
	}

	xCord = 0
	yCord = 0
	wWind = 700
	hWind = 1080

	//fpid, _ := robotgo.FindIds("HD-Player.exe")
	//?????????? ???????????? ????????
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
	end := time.Now()
	fmt.Println("time:", end.Sub(start))

	//TODO: stats.json ???????????? ?? ???????????? sitesC ?? errorsC
	sitesWatched = 0
	errorsCount = 0

	cords = findImage(conf.Icons["vkIcon"])
	if cords[0] != -1 {
		robotgo.MoveClick(cords[0], cords[1])
		log("info", "???????????????? ????")
		time.Sleep(conf.SleepAfter["vkIcon"] * time.Second)

		cords = findImage(conf.Icons["services"])
		if cords[0] != -1 {
			robotgo.MoveClick(cords[0], cords[1])
			log("info", "???????????????? ????????????????")
			time.Sleep(conf.SleepAfter["services"] * time.Second)

			cords = findImage(conf.Icons["appIcon"])
			if cords[0] != -1 {
				robotgo.MoveClick(cords[0], cords[1])
				time.Sleep(conf.SleepAfter["appIcon"] * time.Second)
			} else {
				log("error", "???????????????????? ???? ??????????????")
			}
		} else {
			log("error", "?????????????? ???? ??????????????")
		}

	} else {
		log("error", "???? ???? ????????????")
	}

	//TODO: ???? ?????????? ???????? ????????, ?? ?????????? ?????? ??????????????

	if false {

		//cords = findPixel(75, 179, 75)
		//robotgo.MoveClick(cords[0], cords[1])
		//time.Sleep(sleepAfter["watchBTN"] * time.Second)
		//
		//cords = findPixelInRegion(75, 179, 75, xCord, yCord+hWind/2, wWind, hWind/2)
		//robotgo.MoveClick(cords[0], cords[1])
		//time.Sleep(sleepAfter["siteOpenBtn"] * time.Second)

		cords = findImage(conf.Icons["backArrow"], Region{xCord, yCord + hWind/2, wWind, hWind / 2})
		robotgo.MoveClick(cords[0], cords[1])
		time.Sleep(conf.SleepAfter["backArrow"] * time.Second)

		cords = findImage(conf.Icons["closeAd"], Region{xCord, yCord, wWind, hWind / 2})
		robotgo.MoveClick(cords[0], cords[1])
		time.Sleep(conf.SleepAfter["closeAd"] * time.Second)
	}
}

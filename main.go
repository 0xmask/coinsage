package main

import (
	"fmt"
	"github.com/0xmask/itools/ihelp"
	"github.com/0xmask/itools/ilog"
	"github.com/go-vgo/robotgo"
	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cast"
	"github.com/vcaesar/imgo"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Conf struct {
	ScreenWidth  int
	ScreenHeight int
}

var (
	//BTC行情图片
	btcPng = "./btc.png"
	//屏幕宽度
	ScreenWidth = 1440
	//屏幕高度
	ScreenHeight = 900
	//按钮x
	BtnX = ScreenWidth / 2
	//按钮y
	BtnY = ScreenHeight - 30
	//截图开始位置x
	CaptureX = ScreenWidth/2 + 146
	//截图开始位置y
	CaptureY = ScreenHeight - 60
	//截图宽度
	CaptureW = 110
	//截图高度
	CaptureH = 50
	//刷新按钮x
	RefreshX = 94
	//刷新按钮y
	RefreshY = 87

	t1Map = map[int]int{1: 0, 6: 0, 11: 0, 16: 0, 21: 0, 26: 0, 31: 0, 36: 0, 41: 0, 46: 0, 51: 0, 56: 0}

	//临时value
	btcTmp = float64(0)
)

func init() {
	ilog.Init(false, ilog.DebugLevel, "")
}

func main() {
	t1 := time.NewTicker(time.Second * 1)
	defer t1.Stop()

	for {
		select {
		case <-t1.C:
			m := cast.ToInt(carbon.Now().Format("i"))
			if _, ok := t1Map[m]; ok {
				s := cast.ToInt(carbon.Now().Format("s"))
				if s%5 == 0 {
					if s <= 45 {
						doit()
					} else {
						refresh()
					}
				}
			}
		}
	}
}

// 移动到刷新按钮并刷新
func refresh() {
	robotgo.MilliSleep(100)

	robotgo.Move(RefreshX, RefreshY)

	robotgo.MilliSleep(100)

	robotgo.Click()

	log.Println("刷新完成")
}

// 保持一个终端和一个chrome
// https://huggingface.co/spaces/siddreddy/CoinSage
func doit() {
	defer ihelp.ErrCatch()

	//打开chrome
	err := robotgo.ActiveName("chrome")
	if err != nil {
		ilog.Logger.Error(err)
	}

	//如果没有点击按钮则滑到底部并移动鼠标到按钮上
	robotgo.MilliSleep(100)

	robotgo.ScrollDir(ScreenHeight, "down")

	robotgo.Move(BtnX, BtnY)

	robotgo.Click()

	robotgo.MilliSleep(100)

	robotgo.ScrollDir(ScreenHeight, "down")

	//log.Println("点击按钮后滚到底部")

	//截图
	img := robotgo.CaptureImg(CaptureX, CaptureY, CaptureW, CaptureH)
	if err = imgo.Save(btcPng, img); err != nil {
		ilog.Logger.Error(err)
		return
	}

	//读取btc行情
	txt, err := robotgo.GetText(btcPng)
	if err != nil {
		ilog.Logger.Error(err)
		return
	}

	//todo 保存数据
	txt = strings.TrimSpace(txt)
	txt = strings.Trim(txt, "\n")
	btc, _ := strconv.ParseFloat(txt, 64)
	if btc > 0 && btc != btcTmp {
		saveBtc(btc)
	}
}

func saveBtc(btc float64) {
	yesterday := carbon.CreateFromTimestamp(time.Now().Unix() - 5 + 86400).Format("Y-m-d H:i")
	value := fmt.Sprintf("%s >>> %f\n", yesterday, btc)

	f, err := os.OpenFile("./btc.log", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		ilog.Logger.Error(err)
		return
	}
	if _, err = f.Write([]byte(value)); err != nil {
		ilog.Logger.Error(err)
		return
	}
	btcTmp = btc
	fmt.Println(value)
}

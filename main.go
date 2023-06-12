package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"

	speednet "github.com/Gessar/network_tray_golang/speed_net"
	"github.com/Gessar/network_tray_golang/tray_net"
	"github.com/getlantern/systray"
	"github.com/gookit/ini/v2"
	"github.com/shirou/gopsutil/net"
)

func main() {
	err := ini.LoadExists("test.ini", "not-exist.ini")
	if err != nil {
		panic(err)
	}

	// load more, will override prev data by key
	// 	err = ini.LoadStrings(`
	// age = 100
	// [sec1]
	// newK = newVal
	// some = change val
	// `)
	age := ini.Int("age")           //delme
	age1 := ini.Int("age1")         //delme
	test := ini.String("sec1.key")  //delme
	base := ini.String("sec1.base") //delme
	fmt.Println(age)                //delme
	fmt.Println(age1)               //delme
	fmt.Println(test)               //delme
	fmt.Println(base)               //delme
	systray.Run(onReady, tray_net.OnExit)
}

func onReady() {
	icoFile, err := os.Open("./icon.ico")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer icoFile.Close()
	stat, err := icoFile.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(icoFile).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}
	systray.SetIcon(bs)
	//systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")

	downloadSpeed := make(chan string)

	go func() {
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case speed := <-downloadSpeed:
				systray.SetTooltip(speed)
			}
		}
	}()

	go func() {
		err := ini.LoadExists("test.ini", "not-exist.ini") //Загрузка ini файлов
		if err != nil {
			panic(err)
		}
		base := ini.String("sec1.base") //из секции sec1 значение из base
		var (
			lastBytesSent    uint64
			lastBytesRecv    uint64
			currentBytesSent uint64
			currentBytesRecv uint64
		)

		for {
			runtime.GC()
			ioCounters, err := net.IOCounters(false)
			if err != nil {
				panic(err)
			}

			for _, counter := range ioCounters {
				currentBytesSent += counter.BytesSent
				currentBytesRecv += counter.BytesRecv
			}

			if lastBytesSent != 0 && lastBytesRecv != 0 {
				//currentUploadSpeed := strconv.FormatFloat(float64(currentBytesSent-lastBytesSent)/1024.0, 'f', 2, 64)
				//currentDownloadSpeed := strconv.FormatFloat(float64(currentBytesRecv-lastBytesRecv)/1024.0, 'f', 2, 64)
				//speed := "Upload speed: " + currentUploadSpeed + " KB/s Download speed: " + currentDownloadSpeed + " KB/s"
				speed := formatSpeedText(currentBytesSent, lastBytesSent, currentBytesRecv, lastBytesRecv, base)
				fmt.Printf("\r%s", speed)
				downloadSpeed <- speed
			}

			lastBytesSent = currentBytesSent
			lastBytesRecv = currentBytesRecv
			currentBytesSent = 0
			currentBytesRecv = 0

			time.Sleep(time.Second)
			speednet.SetSpeed()
			fmt.Println(speednet.GetSpeed())
		}
	}()
}

func formatSpeedText(currentBytesSent uint64, lastBytesSent uint64, currentBytesRecv uint64, lastBytesRecv uint64, base string) string {
	uploadString, downloadString := "KB/s", "KB/s"
	var dec int
	if base == "bin" { //система исчисления бинарная или десятичная
		dec = 1024
	} else {
		dec = 1000
	}
	currentUploadSpeed := strconv.FormatFloat(float64(currentBytesSent-lastBytesSent)/float64(dec), 'f', 2, 64)   //Вычитаем из прошлого значения текущее значение отправленных данных
	currentDownloadSpeed := strconv.FormatFloat(float64(currentBytesRecv-lastBytesRecv)/float64(dec), 'f', 2, 64) //Вычитаем из прошлого значения текущее значение загруженных данных
	if float64(currentBytesSent-lastBytesSent)/float64(dec) >= float64(dec) {                                     //Если более чем dec (1000 или 1024, в зависимости от base), то конвертируем в MB\MiB
		currentUploadSpeed = strconv.FormatFloat(float64(currentBytesSent-lastBytesSent)/(float64(dec)*float64(dec)), 'f', 2, 64)
		if dec == 1024 {
			uploadString = "MiB/s"
		} else {
			uploadString = "MB/s"
		}
	}
	if float64(currentBytesRecv-lastBytesRecv)/float64(dec) >= float64(dec) {
		currentDownloadSpeed = strconv.FormatFloat(float64(currentBytesRecv-lastBytesRecv)/(float64(dec)*float64(dec)), 'f', 2, 64)
		if dec == 1024 {
			downloadString = "MiB/s"
		} else {
			downloadString = "MB/s"
		}
	}
	return "Upload speed: " + currentUploadSpeed + " " + uploadString + " Download speed: " + currentDownloadSpeed + " " + downloadString
}

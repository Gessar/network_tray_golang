package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/Gessar/network_tray_golang/tray_net"
	"github.com/getlantern/systray"
	"github.com/shirou/gopsutil/net"
)

func main() {
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
				speed := formatSpeedText(currentBytesSent, lastBytesSent, currentBytesRecv, lastBytesRecv)
				fmt.Printf("\r%s", speed)
				downloadSpeed <- speed
			}

			lastBytesSent = currentBytesSent
			lastBytesRecv = currentBytesRecv
			currentBytesSent = 0
			currentBytesRecv = 0

			time.Sleep(time.Second)
		}
	}()

}

func formatSpeedText(currentBytesSent uint64, lastBytesSent uint64, currentBytesRecv uint64, lastBytesRecv uint64) string {
	uploadString, downloadString := "KB/s", "KB/s"
	currentUploadSpeed := strconv.FormatFloat(float64(currentBytesSent-lastBytesSent)/1024.0, 'f', 2, 64)
	currentDownloadSpeed := strconv.FormatFloat(float64(currentBytesRecv-lastBytesRecv)/1024.0, 'f', 2, 64)
	if float64(currentBytesSent-lastBytesSent)/1024.0 >= 1024 {
		currentUploadSpeed = strconv.FormatFloat(float64(currentBytesSent-lastBytesSent)/1048576.0, 'f', 2, 64)
		uploadString = "MB/s"
	}
	if float64(currentBytesRecv-lastBytesRecv)/1024.0 >= 1024 {
		currentDownloadSpeed = strconv.FormatFloat(float64(currentBytesRecv-lastBytesRecv)/1048576.0, 'f', 2, 64)
		downloadString = "MB/s"
	}
	return "Upload speed: " + currentUploadSpeed + " " + uploadString + " Download speed: " + currentDownloadSpeed + " " + downloadString
}

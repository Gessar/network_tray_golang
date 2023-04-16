package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/shirou/gopsutil/net"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {

	systray.SetIcon(icon.Data)
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
				currentUploadSpeed := strconv.FormatFloat(float64(currentBytesSent-lastBytesSent)/1024.0, 'f', 2, 64)
				currentDownloadSpeed := strconv.FormatFloat(float64(currentBytesRecv-lastBytesRecv)/1024.0, 'f', 2, 64)
				speed := "Upload speed: " + currentUploadSpeed + " KB/s Download speed: " + currentDownloadSpeed + " KB/s"
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

func onExit() {
	fmt.Println("\nExiting...")
}

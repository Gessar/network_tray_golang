package speednet

import "github.com/shirou/gopsutil/net"

var (
	lastBytesSent1   uint64
	lastBytesSent2   uint64
	currentBytesSent uint64
	currentBytesRecv uint64
)

func SetSpeed() {

}

func GetSpeed(currentSent int64, currentRecv int64) int {
	ioCounters, err := net.IOCounters(false)
	if err != nil {
		panic(err)
	}
	for _, counter := range ioCounters {
		currentBytesSent += counter.BytesSent
		currentBytesRecv += counter.BytesRecv
	}
	return 1
}

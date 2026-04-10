package socketcantrols

import (
	"net"
	"os"
	"os/exec"
)

func DialConnection() *net.Conn {
	socket_conection, err := net.Dial("unix", "/tmp/mpvsocket")
	if err != nil {
		panic(err)
	}
	return &socket_conection
}

func MpvInit(url string) *os.Process {
	cmd := exec.Command("mpv", "--input-ipc-server=/tmp/mpvsocket", "--no-video", url)
	cmd.Start()
	proc, err := os.FindProcess(cmd.Process.Pid)
	if err != nil {
		panic(err.Error())
	}
	return proc
}

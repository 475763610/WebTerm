package main


import (
	"net/http"
	"os"
	"io"
	"io/ioutil"
	"encoding/json"
	"syscall"
	"unsafe"
	"github.com/gorilla/websocket"
)
type Action struct {
	Cmd string
	Data string
}
type windowSize struct {
	Rows uint16
	Cols uint16
	X    uint16
	Y    uint16
}
func PassAllOrigin(r *http.Request)  bool{
	return true
}
func BridgePtyToWs(pty *os.File,conn *websocket.Conn){
	defer func(){
		conn.Close()
	}()
	for {
		buf := make([]byte, 1024)
		read, err := pty.Read(buf)
		if err != nil {
			return
		}
		conn.WriteMessage(websocket.BinaryMessage, buf[:read])
	}
}
func BridgeWsToPty(PTY *os.File,conn *websocket.Conn){
	for {
		messageType, reader, err := conn.NextReader()
		if err != nil {
			return
		}
		if messageType == websocket.BinaryMessage{
			_,err:=io.Copy(PTY, reader)
			if err!=nil{
				return
			}
		}
		if messageType == websocket.TextMessage {
			buf,err:=ioutil.ReadAll(reader)
			if err!=nil{
				continue
			}

			action:=&Action{}
			err=json.Unmarshal(buf,action)
			if err!=nil{
				continue
			}

			switch action.Cmd {
			case "resize":
				newWindow:=windowSize{}
				err=json.Unmarshal([]byte(action.Data),&newWindow)
				if err!=nil{
					continue
				}
				syscall.Syscall(syscall.SYS_IOCTL, PTY.Fd(), syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(&newWindow)))
			}

		}
	}
}
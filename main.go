package main

import (
	"github.com/crabkun/MonitorKits"
	"github.com/crabkun/crab"
	"os/exec"
	"os"
	"github.com/kr/pty"
	"github.com/gorilla/websocket"
)

func GetPluginInfo() *MonitorKits.PluginInfo {

	t:=&MonitorKits.PluginInfo{}
	t.Name="WebTerm"
	t.DisplayName="网页虚拟终端"
	t.Author="Crabkun"
	t.Description="Term over Web"
	t.Version="1.0"
	return t

}

func GetPluginRoute() *MonitorKits.PluginRoute {
	t:=&MonitorKits.PluginRoute{}
	t.Add("GET","term","Term")
	return t
}

func LoadPlugin() error {
	return nil
}

func UnloadPlugin() error {
	return nil
}
func PluginIndex(ctx *crab.Context) {
	ctx.Redirect(302,"static/console.html")
}

var ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: PassAllOrigin,
}

func Term(ctx *crab.Context){
	conn, err := ws.Upgrade(ctx.RspWriter, ctx.Req, nil)
	if err != nil {
		panic("无法把连接切换为WebSocket")
		return
	}

	cmd := exec.Command("/bin/bash", "-l")
	cmd.Env = append(os.Environ(), "TERM=xterm-color")

	PTY, err := pty.Start(cmd)
	if err != nil {
		panic("无法启动pty")
		return
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Process.Wait()
		PTY.Close()
		conn.Close()
	}()
	go BridgePtyToWs(PTY,conn)
	BridgeWsToPty(PTY,conn)
}

package service

import (
	"encoding/json"
	"fmt"
	"k8sexec/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// upgrade:websocket: request origin not allowed by Upgrader.CheckOrigin
var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

type s struct {
	ws          *WsConnection
	resizeEvent chan remotecommand.TerminalSize
	mt          int
}
type xtermMessage struct {
	MsgType string `json:"type"`  // 类型:resize客户端调整终端, input客户端输入
	Input   string `json:"input"` // msgtype=input情况下使用
	Rows    uint16 `json:"rows"`  // msgtype=resize情况下使用
	Cols    uint16 `json:"cols"`  // msgtype=resize情况下使用
}

func (s *s) Read(p []byte) (int, error) {
	var xtermMsg xtermMessage
	msg, _ := s.ws.WsRead()
	json.Unmarshal(msg.Data, &xtermMsg)
	// s.mt = mt
	switch xtermMsg.MsgType {
	case "resize":
		s.resizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	case "input":
		size := len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
		return size, nil
	}
	return 0, nil
}

func (s *s) Write(p []byte) (int, error) {
	copyData := make([]byte, len(p))
	copy(copyData, p)
	size := len(p)
	err := s.ws.WsWrite(websocket.TextMessage, copyData)
	return size, err
}

func (handler *s) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.resizeEvent
	size = &ret
	return
}

func ExecShell(c *gin.Context) {
	clientset := utils.K8sInit()
	ws, err := wsInit(c)
	podName := "axe-ui"
	podNs := "default"
	containerName := "axe-ui"
	var restConf *rest.Config
	restConf = utils.GetRestConf()
	Req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(podNs).
		SubResource("exec").
		VersionedParams(&core_v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	fmt.Println(Req.URL())
	executor, err := remotecommand.NewSPDYExecutor(restConf, "POST", Req.URL())
	if err != nil {
		fmt.Println("END")
	}
	fmt.Println("executor:")

	// ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	handler := &s{ws: ws, resizeEvent: make(chan remotecommand.TerminalSize)}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	})

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.wsClose()

	fmt.Println("executor1:")

}

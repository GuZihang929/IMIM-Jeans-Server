package conn

import (
	"IM-Server/global"
	"IM-Server/im"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"

	"net/http"
	"time"
)

type WsServerOptions struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type WsServer struct {
	options  *WsServerOptions
	upgrader websocket.Upgrader
	handler  func(ctx context.Context, conn Connection)
}

func NewWsServer(options *WsServerOptions) WsServer {

	if options == nil {
		options = &WsServerOptions{
			ReadTimeout:  8 * time.Minute,
			WriteTimeout: 8 * time.Minute,
		}
	}
	ws := new(WsServer)
	ws.options = options
	ws.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 65536,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return *ws
}

func (ws *WsServer) handleWebSocketRequest(writer http.ResponseWriter, request *http.Request) {
	//获取用户id
	userId := request.Header.Get("user_id")
	ctx := context.WithValue(context.Background(), "userId", userId)
	i, err2 := strconv.ParseInt(userId, 10, 64)
	if err2 != nil {
		global.Logger.Error("userId数据转换出错：" + err2.Error())
	}
	//在线人数++，用户存放redis
	global.Redis.Set(context.Background(), im.GetRedisKeyMain(i), 123, 0)

	conn, err := ws.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println("出错了")
		return
	}

	proxy := ConnectionProxy{
		conn: NewWsConnection(conn, ws.options),
	}
	ws.handler(ctx, proxy)
}

func (ws *WsServer) SetConnHandler(handler func(ctx context.Context, conn Connection)) {
	ws.handler = handler
}

func (ws *WsServer) Run(host string, port int) error {

	http.HandleFunc("/ws", ws.handleWebSocketRequest)

	addr := fmt.Sprintf("%s:%d", host, port)

	global.Logger.Info("路由添加成功")
	if err := http.ListenAndServe(addr, nil); err != nil {
		global.Logger.Error("监听出错：" + err.Error())
		return err
	}
	return nil
}

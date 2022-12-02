package websocket

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
	"github.com/mlctrez/goapp-audioplayer/model"
	"nhooyr.io/websocket"
	"strings"
	"time"
)

var _ app.Mounter = (*WebSocket)(nil)

func Action(ctx app.Context) *Actions {
	return &Actions{ctx}
}

type Actions struct {
	app.Context
}

func (ac *Actions) Write(message model.WebSocketMessage) {
	ac.NewActionWithValue("websocket.write", message)
}

func (ac *Actions) HandleAction(messageType model.WebSocketMessage, handler func(message model.WebSocketMessage)) {
	ac.Handle("websocket.read."+messageType.WebSocketMessageName(), func(c app.Context, action app.Action) {
		handler(action.Value.(model.WebSocketMessage))
	})
}

func (ac *Actions) handle(webSocket *WebSocket) {
	ac.Handle("websocket.write", webSocket.write)
}

type WebSocket struct {
	app.Compo
	clientId      string
	conn          *websocket.Conn
	wsContext     context.Context
	wsCancel      context.CancelFunc
	earlyMessages []model.WebSocketMessage
}

func (w *WebSocket) Render() app.UI {
	return nodisplay.NoDisplay("websocket")
}

func (w *WebSocket) establishClientId(ctx app.Context) {
	var err error
	err = ctx.LocalStorage().Get("websocket.client.id", &w.clientId)
	if err != nil {
		app.Log("error getting client id from local storage", err)
	}
	if w.clientId == "" {
		w.clientId = uuid.NewString()
		err = ctx.LocalStorage().Set("websocket.client.id", w.clientId)
		if err != nil {
			app.Log("error setting client id to local storage", err)
		}
	}
}

func (w *WebSocket) write(ctx app.Context, action app.Action) {
	if wsm, ok := action.Value.(model.WebSocketMessage); ok {

		data, err := wsm.WebSocketMessage()
		//app.Log("write message ", wsm.WebSocketMessageName())

		if err != nil {
			app.Log("websocket.write serialize error", err)
			return
		}
		if w.conn == nil {
			fmt.Println("w.conn is nil, queueing message")
			w.earlyMessages = append(w.earlyMessages, wsm)
			return
		}

		if err = w.conn.Write(w.wsContext, websocket.MessageBinary, data); err != nil {
			app.Log("websocket.write err", err)
			w.cancelReconnect(ctx)
		}

	}
}

func (w *WebSocket) OnMount(ctx app.Context) {
	fmt.Println("compo/websocket/Websocket.OnMount")
	w.establishClientId(ctx)
	Action(ctx).handle(w)
	ctx.Async(func() { w.connectWebSocket(ctx) })
}

func (w *WebSocket) cancelReconnect(ctx app.Context) {
	w.wsCancel()
	ctx.After(time.Second*10, func(c app.Context) {
		c.Async(func() { w.connectWebSocket(c) })
	})
}

func (w *WebSocket) connectWebSocket(ctx app.Context) {
	w.wsContext, w.wsCancel = context.WithCancel(ctx)
	wsUri := fmt.Sprintf("%s/ws/%s", BaseWs(), w.clientId)
	var err error
	if w.conn, _, err = websocket.Dial(w.wsContext, wsUri, nil); err != nil {
		//app.Log("connectWebSocket websocket.Dial error:", err)
		w.cancelReconnect(ctx)
	} else {
		// bump up the max payload size
		w.conn.SetReadLimit(1024 * 1024)

		go w.readMessage(ctx)
	}
}

func (w *WebSocket) readMessage(ctx app.Context) {
	fmt.Println("compo/websocket/Websocket.readMessage")
	for {
		var err error
		var data []byte
		var wsType websocket.MessageType
		if wsType, data, err = w.conn.Read(w.wsContext); err != nil {
			w.cancelReconnect(ctx)
			return
		}
		if wsType == websocket.MessageBinary {
			var msg model.WebSocketMessage
			if msg, err = model.DecodeResponse(data); err != nil {
				app.Logf("model.DecodeResponse error: %s", err)
				return
			}
			//app.Log("read message ", msg.WebSocketMessageName())
			ctx.NewActionWithValue("websocket.read."+msg.WebSocketMessageName(), msg)
		}
	}
}

// Base returns the base url, removing any trailing slash.
func Base() string {
	href := app.Window().Get("location").Get("href").String()
	return strings.TrimSuffix(href, "/")
}

// BaseWs is the same as Base but transforms http -> ws and https -> wss.
func BaseWs() string {
	return "ws" + strings.TrimPrefix(Base(), "http")
}

package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/model"
	"net/http"
	"nhooyr.io/websocket"
	"strings"
)

func (a *Api) webSocketHandler(ginCtx *gin.Context) {

	var err error
	var conn *websocket.Conn

	clientId := ginCtx.Param("clientId")
	//fmt.Println("websocket connect", clientId, ginCtx.Request.RemoteAddr)

	var options *websocket.AcceptOptions
	// https://github.com/gorilla/websocket/issues/731
	// Compression in certain Safari browsers is broken, turn it off
	if strings.Contains(ginCtx.Request.UserAgent(), "Safari") {
		options = &websocket.AcceptOptions{CompressionMode: websocket.CompressionDisabled}
	}

	if conn, err = websocket.Accept(ginCtx.Writer, ginCtx.Request, options); err != nil {
		_ = ginCtx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	defer cancelFunc()

	// make a response channel
	responseChannel := make(chan []byte, 10)
	defer close(responseChannel)

	go func() {
		for {

			var bytes []byte
			var readErr error

			if _, bytes, readErr = conn.Read(ctx); readErr != nil {
				if !isAcceptedReadError(readErr.Error()) {
					fmt.Println("webSocketHandler readErr", readErr)
				}
				cancelFunc()
				return
			}

			var result []byte
			if result, err = model.InvokeApi(clientId, bytes, a.c); err != nil {
				app.Logf("error invoking api %s", err)
			}
			if result != nil {
				responseChannel <- result
			}
		}
	}()

	var selectRunning = true

	for selectRunning {
		select {
		case data := <-responseChannel:
			if data != nil {
				if writeErr := conn.Write(ctx, websocket.MessageBinary, data); writeErr != nil {
					fmt.Println("webSocketHandler writeErr", writeErr)
					selectRunning = false
				}
			}
		case <-ctx.Done():
			selectRunning = false
		}
	}

	//fmt.Println("websocket disconnect", clientId, ginCtx.Request.RemoteAddr)

}

func isAcceptedReadError(msg string) bool {
	if strings.Contains(msg, "received close frame") {
		return true
	}
	if strings.Contains(msg, "failed to read frame header: EOF") {
		return true
	}

	return false
}

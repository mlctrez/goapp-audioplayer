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

	// TODO: refactor this into struct with connection

	if conn, acceptErr := websocket.Accept(ginCtx.Writer, ginCtx.Request, nil); acceptErr != nil {
		_ = ginCtx.AbortWithError(http.StatusInternalServerError, acceptErr)
	} else {

		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		clientId := ginCtx.Param("clientId")

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

				result, err := model.InvokeApi(clientId, bytes, a.c)
				if err != nil {
					// TODO: send back an error type?
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
	}
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

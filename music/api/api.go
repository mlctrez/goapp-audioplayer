package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mlctrez/goapp-audioplayer/music"
	"github.com/mlctrez/goapp-natsws/proxy"
)

type Api struct {
	c *music.Catalog
}

func New(c *music.Catalog) *Api {
	return &Api{c: c}
}

func (a *Api) Register(engine *gin.Engine) {

	// cover art by release group id
	engine.GET("/cover/:uuid", a.getCover)

	// allows updating cover art with a new image
	engine.POST("/cover/:uuid", a.setCover)

	// flac file keyed by releaseGroupId_disc_track
	engine.GET("/flac/:key", a.releaseGroupTrack)

	// web socket endpoint for websocket api calls, not used - kept for reference
	engine.GET("/ws/:clientId", a.webSocketHandler)

	engine.GET("/natsws/:clientId", gin.WrapH(proxy.New(music.NatsWebsocketURL())))

}

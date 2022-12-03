package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
)

func (a *Api) getCover(ctx *gin.Context) {

	uu, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var pngBytes []byte

	pngBytes, err = a.c.CoverArt(uu)
	if err != nil || len(pngBytes) < 1000 {
		var logoBytes []byte
		logoBytes, err = goapp.WebFs.ReadFile("web/logo-512.png")
		if err != nil {
			log.Println(err)
		}
		ctx.Header("Cache-Control", fmt.Sprintf("max-age=%d", 60*60*3))
		ctx.Data(200, "image/png", logoBytes)
		return
	}

	etag := fmt.Sprintf("%q", uu.String())

	ctx.Header("Etag", etag)
	ctx.Header("Cache-Control", fmt.Sprintf("max-age=%d", 60*60*24*365))

	if ctx.GetHeader("If-None-Match") == etag {
		ctx.Status(http.StatusNotModified)
	} else {
		ctx.Data(200, "image/png", pngBytes)
	}

}

func (a *Api) setCover(ctx *gin.Context) {

	defer func() { _ = ctx.Request.Body.Close() }()

	uu, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		fmt.Println("setCover parse", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var img image.Image
	if img, _, err = image.Decode(ctx.Request.Body); err != nil {
		fmt.Println("setCover decode", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = a.c.SetCoverArt(uu, img); err != nil {
		fmt.Println("setCover catalog.SetCoverArt", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusAccepted)
}

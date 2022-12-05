package api

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"net/http"
	"strconv"
)

func (a *Api) getCover(ctx *gin.Context) {

	uu, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var size int64
	sizeQuery := ctx.Query("size")
	if sizeQuery != "" {
		size, err = strconv.ParseInt(sizeQuery, 10, 64)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
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

	var etag string

	if size > 0 {
		// deny unreasonable service, else it's denial of service
		if size > 800 {
			size = 800
		}
		etag = fmt.Sprintf("%q", fmt.Sprintf("%s_%d", uu.String(), size))

		if pngBytes, err = resizeImage(pngBytes, int(size)); err != nil {
			fmt.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

	} else {
		etag = fmt.Sprintf("%q", uu.String())
	}

	ctx.Header("Etag", etag)
	ctx.Header("Cache-Control", fmt.Sprintf("max-age=%d", 60*60*24*365))

	if ctx.GetHeader("If-None-Match") == etag {
		ctx.Status(http.StatusNotModified)
	} else {
		ctx.Data(200, "image/png", pngBytes)
	}

}

func resizeImage(imageBytes []byte, size int) (result []byte, err error) {
	var img image.Image
	if img, err = png.Decode(bytes.NewReader(imageBytes)); err != nil {
		return
	}
	resize := imaging.Resize(img, int(size), 0, imaging.Lanczos)
	buff := &bytes.Buffer{}
	if err = png.Encode(buff, resize); err != nil {
		return
	}
	return buff.Bytes(), nil
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

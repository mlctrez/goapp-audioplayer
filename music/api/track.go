package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mlctrez/goapp-audioplayer/model"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (a *Api) releaseGroupTrack(context *gin.Context) {
	key := context.Param("key")
	// TODO: some sanity checks on the key format

	var err error
	var md *model.Metadata

	if md, err = a.c.ReleaseDiscTrackMetadata(key); err != nil {
		_ = context.AbortWithError(http.StatusNotFound, err)
	}

	var flacFile *os.File
	if flacFile, err = os.Open(md.Path); err != nil {
		_ = context.AbortWithError(http.StatusInternalServerError, err)
	}
	defer func() { _ = flacFile.Close() }()

	rangeHeader := context.Request.Header.Get("Range")
	rangeParts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")

	// TODO: parse end? it does not ever seem to be sent
	// TODO: support when range is not sent.. Chrome seems to send it always

	var start int64
	if start, err = strconv.ParseInt(rangeParts[0], 10, 64); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = flacFile.Seek(start, io.SeekStart)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	maxSent := int64(2 * 1024 * 1024)

	end := start + maxSent
	if end > md.Size-1 {
		end = md.Size - 1
	}

	context.Status(http.StatusPartialContent)
	context.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, md.Size))
	context.Header("Content-Type", "audio/flac")

	// errors are ignored here, connection reset, broken pipe and eof
	_, _ = io.CopyN(context.Writer, flacFile, maxSent)
}

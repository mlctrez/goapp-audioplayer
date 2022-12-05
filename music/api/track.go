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

	var err error
	var md *model.Metadata

	if md, err = a.c.ReleaseDiscTrackMetadata(key); err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	var flacFile *os.File
	if flacFile, err = os.Open(md.Path); err != nil {
		context.Status(http.StatusNotFound)
		return
	}
	defer func() { _ = flacFile.Close() }()

	rangeHeader := context.Request.Header.Get("Range")
	if rangeHeader == "" {
		// no range then you get the whole file
		context.Header("Content-Type", "audio/flac")

		// errors are ignored here, connection reset, broken pipe and eof
		_, _ = io.Copy(context.Writer, flacFile)
		return
	}

	rangeParts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")

	// TODO range checking on start and end ranges against file size
	var startBytes int64
	if startBytes, err = strconv.ParseInt(rangeParts[0], 10, 64); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	var endBytes int64
	if len(rangeParts[1]) > 0 {
		if endBytes, err = strconv.ParseInt(rangeParts[1], 10, 64); err != nil {
			context.Status(http.StatusBadRequest)
			return
		}
	}

	_, err = flacFile.Seek(startBytes, io.SeekStart)
	if err != nil {
		_ = context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// send 2Mb at a time if no end range provided
	maxSent := int64(2 * 1024 * 1024)

	if endBytes > 0 {
		maxSent = (endBytes - startBytes) + 1
	}

	end := startBytes + maxSent
	if end > md.Size-1 {
		end = md.Size - 1
	}

	debugRange := os.Getenv("GOAPP_DEBUG_RANGE") != ""

	if debugRange {
		fmt.Println("trackSize, startBytes, endBytes =", md.Size, startBytes, endBytes)
	}

	context.Status(http.StatusPartialContent)
	context.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startBytes, end, md.Size))
	context.Header("Content-Type", "audio/flac")

	var written int64
	written, err = io.CopyN(context.Writer, flacFile, maxSent)
	if debugRange {
		fmt.Println("bytesWritten, err =", written, err)
	}
}

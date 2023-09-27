//go:build !wasm

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	brotli "github.com/anargu/gin-brotli"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo"
	"github.com/mlctrez/goapp-audioplayer/model"
	"github.com/mlctrez/goapp-audioplayer/music"
	"github.com/mlctrez/goapp-audioplayer/music/api"
	"github.com/mlctrez/goapp-audioplayer/music/natsapi"
	"github.com/mlctrez/servicego"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"io/fs"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func Entry() {
	compo.Routes()
	servicego.Run(&Service{})
}

var _ servicego.Service = (*Service)(nil)

type Service struct {
	servicego.Defaults
	serverShutdown func(ctx context.Context) error
	db             *music.Catalog
	nats           *server.Server
	serverContext  context.Context
	serverCancel   context.CancelFunc
}

func (s *Service) Start(_ service.Service) (err error) {

	fmt.Println("starting version", goapp.RuntimeVersion())

	s.serverContext, s.serverCancel = context.WithCancel(context.Background())

	if s.db, err = music.OpenCatalog("bolt.db"); err != nil {
		return
	}

	var listener net.Listener
	address := listenAddress()
	if listener, err = net.Listen("tcp4", address); err != nil {
		return
	}
	dev := goapp.IsDevelopment()

	if dev {
		//goland:noinspection HttpUrlsUsage
		fmt.Printf("running on http://%s\n", address)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	// required for go-app to work correctly
	engine.RedirectTrailingSlash = false

	reduceNoise := os.Getenv("GOAPP_LOG_ALL_REQUESTS") == ""

	if dev {
		if reduceNoise {
			logConfig := gin.LoggerConfig{SkipPaths: []string{"/app-worker.js", "/web"}}
			engine.Use(gin.LoggerWithConfig(logConfig), gin.Recovery())
		} else {
			engine.Use(gin.Logger(), gin.Recovery())
		}
	} else {
		engine.Use(gin.Recovery(), brotli.Brotli(brotli.DefaultCompression))
	}

	if err = setupGinStaticHandlers(engine); err != nil {
		return
	}

	if err = s.startEmbeddedNats(); err != nil {
		return
	}

	// other api endpoints can go here
	api.New(s.db).Register(engine)

	var handler *app.Handler
	if handler, err = BuildHandler(); err != nil {
		return
	}

	h := gin.WrapH(handler)
	engine.NoRoute(func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		h(c)
	})

	server := &http.Server{Handler: engine}
	s.serverShutdown = server.Shutdown

	go func() {
		var serveErr error
		if strings.HasSuffix(listener.Addr().String(), ":443") {
			serveErr = server.ServeTLS(listener, "cert.pem", "cert.key")
		} else {
			serveErr = server.Serve(listener)
		}
		if serveErr != nil && serveErr != http.ErrServerClosed {
			_ = s.Log().Error(err)
		}
	}()

	return nil
}

func (s *Service) startEmbeddedNats() (err error) {

	var natsHost string
	var natsPort int
	if natsHost, natsPort, err = music.NatsAddress(); err != nil {
		return
	}

	options := &server.Options{
		ServerName: "audioPlayer",
		Host:       natsHost,
		Port:       natsPort,
		NoSigs:     true,
		Websocket: server.WebsocketOpts{
			Host:  natsHost,
			Port:  natsPort + music.NatsWebsocketPortOffset,
			NoTLS: true,
		},
	}

	if s.nats, err = server.NewServer(options); err != nil {
		return
	}

	// to see nats startup failures, uncomment the next two lines
	// logger := logger.NewTestLogger("nats", false)
	// s.nats.SetLogger(logger, false, false)
	go s.nats.Start()

	if !s.nats.ReadyForConnections(time.Second * 4) {
		return fmt.Errorf("embedded nats did not start correctly")
	}

	// remove logger after successful startup
	s.nats.SetLogger(nil, false, false)

	var natsConn *nats.Conn
	if natsConn, err = nats.Connect("", nats.InProcessServer(s.nats)); err != nil {
		return
	}

	backend := model.NewNatsBackend(s.serverContext, natsConn, &natsapi.Api{Catalog: s.db})

	return backend.Start()
}

func setupGinStaticHandlers(engine *gin.Engine) (err error) {

	var wasmFile fs.File
	if wasmFile, err = goapp.WebFs.Open("web/app.wasm"); err != nil {
		return
	}
	defer func() { _ = wasmFile.Close() }()

	var stat fs.FileInfo
	if stat, err = wasmFile.Stat(); err != nil {
		return
	}
	wasmSize := stat.Size()

	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Wasm-Content-Length", fmt.Sprintf("%d", wasmSize))
		c.Next()
	})

	staticHandler := http.FileServer(http.FS(goapp.WebFs))
	engine.GET("/web/:path", gin.WrapH(staticHandler))

	if _, err = fs.Stat(goapp.WebFs, "web/app.css"); err == nil {
		//  use provided web/app.css instead of app.css provided by go-app
		engine.GET("/app.css", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/web/app.css")
		})
	} else {
		err = nil
	}

	return
}

func (s *Service) Stop(_ service.Service) (err error) {

	if s.db != nil {
		s.db.CloseCatalog()
	}

	if s.serverCancel != nil {
		s.serverCancel()
	}

	if s.nats != nil {
		s.nats.Shutdown()
	}

	if s.serverShutdown != nil {

		stopContext, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = s.serverShutdown(stopContext)
	}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			debug.PrintStack()
			os.Exit(-1)
		}
	} else {
		_ = s.Log().Info("http.Server.Shutdown success")
	}
	return
}

func listenAddress() string {
	if address := os.Getenv("ADDRESS"); address != "" {
		return address
	}
	if port := os.Getenv("PORT"); port == "" {
		return "localhost:8080"
	} else {
		return "localhost:" + port
	}

}

func BuildHandler() (handler *app.Handler, err error) {

	var file fs.File
	if file, err = goapp.WebFs.Open("web/handler.json"); err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	handler = &app.Handler{}
	if err = json.NewDecoder(file).Decode(handler); err != nil {
		return
	}

	appWorkerJs := app.DefaultAppWorkerJS
	for k, v := range appWorkerJsReplace {
		appWorkerJs = strings.Replace(appWorkerJs, k, v, 1)
	}

	handler.ServiceWorkerTemplate = appWorkerJs

	handler.Version = goapp.RuntimeVersion()
	handler.AutoUpdateInterval = goapp.UpdateInterval()
	if goapp.IsDevelopment() {
		handler.Env["DEV"] = "1"
	}
	handler.WasmContentLengthHeader = "Wasm-Content-Length"

	return
}

var appWorkerJsReplace = map[string]string{
	"key !== cacheName": `key !== cacheName || key === "DYNAMIC"`,
	"return response || fetch(event.request);": `if (response) {
        return response;
      }
      if (event.request.url.indexOf("/cover/") === -1) {
        return fetch(event.request);
      }
      return fetch(event.request).then((response) => {
        return caches.open("DYNAMIC").then((cache) => {
          cache.put(event.request.url, response.clone()).then(() => {
          });
          return response;
        })
      })
`,
}

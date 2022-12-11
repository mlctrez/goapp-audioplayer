package generator

import (
	"encoding/json"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/mlctrez/goapp-audioplayer/model/generate/spec"
	"os"
	"path/filepath"
	"strings"
)

func Run() error {
	return (&generator{}).run()
}

type generator struct {
	spec                  *spec.Spec
	qualifier             string
	webSocketMessageTypes map[string]bool
}

func (g *generator) initStruct() (err error) {
	g.webSocketMessageTypes = map[string]bool{}
	g.spec = &spec.Spec{}
	return nil
}

func execute(functions ...func() error) (err error) {
	for _, function := range functions {
		if err = function(); err != nil {
			return
		}
	}
	return nil
}

func (g *generator) run() (err error) {

	err = execute(
		g.initStruct,
		g.readSpec,
		g.fillWebsocketTypes,
		g.generateCommon,
		g.generateRequestResponse,
		g.generateWebsocketApi,
		g.natsClientApi,
		g.natsBackendApi,
	)
	return
}

func (g *generator) readSpec() (err error) {

	var specBytes []byte
	if specBytes, err = os.ReadFile("model/spec.json"); err != nil {
		return
	}

	if err = json.Unmarshal(specBytes, g.spec); err != nil {
		return
	}
	g.qualifier = filepath.Join(g.spec.Module, "model")
	return
}

func (g *generator) newFile() *jen.File {
	jp := jen.NewFilePath(g.qualifier)
	jp.PackageComment("Package model is auto generated from model/spec.json - DO NOT EDIT")
	return jp
}

func (g *generator) fillWebsocketTypes() (err error) {
	for _, m := range g.spec.Methods {
		g.webSocketMessageTypes[m.Request] = true
		g.webSocketMessageTypes[m.Response] = true
	}
	return
}

func (g *generator) generateCommon() (err error) {
	jp := g.newFile()
	for _, t := range g.spec.Types {
		if !g.webSocketMessageTypes[t.Name] {
			jp.Line()
			t.GenerateGo(jp, g.qualifier, false)
		}
	}
	return g.render("common", jp)
}

func (g *generator) generateRequestResponse() (err error) {
	for _, m := range g.spec.Methods {
		jp := g.newFile()
		for _, t := range g.spec.Types {
			if t.Name == m.Response || t.Name == m.Request {
				jp.Line()
				t.GenerateGo(jp, g.qualifier, true)
			}
		}
		if err = g.render(m.Name, jp); err != nil {
			return
		}
	}
	return
}

func (g *generator) render(name string, jf *jen.File) error {

	goFileName := fmt.Sprintf("%s.go", strings.ToLower(name))

	path := filepath.Join(filepath.Base(g.qualifier), goFileName)

	create, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = create.Close() }()

	return jf.Render(create)
}

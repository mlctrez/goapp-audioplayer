package main

import (
	"bytes"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/mlctrez/goapp-audioplayer/internal/gomod"
	"os"
	"path/filepath"
	"strings"
)

func failErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const goApp = "github.com/maxence-charriere/go-app/v9/pkg/app"

func main() {
	compoName := "Test"
	receiver := strings.ToLower(compoName[0:1])
	moduleDir := filepath.Join("goapp", "compo", "test")

	failErr(os.MkdirAll(moduleDir, 0755))

	name, err := gomod.ModuleName()
	failErr(err)

	jf := NewFilePath(filepath.Join(name, moduleDir))

	jf.Type().Id(compoName).Struct(Qual(goApp, "Compo"))

	receiverFunc := func(name string) *Statement {
		jf.Line()
		return jf.Func().Params(Id(receiver).Op("*").Id(compoName)).Id(name)
	}

	ctxParam := func() *Statement { return Id("ctx").Qual(goApp, "Context") }

	receiverFunc("Render").Params().Params(Qual(goApp, "UI")).Block(Return(Qual(goApp, "Div").Params()))

	funcs := []string{"OnMount", "OnNav", "OnUpdate", "OnAppUpdate", "OnAppInstallChange", "OnResize"}
	for _, s := range funcs {
		receiverFunc(s).Params(ctxParam()).Params().Block(Line())
	}

	receiverFunc("OnInit").Params().Params().Block(Line())
	receiverFunc("OnDismount").Params().Params().Block(Line())

	interfaces := []string{"Initializer", "Dismounter", "Mounter", "Navigator", "Updater", "AppUpdater", "Resizer"}
	for _, s := range interfaces {
		jf.Var().Id("_").Qual(goApp, s).Op("=(*").Id(compoName).Op(")(nil)")
	}

	buf := &bytes.Buffer{}
	failErr(jf.Render(buf))
	failErr(os.WriteFile(filepath.Join(moduleDir, "test.go"), buf.Bytes(), 0644))
}

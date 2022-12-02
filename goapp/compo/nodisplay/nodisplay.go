package nodisplay

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func NoDisplay(name string) app.HTMLDiv {
	return app.Div().Style("display", "none").DataSet("name", name)
}

package test

import app "github.com/maxence-charriere/go-app/v9/pkg/app"

type Test struct {
	app.Compo
}

func (t *Test) Render() app.UI {
	return app.Div()
}

func (t *Test) OnMount(ctx app.Context) {

}

func (t *Test) OnNav(ctx app.Context) {

}

func (t *Test) OnUpdate(ctx app.Context) {

}

func (t *Test) OnAppUpdate(ctx app.Context) {

}

func (t *Test) OnAppInstallChange(ctx app.Context) {

}

func (t *Test) OnResize(ctx app.Context) {

}

func (t *Test) OnInit() {

}

func (t *Test) OnDismount() {

}

var _ app.Initializer = (*Test)(nil)
var _ app.Dismounter = (*Test)(nil)
var _ app.Mounter = (*Test)(nil)
var _ app.Navigator = (*Test)(nil)
var _ app.Updater = (*Test)(nil)
var _ app.AppUpdater = (*Test)(nil)
var _ app.Resizer = (*Test)(nil)

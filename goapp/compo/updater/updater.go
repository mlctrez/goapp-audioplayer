package updater

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
)

var _ app.AppUpdater = (*Updater)(nil)

type Updater struct {
	app.Compo
}

func (u *Updater) OnAppUpdate(ctx app.Context) {
	if goapp.IsDevelopment() && ctx.AppUpdateAvailable() {
		ctx.Reload()
	}
}

func (u *Updater) Render() app.UI {
	return nodisplay.NoDisplay("updater").DataSet("runtimeVersion", goapp.RuntimeVersion())
}

package nav

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/actions"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
	"github.com/mlctrez/goapp-audioplayer/model"
	"github.com/mlctrez/goapp-natsws/natsws"
	"strings"
)

var _ app.Mounter = (*Navigation)(nil)

type Navigation struct {
	app.Compo
	natswsConn *natsws.Connection
	width      int
}

func (n *Navigation) OnMount(ctx app.Context) {
	n.natswsConn = &natsws.Connection{}
	natsws.Observe(ctx, n.natswsConn)
}

func (n *Navigation) Render() app.UI {
	return app.Div().Class("navigation").Body(
		app.Div().Class("navigation-left").Body(
			app.Div().Title("mlctrez Music").Body(app.Img().Src("/web/logo-192.png")),
			app.Div().Body(app.Text("mlctrez Music")),
		).OnClick(n.navigationLeft),
		app.Div().Class("navigation-center").Body(
			&Icon{StateName: "navigation.previous", SvgFunc: icon.NavigateBefore48},
			app.Div().Body(app.Text("Randomize")).OnClick(n.requestAlbums),
			&Icon{StateName: "navigation.next", SvgFunc: icon.NavigateNext48},
		).Style("cursor", "pointer"),
		n.version(),
	)
}

func (n *Navigation) requestAlbums(ctx app.Context, _ app.Event) {
	actions.RequestAlbums(ctx, n.natswsConn, &model.AlbumsRequest{})
}

func (n *Navigation) version() app.UI {
	div := app.Div().Class("navigation-right")
	if goapp.IsDevelopment() {
		div.Text(goapp.RuntimeVersion()[0:5])
	} else {
		version := strings.Split(goapp.RuntimeVersion(), "@")[0]
		href := fmt.Sprintf("https://github.com/mlctrez/goapp-audioplayer/tree/%s", version)
		div.Body(app.A().Href(href).Text(version))
	}
	return div
}

func (n *Navigation) navigationLeft(ctx app.Context, _ app.Event) {
	ctx.Reload()
}

type Icon struct {
	app.Compo
	SvgFunc    func() string
	StateName  string
	stateValue string
}

func (i *Icon) Render() app.UI {
	comp := app.Div().Body(app.Raw(i.SvgFunc()))
	if i.stateValue != "" {
		comp.Class("navigation-icon")
		comp.OnClick(i.click)
	} else {
		comp.Class("navigation-icon-disabled")
	}
	return comp
}

func (i *Icon) OnMount(ctx app.Context) {
	ctx.ObserveState(i.StateName).Value(&i.stateValue)
}

func (i *Icon) click(ctx app.Context, _ app.Event) {
	ctx.NewAction(i.StateName)
}

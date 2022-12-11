package generator

import (
	"github.com/dave/jennifer/jen"
	"strings"
)

func LcFirst(in string) string {
	return strings.ToLower(in[0:1]) + in[1:]
}

func (g *generator) natsBackendApi() (err error) {

	var natsApiInterfaceMethods []jen.Code

	// generate individual structs, capturing case statements for api
	for _, m := range g.spec.Methods {

		natsApiInterfaceMethods = append(natsApiInterfaceMethods,
			jen.Line().Id(m.Name).Params(
				jen.Id("request").Op("*").Qual(g.qualifier, m.Request),
			).Params(
				jen.Id("response").Op("*").Qual(g.qualifier, m.Response),
			),
		)
	}

	na := g.newFile()

	apiName := "NatsBackendApi"
	apiBackend := "NatsBackend"
	backendImpl := "natsBackend"

	na.Line().Type().Id(apiName).Interface(natsApiInterfaceMethods...)

	na.Line().Type().Id(apiBackend).Interface(
		jen.Line().Id("Start").Params().Params(jen.Error()),
		jen.Line().Id("Stop").Params(),
	)

	na.Line().Func().Id("New"+apiBackend).Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("conn").Op("*").Qual(NatsPackage, "Conn"),
		jen.Id("api").Qual(g.qualifier, apiName),
	).Params(jen.Id("backend").Qual(g.qualifier, apiBackend)).Block(
		jen.Id("b").Op(":=&").Qual(g.qualifier, backendImpl).Op("{}"),
		jen.Line(),
		jen.Id("b").Dot("ctx").Op("=").Id("ctx"),
		jen.Id("b").Dot("conn").Op("=").Id("conn"),
		jen.Id("b").Dot("api").Op("=").Id("api"),
		jen.Line(),
		jen.Return(jen.Id("b")),
	)

	na.Line().Var().Id("_").Qual(g.qualifier, apiBackend).
		Op("=(*").Id(backendImpl).Op(")(nil)")

	na.Line().Type().Id(backendImpl).Struct(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("conn").Op("*").Qual(NatsPackage, "Conn"),
		jen.Id("api").Qual(g.qualifier, apiName),
		jen.Id("subCtx").Qual("context", "Context"),
		jen.Id("subCancel").Qual("context", "CancelFunc"),
	)

	receiver := func() *jen.Statement {
		return jen.Id("nb").Op("*").Qual(g.qualifier, backendImpl)
	}
	/*
		subCtx, subCancel := context.WithCancel(nb.ctx)

		testChan := make(chan *natsgo.Msg, 10)
		subscribe, err := nb.conn.ChanSubscribe("test", testChan)

	*/
	var initCode []jen.Code
	initCode = append(initCode,
		jen.Id("nb.subCtx").Op(",").Id("nb.subCancel").
			Op("=").Qual("context", "WithCancel").
			Params(jen.Id("nb.ctx")),
		jen.Var().Id("subs").Op("[]*").Qual(NatsPackage, "Subscription"),

		jen.Id("errUnsubAll").Op(":=").Func().Params().Block(
			jen.For(jen.Id("_").Op(",").Id("sub").
				Op(":=").Range().Id("subs")).Block(
				jen.Id("_").Op("=").Id("sub").Dot("Unsubscribe").Params(),
			),
		),
	)

	for _, m := range g.spec.Methods {
		subjectName := LcFirst(m.Name)
		chanName := subjectName + "Chan"
		subName := subjectName + "Sub"

		initCode = append(initCode,
			jen.Line().Id(chanName).Op(":=").Make(
				jen.Chan().Op("*").Qual(NatsPackage, "Msg"), jen.Lit(10),
			),

			jen.Var().Id(subName).Op("*").Qual(NatsPackage, "Subscription"),
			jen.If(
				jen.Id(subName).Op(",").Err().Op("=").
					Id("nb").Dot("conn").Dot("ChanSubscribe").
					Params(jen.Lit(subjectName), jen.Id(chanName)),
			).Op(";").Err().Op("!=").Nil().Block(
				jen.Id("errUnsubAll").Params(),
				jen.Return(),
			),
			jen.Id("subs").Op("=").Append(jen.Id("subs"), jen.Id(subName)),
		)

	}

	var startBlock []jen.Code
	startBlock = append(startBlock, initCode...)

	var selectCases []jen.Code

	selectCases = append(selectCases,
		jen.Case(jen.Op("<-").Id("nb").
			Dot("subCtx").Dot("Done").Params()).
			Block(
				jen.Return(),
			),
	)

	for _, m := range g.spec.Methods {
		chanName := LcFirst(m.Name) + "Chan"
		selectCases = append(selectCases,
			jen.Case(jen.Id("msg").Op(":= <-").Id(chanName)).
				Block(
					jen.If(jen.Id("msg").Op("==").Nil()).Block(jen.Continue()),

					jen.Line().Var().Id("response").Qual(g.qualifier, "WebSocketMessage"),
					jen.If(
						jen.Id("response").Op(",").Err().Op("=").Id("DecodeMessage").
							Params(jen.Id("msg").Dot("Data")),
					).Op(";").Err().Op("!=").Nil().Block(
						jen.Qual("fmt", "Println").Params(jen.Lit("bad message received"), jen.Err()),
					),

					jen.Line().If(
						jen.Id("wsm").Op(",").Id("ok").
							Op(":=").Id("response").Op(".(*").Id(m.Request).Op(")"),
					).Op(";").Id("ok").Block(
						jen.Id("resp").Op(":=").
							Id("nb").Dot("api").Dot(m.Name).Params(jen.Id("wsm")),

						jen.Line().Var().Id("message").Op("[]").Byte(),
						jen.If(
							jen.Id("message").Op(",").Err().Op("=").Id("resp").Dot("WebSocketMessage").Params(),
						).Op(";").Err().Op("!=").Nil().Block(
							jen.Qual("fmt", "Println").Params(jen.Lit("bad message from api"), jen.Err()),
						),

						jen.Line().If(
							jen.Err().Op("=").Id("nb").Dot("conn").Dot("Publish").Params(
								jen.Id("msg").Dot("Reply"),
								jen.Id("message"),
							),
						).Op(";").Err().Op("!=").Nil().Block(
							jen.Qual("fmt", "Println").Params(jen.Lit("cannot publish reply"), jen.Err()),
						),
					),
				),
		)
	}

	/*
		if message, err = resp.WebSocketMessage(); err != nil {
			fmt.Println("bad message api", err)
		}
		if err = nb.conn.Publish(msg.Reply, message); err != nil {
			fmt.Println("cannot publish reply", err)
		}

	*/

	startBlock = append(startBlock,
		jen.Line().Go().Func().Params().Block(
			jen.Defer().Id("errUnsubAll").Params(),
			jen.For().Block(
				jen.Select().Block(
					selectCases...,
				),
			),
		).Params(),
	)

	startBlock = append(startBlock, jen.Line().Return())

	na.Line().Func().Params(receiver()).Id("Start").Params().
		Params(jen.Id("err").Error()).Block(startBlock...)

	na.Line().Func().Params(receiver()).Id("Stop").Params().Block(
		jen.If(jen.Id("nb.subCancel").Op("!=").Nil()).Block(
			jen.Id("nb.subCancel").Params(),
		),
	)

	return g.render("natsbackend", na)

}

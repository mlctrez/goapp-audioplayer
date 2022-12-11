package generator

import (
	"github.com/dave/jennifer/jen"
)

const NatsPackage = "github.com/nats-io/nats.go"

func (g *generator) natsClientApi() (err error) {

	var natsApiInterfaceMethods []jen.Code

	// generate individual structs, capturing case statements for api
	for _, m := range g.spec.Methods {

		natsApiInterfaceMethods = append(natsApiInterfaceMethods,
			jen.Line().Id(m.Name).Params(
				jen.Id("request").Op("*").Qual(g.qualifier, m.Request),
				jen.Id("timeout").Qual("time", "Duration"),
			).Params(
				jen.Id("response").Op("*").Qual(g.qualifier, m.Response),
				jen.Err().Error(),
			),
		)

	}

	na := g.newFile()

	apiName := "NatsClientApi"
	apiImplementation := "natsClientApi"

	na.Type().Id(apiName).Interface(natsApiInterfaceMethods...)

	na.Line()
	na.Func().Id("New" + apiName).Params(
		jen.Id("conn").Op("*").Qual(NatsPackage, "Conn"),
	).Params(
		jen.Id("api").Qual(g.qualifier, apiName),
	).Block(
		jen.Return(jen.Op("&").Qual(g.qualifier, apiImplementation).
			Op("{").Id("conn").Op(":").Id("conn").Op("}")),
	)

	na.Line()
	na.Type().Id(apiImplementation).Struct(jen.Id("conn").Op("*").Qual(NatsPackage, "Conn"))

	/*
	   func (na *natsClientApi) invokeNats(subject string, message WebSocketMessage) (result WebSocketMessage, err error) {
	   	return
	   }
	*/

	receiver := func() *jen.Statement {
		return jen.Id("na").Op("*").Qual(g.qualifier, apiImplementation)
	}

	na.Func().Params(receiver()).Id("invokeNats").Params(
		jen.Id("subject").String(),
		jen.Id("message").Qual(g.qualifier, "WebSocketMessage"),
		jen.Id("timeout").Qual("time", "Duration"),
	).Params(
		jen.Id("result").Qual(g.qualifier, "WebSocketMessage"),
		jen.Err().Error(),
	).Block(
		jen.Var().Id("bytes").Op("[]").Byte(),
		jen.If(
			jen.Id("bytes").Op(",").Err().Op("=").
				Id("message").Dot("WebSocketMessage").Params(),
		).Op(";").Err().Op("!=").Nil().Block(jen.Return()),

		jen.Line(),
		jen.Var().Id("reply").Op("*").Qual(NatsPackage, "Msg"),
		jen.If(
			jen.Id("reply").Op(",").Err().Op("=").
				Id("na").Dot("conn").Dot("Request").Params(
				jen.Id("subject"),
				jen.Id("bytes"),
				jen.Id("timeout"),
			),
		).Op(";").Err().Op("!=").Nil().Block(jen.Return()),

		jen.Line(),

		jen.Return(jen.Qual(g.qualifier, "DecodeMessage").Params(jen.Id("reply").Dot("Data"))),
	)

	for _, m := range g.spec.Methods {

		na.Line()
		na.Func().Params(receiver()).Id(m.Name).Params(
			jen.Id("request").Op("*").Qual(g.qualifier, m.Request),
			jen.Id("timeout").Qual("time", "Duration"),
		).Params(
			jen.Id("response").Op("*").Qual(g.qualifier, m.Response),
			jen.Err().Error(),
		).Block(

			jen.Var().Id("wsm").Qual(g.qualifier, "WebSocketMessage"),
			jen.If(jen.Id("wsm").Op(",").Err().Op("=").Id("na").Dot("invokeNats").Params(
				jen.Lit(LcFirst(m.Name)),
				jen.Id("request"),
				jen.Id("timeout"),
			)).Op(";").Err().Op("!=").Nil().Block(jen.Return()),

			jen.Line(),
			jen.Var().Id("ok").Bool(),

			jen.If(
				jen.Id("response").Op(",").Id("ok").Op("=").
					Id("wsm").Op(".(*").Qual(g.qualifier, m.Response).Op(")"),
			).Op(";").Id("ok").Block(jen.Return()),

			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Params(
				jen.Lit("incorrect message received"),
			)),
		)

	}

	return g.render("natsclient", na)

}

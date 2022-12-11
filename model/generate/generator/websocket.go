package generator

import "github.com/dave/jennifer/jen"

func (g *generator) generateWebsocketApi() (err error) {
	at := g.newFile()

	var invokeCaseStatements []jen.Code
	var interfaceMethods []jen.Code
	var decodeCaseStatements []jen.Code

	// generate individual structs, capturing case statements for api
	for _, m := range g.spec.Methods {

		interfaceMethods = append(interfaceMethods,
			jen.Id(m.Name).Params(
				jen.Id("clientId").String(),
				jen.Id("request").Op("*").Qual(g.qualifier, m.Request),
			).Params(
				jen.Id("response").Op("*").Qual(g.qualifier, m.Response),
				jen.Err().Error(),
			),
		)

		invokeCaseStatements = append(invokeCaseStatements,
			jen.Case(jen.Lit(m.Request)).Block(
				jen.Id("request").Op(":=&").Id(m.Request).Op("{}"),
				jen.If(jen.Err().Op("=").Qual("encoding/json", "Unmarshal").
					Params(
						jen.Id("messageJson"),
						jen.Id("request"),
					).Op(";").Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
				jen.Var().Id("response").Op("*").Qual(g.qualifier, m.Response),
				jen.If(jen.Id("response").Op(",").Err().Op("=").Id("api").Dot(m.Name)).
					Params(
						jen.Id("clientId"),
						jen.Id("request"),
					).Op(";").Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.Return(jen.Id("response").Dot("WebSocketMessage").Params()),
		)

		decodeCaseStatements = append(decodeCaseStatements,
			jen.Case(jen.Lit(m.Request)).Block(
				jen.Id("result").Op("=&").Qual(g.qualifier, m.Request).Op("{}"),
			),
			jen.Case(jen.Lit(m.Response)).Block(
				jen.Id("result").Op("=&").Qual(g.qualifier, m.Response).Op("{}"),
			),
		)
	}

	decodeCaseStatements = append(decodeCaseStatements,
		jen.Default().Block(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Params(
				jen.Lit("message type %q not mapped"),
				jen.Id("messageType"),
			)),
		),
	)

	at.Type().Id("Api").Interface(interfaceMethods...)
	at.Line()

	at.Type().Id("WebSocketMessage").Interface(
		jen.Id("WebSocketMessage").Params().Params(jen.Op("[]").Byte(), jen.Error()),
		jen.Id("WebSocketMessageName").Params().Params(jen.String()),
	)
	at.Line()

	at.Func().Id("InvokeApi").
		Params(
			jen.Id("clientId").String(),
			jen.Id("data").Op("[]").Byte(),
			jen.Id("api").Qual(g.qualifier, "Api")).
		Params(
			jen.Id("result").Op("[]").Byte(),
			jen.Err().Error(),
		).Block(

		jen.Var().Id("messageType").String(),
		jen.Var().Id("messageJson").Op("[]").Byte(),
		jen.For(
			jen.Id("i").Op(":=").Lit(0),
			jen.Id("i").Op("<").Len(jen.Id("data")),
			jen.Id("i").Op("++"),
		).Block(
			jen.If(jen.Id("data").Op("[").Id("i").Op("]==").Lit(0)).Block(
				jen.Id("messageType").Op("=").String().Params(
					jen.Id("data").Op("[").Lit(0).Op(":").Id("i").Op("]"),
				),
				jen.Id("messageJson").Op("=").Id("data").Op("[").Id("i").Op("+1:]"),
				jen.Break(),
			),
		),
		jen.Line(),
		jen.Switch(jen.Id("messageType")).Block(invokeCaseStatements...),

		jen.Line(),

		jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Params(
			jen.Lit("message type %q not mapped"),
			jen.Id("messageType"),
		)),
	)

	at.Line()
	at.Func().Id("DecodeMessage").Params(
		jen.Id("data").Op("[]").Byte(),
	).Params(
		jen.Id("result").Qual(g.qualifier, "WebSocketMessage"),
		jen.Err().Error(),
	).Block(
		jen.Var().Id("messageType").String(),
		jen.Var().Id("messageJson").Op("[]").Byte(),
		at.Line(),
		jen.For(
			jen.Id("i").Op(":=").Lit(0),
			jen.Id("i").Op("<").Len(jen.Id("data")),
			jen.Id("i").Op("++"),
		).Block(
			jen.If(jen.Id("data").Op("[").Id("i").Op("]==").Lit(0)).Block(
				jen.Id("messageType").Op("=").String().Params(
					jen.Id("data").Op("[").Lit(0).Op(":").Id("i").Op("]"),
				),
				jen.Id("messageJson").Op("=").Id("data").Op("[").Id("i").Op("+1:]"),
				jen.Break(),
			),
		),
		at.Line(),
		jen.Switch(jen.Id("messageType")).Block(decodeCaseStatements...),

		jen.If(jen.Err().Op("=").Qual("encoding/json", "Unmarshal").
			Params(
				jen.Id("messageJson"),
				jen.Id("result"),
			).Op(";").Err().Op("!=").Nil()).
			Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
		jen.Return(),
	)

	return g.render("api", at)
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/dave/jennifer/jen"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	file, err := os.ReadFile("model/spec.json")
	if err != nil {
		panic(err)
	}
	spec := &Spec{}
	err = json.Unmarshal(file, spec)
	if err != nil {
		panic(err)
	}

	qualifier := filepath.Join(spec.Module, "model")

	var invokeCaseStatements []jen.Code
	var interfaceMethods []jen.Code
	var decodeCaseStatements []jen.Code

	webSocketMessageTypes := map[string]bool{}

	// generate individual structs, capturing case statements for api
	for _, m := range spec.Methods {

		webSocketMessageTypes[m.Request] = true
		webSocketMessageTypes[m.Response] = true

		interfaceMethods = append(interfaceMethods,
			jen.Id(m.Name).Params(
				jen.Id("clientId").String(),
				jen.Id("request").Op("*").Qual(qualifier, m.Request),
			).Params(
				jen.Id("response").Op("*").Qual(qualifier, m.Response),
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
				jen.Var().Id("response").Op("*").Qual(qualifier, m.Response),
				jen.If(jen.Id("response").Op(",").Err().Op("=").Id("api").Dot(m.Name)).
					Params(
						jen.Id("clientId"),
						jen.Id("request"),
					).Op(";").Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.Return(jen.Id("response").Dot("WebSocketMessage").Params()),
		)

		decodeCaseStatements = append(decodeCaseStatements,
			jen.Case(jen.Lit(m.Response)).Block(
				jen.Id("response").Op("=&").Qual(qualifier, m.Response).Op("{}"),
				jen.If(jen.Err().Op("=").Qual("encoding/json", "Unmarshal").
					Params(
						jen.Id("messageJson"),
						jen.Id("response"),
					).Op(";").Err().Op("!=").Nil()).
					Block(
						jen.Return(jen.Nil(), jen.Err()),
					),
				jen.Return(),
			),
		)

	}

	for _, m := range spec.Methods {
		jp := jen.NewFilePath(qualifier)
		jp.PackageComment("Package model is auto generated from model/spec.json - DO NOT EDIT")
		for _, t := range spec.Types {
			if t.Name == m.Response || t.Name == m.Request {
				jp.Line()
				t.generateGo(jp, qualifier, true)
			}
		}
		render(qualifier, m.Name, jp)
	}

	jp := jen.NewFilePath(qualifier)
	jp.PackageComment("Package model is auto generated from model/spec.json - DO NOT EDIT")

	for _, t := range spec.Types {
		if !webSocketMessageTypes[t.Name] {
			jp.Line()
			t.generateGo(jp, qualifier, webSocketMessageTypes[t.Name])
		}
	}
	render(qualifier, "common", jp)

	at := jen.NewFilePath(qualifier)
	at.PackageComment("Package model is auto generated from model/spec.json - DO NOT EDIT")
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
			jen.Id("api").Qual(qualifier, "Api")).
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
	at.Func().Id("DecodeResponse").Params(
		jen.Id("data").Op("[]").Byte(),
	).Params(
		jen.Id("response").Qual(qualifier, "WebSocketMessage"),
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

		jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Params(
			jen.Lit("unknown message type %q"),
			jen.Id("messageType"),
		)),
	)

	render(qualifier, "api", at)

}

func (t *Type) generateGo(jp *jen.File, qualifier string, websocketMessage bool) {

	var fields []jen.Code

	for _, field := range t.Fields {

		f := jen.Id(field.Name)

		switch field.Type {
		case "string":
			f.String()
		case "int64":
			f.Int64()
		case "uint64":
			f.Uint64()
		case "time.Time":
			f.Qual("time", "Time")

		default:
			fieldType := field.Type

			for _, op := range []string{"[]", "*"} {
				if strings.HasPrefix(fieldType, op) {
					fieldType = strings.TrimPrefix(fieldType, op)
					f.Op(op)
				}
			}
			f.Qual(qualifier, fieldType)
		}
		jsonTag := strings.ToLower(field.Name)
		if field.Json != "" {
			jsonTag = field.Json
		}
		f.Tag(map[string]string{"json": jsonTag + ",omitempty"})
		fields = append(fields, f)
	}
	jp.Type().Id(t.Name).Struct(fields...)

	receiver := jen.Id("m").Op("*").Id(t.Name)
	returns := []jen.Code{jen.Op("[]").Byte(), jen.Error()}

	if websocketMessage {
		jp.Var().Id("_").Qual(qualifier, "WebSocketMessage").Op("=").Op("(*").Id(t.Name).Op(")(").Nil().Op(")")

		jp.Func().Params(receiver).Id("WebSocketMessage").Params().Params(returns...).Block(
			jen.Id("result").Op(":=").Qual("bytes", "NewBufferString").Params(jen.Lit(t.Name)),
			jen.Id("result").Dot("WriteByte").Params(jen.Lit(0)),
			jen.Line(),
			jen.Err().Op(":=").Qual("encoding/json", "NewEncoder").
				Params(jen.Id("result")).Dot("Encode").Params(jen.Id("m")),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.Line(),
			jen.Return(jen.Id("result").Dot("Bytes").Params(), jen.Nil()),
		)

		jp.Line()

		jp.Func().Params(receiver).Id("WebSocketMessageName").Params().Params(jen.String()).Block(
			jen.Return(jen.Lit(t.Name)),
		)
	}

}

func render(qualifier, name string, jf *jen.File) {

	goFileName := fmt.Sprintf("%s.go", strings.ToLower(name))

	path := filepath.Join(filepath.Base(qualifier), goFileName)

	create, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func() { _ = create.Close() }()

	err = jf.Render(create)
	if err != nil {
		panic(err)
	}

}

type Spec struct {
	Module  string   `json:"module"`
	Methods []Method `json:"methods"`
	Types   []Type   `json:"types"`
}

type Method struct {
	Name     string `json:"name"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

type Type struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Json string `json:"json,omitempty"`
}

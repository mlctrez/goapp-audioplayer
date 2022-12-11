package spec

import (
	"github.com/dave/jennifer/jen"
	"strings"
)

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

func (t *Type) GenerateGo(jp *jen.File, qualifier string, websocketMessage bool) {

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

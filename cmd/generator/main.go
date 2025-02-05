package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"strings"
)

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

func typeJSFicator(t string, opts ...string) string {
	if strings.HasPrefix(t, "[]") {
		return typeJSFicator(t[2:], "[]")
	}
	if strings.Contains(t, "int") || strings.Contains(t, "float") {
		if len(opts) == 0 {
			return "number"
		}
		return fmt.Sprintf("number%s", opts[0])
	}
	if strings.Contains(t, "bool") {
		if len(opts) == 0 {
			return "boolean"
		}
		return fmt.Sprintf("boolean%s", opts[0])
	}
	if strings.Contains(t, "&") {
		r1 := strings.NewReplacer("&", "", "{", "", "}", "")
		root := r1.Replace(t)
		_type := strings.Split(root, " ")
		if len(opts) == 0 {
			return _type[1]
		}
		return fmt.Sprintf("%s%s", _type[1], opts[0])
	}
	if len(opts) == 1 {
		return fmt.Sprintf("%s%s", t, opts[0])
	}
	return t
}

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "internal/application/dto/user.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	var structs []Struct
	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if !strings.HasSuffix(strings.ToLower(typeSpec.Name.Name), "dto") {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		var fields []Field
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				fields = append(fields, Field{
					Name: strings.ToLower(name.Name),
					Type: typeJSFicator(fmt.Sprintf("%s", field.Type)),
				})
			}
		}

		structs = append(structs, Struct{
			Name:   typeSpec.Name.Name[:len(typeSpec.Name.Name)-3],
			Fields: fields,
		})
		return true
	})

	tmpl := `
		{{- range . }}
		type {{ .Name }} = {
		{{- range .Fields }}
			{{ .Name }}: {{ .Type }};
		{{- end }}
		};
		{{ end }}
	`

	t := template.Must(template.New("code").Parse(tmpl))
	if err := t.Execute(os.Stdout, structs); err != nil {
		fmt.Println(err)
		return
	}
}

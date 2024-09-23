package tygo

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func (g *PackageGenerator) Generate() (string, error) {
	s := new(strings.Builder)

	g.writeFileCodegenHeader(s)
	g.writeFileFrontmatter(s)

	filepaths := g.GoFiles

	availableImports := make([]string, 0)
	typespecmap := make(map[string]*ast.TypeSpec)

	for i, file := range g.pkg.Syntax {

		if g.conf.IsFileIgnored(filepaths[i]) {
			continue
		}

		first := true

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {

			// GenDecl can be an import, type, var, or const expression
			case *ast.GenDecl:
				if x.Tok == token.IMPORT {
					return false
				}
				isEmit := false
				if x.Tok == token.VAR {
					isEmit = g.isEmitVar(x)
					if !isEmit {
						return false
					}
				}

				if first {
					g.writeFileSourceHeader(s, filepaths[i], file)
					first = false
				}
				if isEmit {
					g.emitVar(s, x)
					return false
				}
				g.writeGroupDecl(s, x)

				for _, spec := range x.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if ok && ts.Name.IsExported() {

						typespecmap[ts.Name.Name] = ts
						availableImports = append(availableImports, ts.Name.Name)

					}
				}
				return false
			}
			return true

		})

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				a := &FunctionDoc{}
				a.ParseFromFuncDecl(x)
				a.Generate(s)
				availableImports = append(availableImports, a.LowerName)
			}
			return true
		})

	}

	return s.String(), nil
}

type FunctionDoc struct {
	Name                    string
	LowerName               string
	Method                  string
	ContentType             string
	Body                    string
	Response                string
	Path                    string
	SkipRequestGeneration   bool
	ShouldGenerateForm      bool
	FetchParams             string
	FetchParamsWithOutTypes string
	JsPath                  string
}

func (g *FunctionDoc) ParseFromFuncDecl(f *ast.FuncDecl) {
	if f == nil {
		g.SkipRequestGeneration = true
		return
	}

	if f.Doc == nil {
		g.SkipRequestGeneration = true
		return
	}

	if f.Doc.Text() == "" {
		g.SkipRequestGeneration = true
		return
	}

	for _, c := range f.Doc.List {

		text := strings.Replace(c.Text, "//", "", 1)

		if strings.Contains(text, "@method") {
			g.Method = strings.TrimSpace(strings.Replace(text, "@method", "", -1))
		}

		if strings.Contains(text, "@content-type") {
			g.ContentType = strings.TrimSpace(strings.Replace(text, "@content-type", "", -1))
		}

		if strings.Contains(text, "@body") {
			g.Body = strings.TrimSpace(strings.Replace(text, "@body", "", -1))
		}

		if strings.Contains(text, "@response") {
			g.Response = strings.TrimSpace(strings.Replace(text, "@response", "", -1))
		}

		if strings.Contains(text, "@path") {
			g.Path = strings.TrimSpace(strings.Replace(text, "@path", "", -1))
		}

		if strings.Contains(text, "@name") {
			g.Name = strings.TrimSpace(strings.Replace(text, "@name", "", -1))
		} else {
			name := f.Name.Name
			firstWord := string(name[0])
			rest := name[1:]
			g.LowerName = strings.ToLower(firstWord) + rest
			g.Name = name
		}

		if strings.Contains(text, "@skip") {
			g.SkipRequestGeneration = true
		}

		if strings.Contains(text, "@generateform") {
			g.ShouldGenerateForm = true
		}

	}
	params := ""
	paramsWithoutTypes := ""

	if g.Body != "" {
		params = params + "data : " + g.Body + ","
		paramsWithoutTypes = paramsWithoutTypes + "data,"
	}

	if g.Method == "" {
		g.Method = "GET"
	}

	splittedPaths := strings.Split(g.Path, "/")
	jsPaths := make([]string, 0)
	pathParams := make([]string, 0)

	for _, v := range splittedPaths {
		if strings.HasPrefix(v, ":") {
			param := strings.Replace(v, ":", "", -1)
			jsPaths = append(jsPaths, fmt.Sprintf(`${%s}`, param))
			pathParams = append(pathParams, param)

		} else {
			jsPaths = append(jsPaths, v)
		}
	}

	jsPath := ""

	for _, v := range jsPaths {
		if v != "" {
			jsPath = jsPath + "/" + v
		}
	}

	if len(pathParams) > 0 {
		for _, v := range pathParams {
			params = params + v + ":string ,"
			paramsWithoutTypes = paramsWithoutTypes + v + ","
		}
	}

	g.JsPath = jsPath
	g.FetchParams = params
	g.FetchParamsWithOutTypes = paramsWithoutTypes

}

func (g *FunctionDoc) Generate(s *strings.Builder) {
	responseType := ""
	if g.Response != "" {
		responseType = "<" + g.Response + ">"
	}
	s.WriteString(fmt.Sprintf("export const %s = async (%s options?: any) => {\n", g.LowerName, g.FetchParams))

	body := ""
	if g.Body != "" {
		if g.ContentType == "multipart/formdata" {
			s.WriteString("const formData = ShouldGenerateFormData(data)\n")
			body = "formData"
		} else {
			body = "data"
		}
	}

	s.WriteString(fmt.Sprintf("return await SafeParse(Fetcher%s({\n", responseType))
	s.WriteString("url :" + fmt.Sprintf("`%s`", g.JsPath) + ",\n")

	if body != "" {
		s.WriteString("body : " + body + ",\n")
	}
	s.WriteString("method :" + fmt.Sprintf(`"%s"`, g.Method) + ",\n")
	s.WriteString("headers : {\n")
	if g.ContentType == "multipart/formdata" {
		s.WriteString(`"Content-Type" : "multipart/formdata",` + "\n")
	}
	s.WriteString("...options?.headers\n")
	s.WriteString("},\n")
	s.WriteString("...options\n")
	s.WriteString("}))}\n")

}

// Extract JSON tag from struct field
func extractTags(tags string, tagName string) string {
	if strings.Contains(tags, tagName+":") {
		removedPrefix := strings.Index(tags, tagName+":")
		removed, found := strings.CutPrefix(tags[removedPrefix:], tagName+`:"`)

		if !found {
			return ""
		}
		return strings.Split(removed, `"`)[0]
	}
	return ""
}

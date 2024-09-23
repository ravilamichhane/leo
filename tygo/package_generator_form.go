package tygo

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func (g *PackageGenerator) GenerateForm(name string) (string, error) {
	s := new(strings.Builder)

	g.writeFileCodegenHeader(s)
	g.writeFileFrontmatter(s)

	filepaths := g.GoFiles

	availableImports := make([]string, 0)
	typespecmap := make(map[string]*ast.TypeSpec)
	functionsmap := make(map[string]*FunctionDoc)

	for i, file := range g.pkg.Syntax {

		if g.conf.IsFileIgnored(filepaths[i]) {
			continue
		}

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

				if isEmit {
					g.emitVar(s, x)
					return false
				}
				for _, spec := range x.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if ok && ts.Name.IsExported() {

						typespecmap[ts.Name.Name] = ts
						if ts.Name.Name != "" {
							availableImports = append(availableImports, ts.Name.Name)
						}
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
				functionsmap[a.Name] = a

				if a.LowerName != "" {

					availableImports = append(availableImports, a.LowerName)
				}
			}
			return true
		})

	}

	if v, ok := functionsmap[name]; ok {
		if v.ShouldGenerateForm {
			ts, ok := typespecmap[v.Body]
			if ok {
				st, isStruct := ts.Type.(*ast.StructType)
				if isStruct {

					s.WriteString(fmt.Sprintf(`import {%s} from "./api" 
					`, strings.Join(availableImports, ",")))

					s.WriteString(fmt.Sprintf("export const %sForm = () => {\n", v.Body))
					s.WriteString(fmt.Sprintf("const form = useForm<%s>({})\n", v.Body))

					s.WriteString(fmt.Sprintf(`
					  const %sMutation = useMutation({
mutationFn: async ({%s}:{%s}) => {
return %s(%s).then(([res, err]) => {
if (res) {
}
if(err) {
handleResponseError(err, form)
}
});

}});

const onSubmit = form.handleSubmit((data) => {
%sMutation.mutate({data});
});

return (
<Form {...form}>
<form onSubmit={onSubmit} className="grid grid-cols-1 lg:grid-cols-2">
`, v.LowerName, v.FetchParamsWithOutTypes, v.FetchParams, v.LowerName, v.FetchParamsWithOutTypes, v.LowerName))

					for _, field := range st.Fields.List {
						formtype := extractTags(field.Tag.Value, "formtype")
						fieldName := extractTags(field.Tag.Value, "json")
						switch formtype {
						case "string":
							s.WriteString(
								fmt.Sprintf(`<TextInput control={form.control} name={"%s"} label="%s"/>`, fieldName, fieldName))
						case "number":
							s.WriteString(
								fmt.Sprintf(`<TextInput control={form.control} type={"number"} name={"%s"} label="%s"/>`, fieldName, fieldName))
						case "email":
							s.WriteString(
								fmt.Sprintf(`<TextInput control={form.control} type={"email"} name={"%s"} label="%s"/>`, fieldName, fieldName))
						case "password":
							s.WriteString(
								fmt.Sprintf(`<TextInput control={form.control} type={"password"} name={"%s"} label="%s"/>`, fieldName, fieldName))
						case "date":
							s.WriteString(
								fmt.Sprintf(`<DateInput control={form.control}  name={"%s"} label="%s"/>`, fieldName, fieldName))

						case "select":
							s.WriteString(
								fmt.Sprintf(`  <SelectInput  control={form.control} name={"%s"}  placeholder="Select %s" options={[ ]} label="%s" />`, fieldName, fieldName, fieldName))

						case "switch":
							s.WriteString(fmt.Sprintf(` <SwitchInput control={form.control} name="%s" label="%s" />`, fieldName, fieldName))

						}
					}

					s.WriteString(fmt.Sprintf(`			
	<Button variant="outline" loading={ %sMutation.isPending} className="w-full col-span-2">
	  Submit
	</Button>
  </form>
</Form>)`, v.LowerName))

					s.WriteString("}\n")

				}
			}
		}
	}

	return s.String(), nil

}

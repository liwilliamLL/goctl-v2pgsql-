package template

// Field defines a filed template for types
var Field = `{{.name}} {{.type}} {{.tag}} {{if .hasComment}}// {{.comment}}{{end}}`


var FactoryFiled = `{{.name}}Model *model.{{.name}}Model`


var FactoryFuncFiled = `{{.name}}Model: model.New{{.name}}Model(dataSource) , `

package template

// Field defines a filed template for types
var Field = `{{.name}} {{.type}} {{.tag}} {{if .hasComment}}// {{.comment}}{{end}}`


var FactoryFiled = `{{.name}}Model *dto.{{.name}}Model`


var FactoryFuncFiled = `{{.name}}Model: dto.New{{.name}}Model(dataSource) , `

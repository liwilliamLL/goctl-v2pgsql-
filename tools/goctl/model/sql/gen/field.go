package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func genFields(fields []*parser.Field, primaryKey *parser.Field) (string, error) {
	var list []string

	for _, field := range fields {
		result, err := genField(field, primaryKey)
		if err != nil {
			return "", err
		}

		list = append(list, result)
	}

	return strings.Join(list, "\n"), nil
}

func genField(field *parser.Field, primaryKey *parser.Field) (string, error) {
	tag, err := genTag(field.Name.Source(), primaryKey.Name == field.Name)
	if err != nil {
		return "", err
	}

	text, err := util.LoadTemplate(category, fieldTemplateFile, template.Field)
	if err != nil {
		return "", err
	}

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"name":       field.Name.ToCamel(),
			"type":       field.DataType,
			"tag":        tag,
			"hasComment": field.Comment != "",
			"comment":    field.Comment,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func genFactoryFields(fields map[string]*model.Table)(string ,error){
	var list []string

	for _, field := range fields {
		result, err := genFactoryField(field)
		if err != nil {
			return "", err
		}

		list = append(list, result)
	}

	return strings.Join(list, "\n"), nil
}


func genFactoryField(field *model.Table)(string, error){
	text, err := util.LoadTemplate(Factory, factoryFiledsfile, template.FactoryFiled)
	if err != nil {
		return "", err
	}

	list, err := parser.ConvertDataType(field)
	if err != nil {
		return "", err
	}
	//println(tables.Name.ToCamel())

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"name":       list.Name.ToCamel(),
			//"type":       field.DataType,
			//"tag":        tag,
			//"hasComment": field.Comment != "",
			"comment":    field.Comment,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil

}



func genFactoryFuncFields(fields map[string]*model.Table)(string ,error){
	var list []string

	for _, field := range fields {
		result, err := genFactoryFuncField(field)
		if err != nil {
			return "", err
		}

		list = append(list, result)
	}

	return strings.Join(list, "\n"), nil
}


func genFactoryFuncField(field *model.Table)(string, error){
	text, err := util.LoadTemplate(Factory, factoryFuncFiledFile, template.FactoryFuncFiled)
	if err != nil {
		return "", err
	}

	list, err := parser.ConvertDataType(field)
	if err != nil {
		return "", err
	}
	//println(tables.Name.ToCamel())

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"name":       list.Name.ToCamel(),
			//"type":       field.DataType,
			//"tag":        tag,
			//"hasComment": field.Comment != "",
			//"comment":    field.Comment,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil

}
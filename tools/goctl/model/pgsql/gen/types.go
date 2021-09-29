package gen

import (
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/model/pgsql/model"
	"github.com/tal-tech/go-zero/tools/goctl/model/pgsql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func genTypes(table PgTable, methods, comment string, withCache bool) (string, error) {
	fields := table.Fields
	fieldsString, err := genFields(fields, &table.PrimaryKey.Field)
	if err != nil {
		return "", err
	}

	text, err := util.LoadTemplate(category, typesTemplateFile, template.Types)
	if err != nil {
		return "", err
	}

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"method":                methods,
			"upperStartCamelObject": table.Name.ToCamel(),
			"fields":                fieldsString,
			"comment":               comment,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func genPgFactoryTypes(pkg string, table map[string]*model.PgTable) (string, error) {
	//for _,k:=range table {
	//	tables, err := parser.ConvertDataType(k)
	//	if err != nil {
	//		return "", err
	//	}
	//}

	fieldsString, err := genPgFactoryFields(table)
	if err != nil {
		return "", err
	}

	text, err := util.LoadTemplate(Factory, factoryTypesFile, template.FactoryTypes)
	if err != nil {
		return "", err
	}

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			//"upperStartCamelObject": table.Name.ToCamel(),
			//"upperStartModelObject": table.Name.ToCamel(),
			"upperStartCamelObject": fmt.Sprintf("%s%s", UpdateUpper(pkg), "Dao"),
			"fields":                fieldsString,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}



func UpdateUpper(a string) string {
	vv := []rune(a)
	if len(vv) != 0 && vv[0] >= 97 && vv[0] <= 132 {
		vv[0] = vv[0] - 32
	}
	return string(vv)
}

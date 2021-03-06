package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/pgsql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func genNew(table PgTable, withCache bool) (string, error) {
	text, err := util.LoadTemplate(category, modelNewTemplateFile, template.New)
	if err != nil {
		return "", err
	}

	output, err := util.With("new").
		Parse(text).
		Execute(map[string]interface{}{
			"table":                 wrapWithRawString(table.Name.Source()),
			"withCache":             withCache,
			"upperStartCamelObject": table.Name.ToCamel(),
			"originTable":           table.Name.Source(),
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

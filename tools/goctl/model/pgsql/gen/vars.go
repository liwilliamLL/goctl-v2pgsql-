package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/pgsql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func genVars(table PgTable, withCache bool) (string, error) {
	keys := make([]string, 0)
	keys = append(keys, table.PrimaryCacheKey.VarExpression)
	for _, v := range table.UniqueCacheKey {
		keys = append(keys, v.VarExpression)
	}

	camel := table.Name.ToCamel()
	text, err := util.LoadTemplate(category, varTemplateFile, template.Vars)
	if err != nil {
		return "", err
	}

	output, err := util.With("var").Parse(text).
		GoFmt(true).Execute(map[string]interface{}{
		"lowerStartCamelObject": stringx.From(camel).Untitle(),
		"upperStartCamelObject": camel,
		"cacheKeys":             strings.Join(keys, "\n"),
		"autoIncrement":         table.PrimaryKey.AutoIncrement,
		"originalPrimaryKey":    wrapWithRawString(table.PrimaryKey.Name.Source()),
		"withCache":             withCache,
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

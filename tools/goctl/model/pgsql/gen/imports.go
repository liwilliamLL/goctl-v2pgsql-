package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/pgsql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"strings"
)

func genImports(table PgTable, withCache, timeImport, status bool) (string, error) {

	sql, importGorm := false, false

	for _, f := range table.Fields {
		if strings.Index(f.DataType, "sql.") != -1{
			sql = true
		}
		if f.Name.ToCamel() == "DeletedAt" {
			importGorm = true
		}
	}

	if withCache {
		text, err := util.LoadTemplate(category, importsTemplateFile, template.Imports)
		if err != nil {
			return "", err
		}

		buffer, err := util.With("import").Parse(text).Execute(map[string]interface{}{
			"time":   timeImport,
			"status": status,
			"sql":    sql,
			"gorm":   importGorm,
		})
		if err != nil {
			return "", err
		}

		return buffer.String(), nil
	}

	text, err := util.LoadTemplate(category, importsWithNoCacheTemplateFile, template.ImportsNoCache)
	if err != nil {
		return "", err
	}

	buffer, err := util.With("import").Parse(text).Execute(map[string]interface{}{
		"time":   timeImport,
		"status": status,
		"sql":    sql,
		"gorm":   importGorm,
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func genFactoryImport(pkg string) (string, error) {
	text, err := util.LoadTemplate(Factory, factoryImportsFile, template.FactoryImport)
	if err != nil {
		return "", err
	}

	buffer, err := util.With("import").Parse(text).Execute(map[string]interface{}{
		"time": false,
		"pkg":  pkg,
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

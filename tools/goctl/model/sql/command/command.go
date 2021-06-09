package command

import (
	"errors"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/urfave/cli"
)

const (
	flagSrc    = "src"
	flagDir    = "dir"
	flagOption = "o"
	flagCache  = "cache"
	flagIdea   = "idea"
	flagURL    = "url"
	flagTable  = "table"
	flagStyle  = "style"
	flagProto  = "proto"
	flagPPack  = "ppack"
)

var errNotMatched = errors.New("sql not matched")

// MysqlDDL generates model code from ddl
func MysqlDDL(ctx *cli.Context) error {
	src := ctx.String(flagSrc)
	dir := ctx.String(flagDir)
	cache := ctx.Bool(flagCache)
	idea := ctx.Bool(flagIdea)
	style := ctx.String(flagStyle)
	proto := ctx.String(flagProto)
	ppack := ctx.String(flagPPack)
	cfg, err := config.NewConfig(style, false)
	if err != nil {
		return err
	}

	return fromDDl(src, dir, proto, ppack, cfg, cache, idea)
}

// MyDataSource generates model code from datasource
func MyDataSource(ctx *cli.Context) error {
	url := strings.TrimSpace(ctx.String(flagURL))
	dir := strings.TrimSpace(ctx.String(flagDir))
	option := strings.TrimSpace(ctx.String(flagOption))
	cache := ctx.Bool(flagCache)
	idea := ctx.Bool(flagIdea)
	style := ctx.String(flagStyle)
	proto := ctx.String(flagProto)
	ppack := ctx.String(flagPPack)
	pattern := strings.TrimSpace(ctx.String(flagTable))
	cfg, err := config.NewConfig(style, false)
	if err != nil {
		return err
	}

	return fromDataSource(url, pattern, dir, proto, ppack, option, cfg, cache, idea)
}

func fromDDl(src, dir, proto, ppack string, cfg *config.Config, cache, idea bool) error {
	log := console.NewConsole(idea)
	src = strings.TrimSpace(src)
	if len(src) == 0 {
		return errors.New("expected path or path globbing patterns, but nothing found")
	}

	files, err := util.MatchFiles(src)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errNotMatched
	}

	var source []string
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		source = append(source, string(data))
	}

	generator, err := gen.NewDefaultGenerator(dir, proto, ppack, cfg, gen.WithConsoleOption(log))
	if err != nil {
		return err
	}

	return generator.StartFromDDL(strings.Join(source, "\n"), cache)
}

func fromDataSource(url, pattern, dir, proto, ppack, option string, cfg *config.Config, cache, idea bool) error {
	log := console.NewConsole(idea)
	if len(url) == 0 {
		log.Error("%v", "expected data source of mysql, but nothing found")
		return nil
	}

	if len(pattern) == 0 {
		log.Error("%v", "expected table or table globbing patterns, but nothing found")
		return nil
	}

	dsn, err := mysql.ParseDSN(url)
	if err != nil {
		return err
	}

	logx.Disable()
	databaseSource := strings.TrimSuffix(url, "/"+dsn.DBName) + "/information_schema"
	db := sqlx.NewMysql(databaseSource)
	im := model.NewInformationSchemaModel(db)

	tables, err := im.GetAllTables(dsn.DBName)
	if err != nil {
		return err
	}

	matchTables := make(map[string]*model.Table)
	for _, item := range tables {
		match, err := filepath.Match(pattern, item.TABLE_NAME)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		columnData, err := im.FindColumns(dsn.DBName, item.TABLE_NAME)
		if err != nil {
			return err
		}

		table, err := columnData.Convert()
		if err != nil {
			return err
		}

		//println(option)
		table.Comment = item.TABLE_COMMENT
		matchTables[item.TABLE_NAME] = table
	}

	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	for _, k := range matchTables {
		table, err := parser.ConvertDataType(k)
		if err != nil {
			return err
		}
		println(table.Name.ToCamel())
	}

	generator, err := gen.NewDefaultGenerator(dir, proto, ppack, cfg, gen.WithConsoleOption(log))
	if err != nil {
		return err
	}
	gens, err := gen.NewDefaultGenerator(option, proto, ppack, cfg, gen.WithConsoleOption(log))
	if err != nil {
		return err
	}
	err = gens.GenFactory(option, matchTables)
	return generator.StartFromInformationSchema(matchTables, cache)
}

package model

import (
	"gorm.io/gorm"
	"sort"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

const indexPri = "PRIMARY"

type (
	// InformationSchemaModel defines information schema model
	InformationSchemaModel struct {
		conn sqlx.SqlConn
	}

	InforPgmationSchemaModel struct {
		conn *gorm.DB
	}


	// Column defines column in table
	Column struct {
		*DbColumn
		Index *DbIndex
	}

	PgColumn struct {
		*PgDbColumn
		Index *DbIndex
	}

	// DbColumn defines column info of columns
	DbColumn struct {
		Name            string      `db:"COLUMN_NAME"`
		DataType        string      `db:"DATA_TYPE"`
		Extra           string      `db:"EXTRA"`
		Comment         string      `db:"COLUMN_COMMENT"`
		ColumnDefault   interface{} `db:"COLUMN_DEFAULT"`
		IsNullAble      string      `db:"IS_NULLABLE"`
		OrdinalPosition int         `db:"ORDINAL_POSITION"`
	}
	PgDbColumn struct {
		Name     	    string      `gorm:"column:column_name"`
		DataType        string      `gorm:"column:data_type"`
		Extra           string      `db:"EXTRA"`
		Comment         string      `db:"COLUMN_COMMENT"`
		ColumnDefault   interface{} `gorm:"column:column_default"`
		IsNullAble      string      `gorm:"column:is_nullable"`
		OrdinalPosition int         `gorm:"column:ordinal_position"`
	}

	// DbIndex defines index of columns in information_schema.statistic
	DbIndex struct {
		IndexName  string `db:"INDEX_NAME"`
		NonUnique  int    `db:"NON_UNIQUE"`
		SeqInIndex int    `db:"SEQ_IN_INDEX"`
	}

	// ColumnData describes the columns of table
	ColumnData struct {
		Db      string
		Table   string
		Columns []*Column
	}

	PgColumnData struct {
		Db      string
		Table   string
		Columns []*PgColumn
	}

	// Table describes mysql table which contains database name, table name, columns, keys
	Table struct {
		Db      string
		Table   string
		Comment string
		Columns []*Column
		// Primary key not included
		UniqueIndex map[string][]*Column
		PrimaryKey  *Column
		NormalIndex map[string][]*Column
	}

	PgTable struct {
		Db      string
		Table   string
		Comment string
		Columns []*PgColumn
		// Primary key not included
		UniqueIndex map[string][]*PgColumn
		PrimaryKey  *PgColumn
		NormalIndex map[string][]*PgColumn
	}

	// IndexType describes an alias of string
	IndexType string

	// Index describes a column index
	Index struct {
		IndexType IndexType
		Columns   []*Column
	}

	DBTable struct {
		TABLE_NAME string
		TABLE_COMMENT string
	}
	PgClass struct {
		TableName string `json:"table_name" xorm:"talbe_name"`
		TABLE_COMMENT string
	}
)

// NewInformationSchemaModel creates an instance for InformationSchemaModel
func NewInformationSchemaModel(conn sqlx.SqlConn) *InformationSchemaModel {
	return &InformationSchemaModel{conn: conn}
}
// NewInformationSchemaModel creates an instance for InformationSchemaModel
func NewPGInformationSchemaModel(conn *gorm.DB) *InforPgmationSchemaModel {
	return &InforPgmationSchemaModel{conn: conn}
}



// GetAllTables selects all tables from TABLE_SCHEMA
func (m *InformationSchemaModel) GetAllTables(database string) ([]*DBTable, error) {
	query := `select TABLE_NAME, TABLE_COMMENT from TABLES where TABLE_SCHEMA = ?`
	var tables []*DBTable
	err := m.conn.QueryRows(&tables, query, database)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

// GetAllTables selects all tables from TABLE_SCHEMA
func (m *InforPgmationSchemaModel) SGetAllTables(database string) ([]*PgClass, error) {
	//query := `select relname as TABLE_NAME,col_description(c.oid, 0) as  TABLE_COMMENT from pg_class c where relkind = 'r' and relname  like 'o_%' order by relname`
	query := `select relname as TABLE_NAME,col_description(c.oid, 0) as  TABLE_COMMENT from pg_class c where relkind = 'r' and relname  = 'o_account' order by relname`
	tables:=[]*PgClass{}
	err := m.conn.Raw(query).Find(&tables).Error
	if err != nil {
		return nil, err
	}
	return tables, nil
}

// FindColumns return columns in specified database and table
func (m *InforPgmationSchemaModel) FindColumns(db, table string) (*PgColumnData, error) {
	//querySql := `SELECT c.COLUMN_NAME,c.DATA_TYPE,EXTRA,c.COLUMN_COMMENT,c.COLUMN_DEFAULT,c.IS_NULLABLE,c.ORDINAL_POSITION from COLUMNS c WHERE c.TABLE_SCHEMA = ? and c.TABLE_NAME = ? `
	querySql :=`SELECT * from information_schema.columns c WHERE c.TABLE_NAME = ? `

	var reply []*PgDbColumn
	err := m.conn.Raw(querySql,table).Find(&reply).Error
	//err := m.conn.QueryRowsPartial(&reply, querySql, db, table)
	if err != nil {
		return nil, err
	}

	var list []*PgColumn
	for _, item := range reply {
		//index, err := m.FindIndex(db, table, item.Name)
		//if err != nil {
		//	if err != sqlx.ErrNotFound {
		//		return nil, err
		//	}
		//	continue
		//}
		//
		//if len(index) > 0 {
		//	for _, i := range index {
		//		list = append(list, &PgColumn{
		//			PgDbColumn: item,
		//			Index:    i,
		//		})
		//	}
		//} else {
		//	list = append(list, &PgColumn{
		//		PgDbColumn: item,
		//	})
		//}

		list =append(list,&PgColumn{
			PgDbColumn: item,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].OrdinalPosition < list[j].OrdinalPosition
	})

	var columnData PgColumnData
	columnData.Db = db
	columnData.Table = table
	columnData.Columns = list
	return &columnData, nil
}

// FindIndex finds index with given db, table and column.
func (m *InforPgmationSchemaModel) FindIndex(db, table, column string) ([]*DbIndex, error) {
	querySql := `SELECT s.INDEX_NAME,s.NON_UNIQUE,s.SEQ_IN_INDEX from  STATISTICS s  WHERE s.TABLE_NAME = ? and s.COLUMN_NAME = ?`
	var reply []*DbIndex
	err := m.conn.Raw(querySql,table,column).Find(&reply).Error
	//err := m.conn.QueryRowsPartial(&reply, querySql, db, table, column)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Convert converts column data into PgTable
func (c *PgColumnData) Convert() (*PgTable, error) {
	var table PgTable
	table.Table = c.Table
	table.Db = c.Db
	table.Columns = c.Columns
	table.UniqueIndex = map[string][]*PgColumn{}
	table.NormalIndex = map[string][]*PgColumn{}

	m := make(map[string][]*PgColumn)
	for _, each := range c.Columns {
		if each.Index != nil {
			m[each.Index.IndexName] = append(m[each.Index.IndexName], each)
		}
	}

	//primaryColumns := m[indexPri]
	//if len(primaryColumns) == 0 {
	//	return nil, fmt.Errorf("db:%s, table:%s, missing primary key", c.Db, c.Table)
	//}
	//
	//if len(primaryColumns) > 1 {
	//	return nil, fmt.Errorf("db:%s, table:%s, joint primary key is not supported", c.Db, c.Table)
	//}

	//table.PrimaryKey = primaryColumns[0]
	//for indexName, columns := range m {
	//	if indexName == indexPri {
	//		continue
	//	}
	//
	//	for _, one := range columns {
	//		if one.Index != nil {
	//			if one.Index.NonUnique == 0 {
	//				table.UniqueIndex[indexName] = columns
	//			} else {
	//				table.NormalIndex[indexName] = columns
	//			}
	//		}
	//	}
	//}

	return &table, nil
}

package template

var (
	// Imports defines a import template for model in cache case
	Imports = `import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)
`
	// ImportsNoCache defines a import template for model in normal case
	ImportsNoCache = `import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)
`
	FactoryImport =`import (
	"erp-server/lib/provider/mysql"
	"erp-server/src/dao/{{.pkg}}/model"
)
`


)

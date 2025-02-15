package tables

import "embed"

//go:embed *.sql
var CreateTableSqls embed.FS

package main

import "fmt"

type tableInfo struct {
	OID          uint32
	Name         string
	OrderColumns []columnInfo
}
type columnInfo struct {
	Name     string
	OrderNum uint64
}

func setTableInfo() {
	// Get list of table
	var (
		queryTableName = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
		tableNames     []string
	)
	rows, err := dbConn.Query(queryTableName)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		tableNames = append(tableNames, name)
	}
	rows.Close()

	// Get the info's columns of table
	var (
		queryTableColumns = func(name string) string {
			return fmt.Sprintf("SELECT attrelid, attname, attnum FROM pg_attribute WHERE attrelid = 'public.%v'::regclass AND NOT attisdropped AND attnum >= 0 ORDER BY attnum", name)
		}
	)
	for idx := range tableNames {
		var (
			tableName       = tableNames[idx]
			query           = queryTableColumns(tableName)
			columnRows, err = dbConn.Query(query)
		)
		if err != nil {
			panic(err)
		}

		for columnRows.Next() {
			var (
				attrelid uint32
				attname  string
				attnum   uint64
			)
			err := columnRows.Scan(&attrelid, &attname, &attnum)
			if err != nil {
				panic(err)
			}
			var table, ok = mapTable[attrelid]
			if !ok {
				table = &tableInfo{
					OID:  attrelid,
					Name: tableName,
				}
				mapTable[attrelid] = table
			}
			table.OrderColumns = append(table.OrderColumns, columnInfo{
				Name:     attname,
				OrderNum: attnum,
			})
		}
		columnRows.Close()
	}
}

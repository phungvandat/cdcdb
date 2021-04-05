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
	fmt.Println("table config loading...")
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
		var (
			tableID     uint32
			columnInfos = []columnInfo{}
		)
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
			tableID = attrelid
			columnInfos = append(columnInfos, columnInfo{
				Name:     attname,
				OrderNum: attnum,
			})
		}
		columnRows.Close()
		mapTable[tableID] = &tableInfo{
			OID:          tableID,
			Name:         tableName,
			OrderColumns: columnInfos,
		}
	}

	fmt.Println("table config loaded!")
}

package core

import "fmt"

// TODO: prepared statements
func createDatabaseSQL() string {
	return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", databaseName)
}

func createTableSQL() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT PRIMARY KEY AUTO_INCREMENT, %s VARCHAR(255) NOT NULL);", table, idField, dataField)
}

func insertSQL(data string) string {
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES (\"%s\");", table, dataField, data)
}

func queryOneSQL(order orderType) string {
	orderString := "DESC"
	if order == AscendingOrder {
		orderString = "ASC"
	}
	return fmt.Sprintf("SELECT * FROM %s ORDER BY %s %s LIMIT 1;", table, idField, orderString)
}

func deleteSQL(id int) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s=%v;", table, idField, id)
}

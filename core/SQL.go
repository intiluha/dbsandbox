package main

import "fmt"

func createDatabaseSQL() string {
	return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", DatabaseName)
}

func createTableSQL() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT PRIMARY KEY AUTO_INCREMENT, %s VARCHAR(255) NOT NULL);", Table, IdField, DataField)
}

func insertSQL(data string) string {
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES (\"%s\");", Table, DataField, data)
}

func queryOneSQL() string {
	return fmt.Sprintf("SELECT * FROM %s ORDER BY %s %s LIMIT 1;", Table, IdField, Order)
}

func deleteSQL(id int) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s=%v;", Table, IdField, id)
}

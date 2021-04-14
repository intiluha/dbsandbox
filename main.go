package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/intiluha/dbsandbox/core"
	"log"
)

func main() {
	err := core.CreateDatabase()
	if err != nil {
		log.Fatal(err)
	}
	err = core.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	errs := make(chan error)
	nWriters, nReaders, nOperations := 2, 2, 100
	for i := 0; i < nWriters; i++ {
		go core.Writer(string(rune('a'+i)), nOperations, errs)
	}
	for i := 0; i < nReaders; i++ {
		go core.Reader(nOperations, errs)
	}
	for i := 0; i < nWriters + nReaders; i++ {
		log.Print(<- errs)
	}
}

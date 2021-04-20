package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/intiluha/dbsandbox/core"
	"sync"
)

func main() {
	err := core.CreateDatabase()
	core.Assert(err)
	err = core.CreateTable()
	core.Assert(err)

	nWriters, nReaders, nOperations := 5, 5, 10
	wg := new(sync.WaitGroup)
	wg.Add(nWriters+nReaders)
	for i := 0; i < nWriters; i++ {
		go core.Writer(string(rune('a'+i)), nOperations, wg)
	}
	for i := 0; i < nReaders; i++ {
		go core.Reader(core.DescendingOrder, nOperations, wg)
	}
	wg.Wait()
}

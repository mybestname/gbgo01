package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

// You Application logic here
func main() {
	err := doBusiness();
	if err!=nil {
		if errors.Is(err, sql.ErrNoRows) {      // pkg errors also support go1.13 compatible API
			fmt.Printf("%+v\n", err)     // pkg.errors can print call stack
		}
		if errors.Cause(err) == sql.ErrNoRows { // the original pkg.errors Cause() API
			// do handle the sentinel error
		}
	}
}

func doDao() error {
	return errors.Wrap(sql.ErrNoRows, "DAO failed") // wrap with dao context
}

func doBusiness() error {
	// handle error
	if err:= doDao(); err!= nil {
		return errors.WithMessagef(err,"do Business failed"); //add business context
	}
	// do stuff
	return nil
}


package main

import (
	"database/sql"
	"errors"
	"fmt"
)

// You Application logic here
func main() {
 	err := doBusiness();
 	if err!=nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Can still handle sentinel error
			fmt.Printf("%+v\n", err)
		}else{
			// unknown error
			fmt.Printf("%+v\n", err)
		}
	}
}

func doDao() error {
	return fmt.Errorf("DAO failed: %w", sql.ErrNoRows) // wrap with dao context
}

func doBusiness() error {
	// handle error
	if err:= doDao(); err!= nil {
		return fmt.Errorf("do Business failed : %w", err); //wrap with business context
	}
	// do stuff
	return nil
}

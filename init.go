// init
package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

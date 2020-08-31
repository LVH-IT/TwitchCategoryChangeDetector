package main

import (
	"log"
	"os"
)

func checkError(errorVar error) {
	if errorVar != nil {
		log.Fatal(errorVar)
		os.Exit(1)
	}
}

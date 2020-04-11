package main

import (
	"flag"
	"log"

	"github.com/kingzbauer/gmail"
)

var (
	cnfFile = flag.String("c", "", "path to credentials file")
	tknFile = flag.String("t", "", "path to token file")
)

func chk(msg string, err error) {
	if err != nil {
		log.Fatalf("Error %s => %s", msg, err)
	}
}

func main() {
	flag.Parse()

	if len(*cnfFile) == 0 {
		flag.Usage()
		log.Fatal("c flag is required")
	}

	tkn, err := gmail.GetUserToken(*cnfFile, *tknFile)
	chk("Retrieve user token", err)
	err = gmail.TokenToFile(*tknFile, tkn)
	chk("Save token", err)
}

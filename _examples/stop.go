package main

import (
	"context"
	"flag"
	"log"

	"github.com/kingzbauer/gmail/auth"
)

var (
	c = flag.String("c", "", "config file")
)

func chk(msg string, err error) {
	if err != nil {
		log.Fatalf("Error %s => %s", msg, err)
	}
}

func main() {
	flag.Parse()
	if len(*c) == 0 {
		flag.Usage()
		log.Fatal("c flag is required")
	}

	cnf, err := auth.GetJWTConfig(*c)
	chk("Get client config", err)
	cnf.Subject = "support@wezago.com"
	srv, err := auth.NewService(cnf.TokenSource(context.Background()))
	chk("Get gmail service", err)

	call := srv.Users.Stop("support@wezago.com")
	chk("Stop watch", call.Do())
}

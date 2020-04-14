package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/gmail/v1"

	"github.com/kingzbauer/gmail/auth"
)

var (
	c = flag.String("c", "", "config file")
	t = flag.String("t", "", "token file")
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

	watchReq := &gmail.WatchRequest{
		LabelIds:  []string{"UNREAD"},
		TopicName: "projects/wezago/topics/emails",
	}

	call := srv.Users.Watch("support@wezago.com", watchReq)
	watchRes, err := call.Do()
	chk("Watch Response", err)

	// Convert the milli seconds into seconds
	secs := watchRes.Expiration / 1000
	t := time.Unix(secs, 0)
	nanos := watchRes.Expiration * 1000000
	tM := time.Unix(0, nanos)
	fmt.Printf("Expiration: %s\n", t)
	fmt.Printf("Expiration2: %s\n", tM)
	fmt.Printf("HistoryId: %d\n", watchRes.HistoryId)
}

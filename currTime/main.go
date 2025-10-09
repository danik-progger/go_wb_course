package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func PrintTime() {
	for {
		t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
		if err != nil {
			fmt.Fprintf(os.Stderr, "🔴 Error queriing time from NTP: %s", err)
			os.Exit(0)
		}

		fmt.Printf("Current time is: %v\n", t.Second())
		time.Sleep(1 * time.Second)
	}
}

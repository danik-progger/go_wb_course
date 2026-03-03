package main

import (
	"fmt"
	"currTime/pkg"
	"os"
)

func main() {
	t, err := pkg.GetNTPTime("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "🔴 Error querying time from NTP: %s\n", err)
		os.Exit(1)
	}

	pkg.PrintTime(t)
}

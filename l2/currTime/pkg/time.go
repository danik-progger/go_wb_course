package pkg

import (
	"fmt"
	"time"

	"github.com/beevik/ntp"
)

func GetNTPTime(server string) (time.Time, error) {
	return ntp.Time(server)
}

func PrintTime(t time.Time) {
	fmt.Printf("Current time is: %s\n", t.Format(time.RFC3339))
}

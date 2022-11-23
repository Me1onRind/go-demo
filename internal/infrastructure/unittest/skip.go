package unittest

import (
	"os"
)

func SkipCauseIO() {
	if os.Getenv("skip_io") == "1" {
		os.Exit(0)
	}
}

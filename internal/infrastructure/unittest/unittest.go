package unittest

import "os"

func SkipCauseExternalIO() {
	if os.Getenv("skip_external_io") == "1" {
		os.Exit(0)
	}
}

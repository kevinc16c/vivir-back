package config_time_zone

import (
	"os"
)

func init() {
	os.Setenv("TZ", "America/Buenos_Aires")
}

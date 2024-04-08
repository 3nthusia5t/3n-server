package log

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

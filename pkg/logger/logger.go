package logger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"convenienceStore/pkg/config"
)

// FromConfig builds a ready to use logger from configuration.
func FromConfig(cfg config.LoggingConfig) *log.Logger {
	flags := log.LstdFlags
	if strings.EqualFold(cfg.Format, "detailed") {
		flags |= log.Lshortfile
	}

	prefix := "[convenienceStore] "
	if cfg.Level != "" {
		prefix = fmt.Sprintf("[convenienceStore][%s] ", strings.ToUpper(cfg.Level))
	}

	return log.New(os.Stdout, prefix, flags)
}

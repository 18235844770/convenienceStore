package logger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"convenienceStore/pkg/config"
)

// FromConfig 根据配置构建一套开箱即用的日志记录器。
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

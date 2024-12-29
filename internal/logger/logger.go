package logger

import (
	"fmt"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
)

func Debug(format string, args ...interface{}) {
	if config.DEBUG {
		fmt.Printf(format, args...)
	}
}

package debug

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const watchFile string = "/debug-attach"

func Wait(log *slog.Logger) {
	log.Info("Waiting for debugger to be attached...")
	for {
		_, err := os.Stat(watchFile)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Error(fmt.Sprintf("Error watching %s", watchFile), "error", err.Error())
				os.Exit(-1)
			}
		} else {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	err := os.Remove(watchFile)
	if err != nil {
		log.Error(fmt.Sprintf("Error removing %s", watchFile), "error", err.Error())
		os.Exit(-1)
	}
	log.Info("Debugger attached")
}

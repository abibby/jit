package jitlog

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"
)

func Logger(name string) (*slog.Logger, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %v", err)
	}

	dir := path.Join(home, ".config/jit/logs")
	err = os.MkdirAll(dir, 0o777)
	if err != nil {
		return nil, fmt.Errorf("failed to create log dir: %v", err)
	}

	day := time.Now().Format(time.DateOnly)
	f, err := os.OpenFile(path.Join(dir, day+"-"+name+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log: %v", err)
	}
	return slog.New(slog.NewJSONHandler(f, nil)), nil
}

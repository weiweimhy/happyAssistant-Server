package logger

import (
	"github.com/sirupsen/logrus"
	"happyAssistant/internal/config"
	"log"
	"os"
)

func InitLogger(config config.LogConfig) {
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// 设置日志输出文件
	if config.File != "" {
		f, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("open logger file error: %v", err)
		}
		log.SetOutput(f)
	}
}

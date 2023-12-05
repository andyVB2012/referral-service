// logger/logger.go
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	// Initialize logger
	Log.Out = os.Stdout
	Log.SetFormatter(&logrus.JSONFormatter{})
}

package common

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func ParseFlags() {
	err := godotenv.Load()
	if err != nil {
		logrus.Debugf("loading .env file %v", err)
	}

	flag.VisitAll(func(f *flag.Flag) {
		name := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if value, ok := os.LookupEnv(name); ok {
			err = flag.Set(f.Name, value)
			if err != nil {
				logrus.Fatalf("failed setting flag from environment: %v", err)
				return
			}
		}
	})

	flag.Parse()
}

func SetupLogger(level string) (*logrus.Logger, error) {
	log := logrus.StandardLogger()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	log.SetLevel(lvl)
	return log, err
}

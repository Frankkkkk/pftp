package main

import (
	"os"
	"os/signal"
	"syscall"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/pyama86/pftp/pftp"
	"github.com/sirupsen/logrus"
)

var ftpServer *pftp.FtpServer

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	stackLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
	logrus.AddHook(logrus_stack.NewHook(stackLevels, stackLevels))
}

func main() {
	confFile := "./example.toml"

	ftpServer, err := pftp.NewFtpServer(confFile)
	if err != nil {
		logrus.Fatal(err)
	}

	ftpServer.Use("user", User)
	if err := ftpServer.ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}
}

func signalHandler() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	for {
		switch <-ch {
		case syscall.SIGTERM:
			ftpServer.Stop()
			break
		}
	}
}

func User(c *pftp.Context, param string) error {
	c.RemoteAddr = "192.168.33.2:21"
	return nil
}

package logging

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

func NewLogger() *logrus.Logger {
	Logger = logrus.New()
	Logger.SetFormatter(&formatter{})
	Logger.SetOutput(newService())

	return Logger
}

type service struct{}

func newService() *service { return &service{} }

func (s *service) Write(b []byte) (n int, err error) {
	cmd := exec.Command("notify-send", "hlwm", string(b))
	return 0, cmd.Start()
}

type formatter struct{}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

package testutil

import (
	"io"

	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = io.Discard
	return mute.WithField("test", true)
}

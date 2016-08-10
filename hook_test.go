package xmpprus

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/facebookgo/ensure"
)

func TestLevels(t *testing.T) {
	lvls := levelsFrom(logrus.DebugLevel)
	ensure.DeepEqual(t, lvls, []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel})

	lvls = levelsFrom(logrus.ErrorLevel)
	ensure.DeepEqual(t, lvls, []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel})
}

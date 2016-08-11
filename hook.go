package xmpprus

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/agl/xmpp-client/xmpp"
)

// Hook is a Logrus Hook that sends logs via XMPP.
type Hook struct {
	levels    []logrus.Level
	formatter logrus.Formatter // ignoring default formatter, otherwise TTY causes problems while sending

	xmppConn    *xmpp.Conn
	xmppConnMtx sync.Mutex

	receiver string
}

// NewHook initializes a Hook that can send logs via XMPP.
// serverHost may be empty, then auto-detection will be attempted.
func NewHook(minLogLevel logrus.Level, receiver, sender, password, serverHost string) (*Hook, error) {
	addressParts := strings.Split(sender, "@")
	if len(addressParts) != 2 {
		return nil, errors.New("invalid XMPP user")
	}

	if serverHost == "" {
		host, port, err := xmpp.Resolve(addressParts[1])
		if err != nil {
			return nil, fmt.Errorf("cannot resolve server name automatically, please specify server host (%v)", err)
		}
		serverHost = fmt.Sprintf("%s:%d", host, port)
	}

	conn, err := xmpp.Dial(serverHost, addressParts[0], addressParts[1], "logrus", password, &xmpp.Config{})
	if err != nil {
		return nil, err
	}
	return &Hook{
		levels:    levelsFrom(minLogLevel),
		xmppConn:  conn,
		receiver:  receiver,
		formatter: &logrus.TextFormatter{DisableColors: true},
	}, nil
}

func (h *Hook) Levels() []logrus.Level {
	return h.levels
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	buf, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	h.xmppConnMtx.Lock()
	err = h.xmppConn.Send(h.receiver, string(buf))
	h.xmppConnMtx.Unlock()
	return err
}

var levels = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

func levelsFrom(l logrus.Level) []logrus.Level {
	for i := range levels {
		if levels[i] == l {
			return levels[i:]
		}
	}
	return []logrus.Level{}
}

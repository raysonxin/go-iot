package drivers

import "github.com/sirupsen/logrus"

type Driver struct {
	Name string //instance identifier
}

type Logger struct {
	logger *logrus.Entry
}

type Adaptor interface {
	Name() string
	SetName(name string)
	SetLogger(log *logrus.Entry)
}

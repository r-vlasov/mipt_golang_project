package targethandlers

import (
	"time"
)

type TargetInterface interface {
	Init(interface{})
	Monitor(time time.Duration, alertchannel chan<- map[string]string)
}

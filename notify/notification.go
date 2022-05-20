package notificator

import (
	"time"
	"sync"
)

type NotifierUnitInterface interface {
	Notify(msg string)
}

type Notifier struct {
        AlertChannel	chan map[string]string
	LastNotified	map[string]time.Time
	MutexMap	sync.RWMutex // concurrent read/write from map
	Delay		time.Duration
	Premsg		string
	Notificators	[]NotifierUnitInterface
}

func (n *Notifier) Init(delay int, premsg string, notificators []NotifierUnitInterface) {
	n.AlertChannel = make(chan map[string]string)
	n.Premsg = premsg
	n.Delay = time.Duration(delay) * time.Second
	n.LastNotified = make(map[string]time.Time)
	n.Notificators = notificators
}

func (n *Notifier) updateTimeNotification(n_id string, ts time.Time) {
	n.MutexMap.Lock()
	n.LastNotified[n_id] = ts
	n.MutexMap.Unlock()
}

func (n *Notifier) GetAlertChannel() (chan map[string]string) {
	return n.AlertChannel
}


// true -> already notified
// false -> should be notified
func (n *Notifier) checkAlreadyNotified(n_id string) (bool) {
	n.MutexMap.RLock()
	ts, found := n.LastNotified[n_id]
	n.MutexMap.RUnlock()

	if found {
		ts_now := time.Now()
		if (ts_now.Sub(ts) > n.Delay) {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

func (n *Notifier) handleAlert(n_id string, msg *string) {
	if !n.checkAlreadyNotified(n_id) {
		for _, notif_unit := range n.Notificators {
			go notif_unit.Notify(n.Premsg + *msg)
		}
		n.updateTimeNotification(n_id, time.Now())
	}
}

func (n *Notifier) RunNotify() {
	for {
		select {
		case alert := <-n.AlertChannel:
			for id, msg := range alert {
				go n.handleAlert(id, &msg)
			}
		}
	}
}




/*
func main() {
	n := new(Notifier)
	tn := new(TelegramNotifier)
	tn.Init("5132582314:AAEVW_xw4JvRPWb-HQ9_VU8i62L3-QFJQ8U", 5055419231)

	n.Init(30, "asd", []NotifierUnitInterface{&fmtNotifier{}, tn})
	go n.RunNotify()
	for {
		n.AlertChannel <- map[string]string {"blyaha" : "blya"}
	}
}
*/

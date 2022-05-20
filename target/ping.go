package targethandlers

import (
	"time"
	"sync"
	"fmt"
	"github.com/go-ping/ping"
	config"monitoring/config"
)

type PingTarget struct {
	EndPoints 	[]config.PingTargetConfig
	Pingers 	[]*ping.Pinger
}


func (p *PingTarget) Init(raw_conf []config.PingTargetConfig) {
	p.EndPoints = raw_conf

	for _, conf := range p.EndPoints {
		p.Pingers = append(p.Pingers, p.createPinger(conf))
	}
}

func (p *PingTarget) sendAlert(name string, host string, problem string, alertchannel chan<- map[string]string) {
	msg := fmt.Sprintf("[PING] %s:%s is unreachable - %s", name, host, problem)
	alertchannel <- map[string]string {
		name + ":" + host : msg,
	}
}

func (p *PingTarget) createPinger(conf config.PingTargetConfig) (*ping.Pinger) {
	
	host := conf.Host
	repeat := conf.Repeat
	timeout := conf.Timeout

	// create handler
	ping_handler, err := ping.NewPinger(host)
	if err != nil {
		return nil
	}

	interval_between := time.Duration(1)

	ping_handler.Count = repeat
	ping_handler.Interval = interval_between * time.Second
	ping_handler.Timeout = time.Duration(timeout) * time.Second
	ping_handler.SetPrivileged(true)
	return ping_handler
}

func (p *PingTarget) Monitor(delay time.Duration, alertchannel chan<- map[string]string) {
	var wg sync.WaitGroup

	runnerNotifier := func(pinger *ping.Pinger, conf config.PingTargetConfig) {
		defer wg.Done()

		recv_prev := pinger.Statistics().PacketsRecv
		_ = pinger.Run()
		recv_cur := pinger.Statistics().PacketsRecv
		if recv_cur <= recv_prev {
			//fmt.Println("Unreachable ping", pinger.Addr())
			p.sendAlert(conf.Name, conf.Host, "ping is unreachable", alertchannel)
		}
	}
	
	for {
		for i, pngr := range p.Pingers {
			if pngr != nil {
				wg.Add(1)
				go runnerNotifier(pngr, p.EndPoints[i])
			}
		}
		wg.Wait()

		// refresh pingers array (there may be balancers that give ip by round-robin or smth else)
		for i, conf := range p.EndPoints {
			p.Pingers[i] = p.createPinger(conf)
			if p.Pingers[i] == nil {
				p.sendAlert(conf.Name, conf.Host, "can't resolve hostname", alertchannel)
				//fmt.Println("Failed to start pinger (can't resolve domain name)")
			}
		}
		time.Sleep(delay)
	}
}


/*
func main() {

	var p PingTarget;
	a := map[string]interface{} {
		"Hostname" : "www.google.com",
		"Timeout" : 2,
		"Repeat" : 2,
	}
	c := map[string]interface{} {
		"Hostname" : "1.1.13.104",
		"Timeout" : 1,
		"Repeat" : 2, 
	}
	p.Init([]map[string]interface{}{a, c})
	p.Monitor(10 * time.Second);
}
*/

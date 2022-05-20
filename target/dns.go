package targethandlers

import (
	"time"
	"context"
	"sync"
	"fmt"
	"net"
	config "monitoring/config"
)

type DnsTarget struct {
	EndPoints 	[]config.DnsTargetConfig
	Resolvers 	[]*net.Resolver
}


//func (d *DnsTarget) Init(raw_conf []map[string]interface{}) {

func (d *DnsTarget) Init(raw_conf []config.DnsTargetConfig) {
	d.EndPoints = raw_conf
	//fmt.Println(d.EndPoints)
	
	for _, conf := range d.EndPoints {
		timeout := conf.Timeout
		dns_server := conf.Server
		d.Resolvers = append(d.Resolvers, &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				dialer := net.Dialer{
					Timeout: time.Duration(timeout) * time.Second,
				}
				return dialer.DialContext(ctx, network, dns_server)
			},
		})
        }
}

func (d *DnsTarget) sendAlert(name string, hostname string, server string, alertchannel chan<- map[string]string) {
	msg := fmt.Sprintf("[DNS] %s:%s is not resolved. DNS server: %s", name, hostname, server)
	alertchannel <- map[string]string {
		name + ":" + hostname : msg,
	}
}

func (d *DnsTarget) runnerNotifier(
		wg *sync.WaitGroup,
		conf config.DnsTargetConfig, 
		resolver *net.Resolver,
		alertchannel chan<- map[string]string) {
	
	defer wg.Done()
	hostname := conf.Hostname
	server := conf.Server
	name := conf.Name
	repeat := conf.Repeat
	timeout := conf.Timeout

	for i := 0; i < repeat; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout) * time.Second)
		defer cancel()

		ip, _ := resolver.LookupHost(ctx, hostname)
		if len(ip) > 0 {
			return
		}

	}
	d.sendAlert(name, hostname, server, alertchannel)
}

func (d *DnsTarget) Monitor(delay time.Duration, alertchannel chan<- map[string]string) {
	var wg sync.WaitGroup
	
	for {
		for i, rslv := range d.Resolvers {
			if rslv != nil {
				wg.Add(1)
				go d.runnerNotifier(&wg, d.EndPoints[i], rslv, alertchannel)
			}
		}
		wg.Wait()
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

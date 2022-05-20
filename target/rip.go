package targethandlers

import (
	"time"
	"fmt"
	"net"
	"strconv"
	"io"
	config "monitoring/config"
)

type RipTarget struct {
	EndPoints 		[]config.RipTargetConfig
	ListenSock		net.PacketConn
	LastAnnouncedTime	map[string]time.Time
}


//func (d *DnsTarget) Init(raw_conf []map[string]interface{}) {

func (d *RipTarget) Init(raw_conf []config.RipTargetConfig) {
	d.EndPoints = raw_conf
	sock, err := net.ListenPacket("udp4", ":520")
	if err != nil {
		panic("Can't listen 520 port for RIP")
	}
	d.ListenSock = sock
}

func (d *RipTarget) sendAlert(name string, address string, alertchannel chan<- map[string]string) {
	msg := fmt.Sprintf("[RIP] %s:%s is not anounced", name, address)
	alertchannel <- map[string]string {
		name + ":" + address : msg,
	}
}

//func (d *RipTarget) runnerNotifier(alertchannel chan<- map[string]string) {


func (r *RipTarget) runnerListener() {
	buf := make([]byte, 24 * 8) // protocol RIP
	for {
		n, addr, err := r.ListenSock.ReadFrom(buf)
		fmt.Println(n, addr)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read data from socket error")
			}
		}
		ip := buf[8:12]
		fmt.Println(strconv.Itoa(int(ip[0])) + "." + strconv.Itoa(int(ip[1])) +
			"." + strconv.Itoa(int(ip[2])) + "." + strconv.Itoa(int(ip[3])))
	}
}

func (r *RipTarget) Monitor(_ time.Duration, alertchannel chan<- map[string]string) {
	//go d.runnerNotifier(alertchannel)
	go r.runnerListener()
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
